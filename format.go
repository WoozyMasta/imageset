// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	// defaultIndent is used when FormatOptions.Indent is empty.
	defaultIndent = "\t"
)

// FormatOptions controls .imageset formatting behavior.
type FormatOptions struct {
	// Indentation string for one level.
	Indent string `json:"indent,omitempty" yaml:"indent,omitempty"`

	// Normalize names as CamelCase.
	UseCamelCaseNames bool `json:"camel_case,omitempty" yaml:"camel_case,omitempty"`
}

// writerState keeps buffered writer and format options.
type writerState struct {
	writer       io.Writer       // Raw output writer.
	stringWriter io.StringWriter // Fast string writer.
	byteWriter   io.ByteWriter   // Optional byte writer for small tokens.
	flusher      *bufio.Writer   // Optional buffered flusher.
	opts         FormatOptions   // Effective format options.
	intScratch   [20]byte        // Reused integer formatting buffer.
}

// Write serializes document into .imageset text form.
func Write(writer io.Writer, document *Document, opts *FormatOptions) error {
	if document == nil {
		return ErrNilDocument
	}

	state := newWriterState(writer, resolveFormatOptions(opts))
	if err := writeDocument(state, document, 0); err != nil {
		return err
	}

	if state.flusher != nil {
		if err := state.flusher.Flush(); err != nil {
			return fmt.Errorf("flush output: %w", err)
		}
	}

	return nil
}

// Format serializes document and returns encoded bytes.
func Format(document *Document, opts *FormatOptions) ([]byte, error) {
	if document == nil {
		return nil, ErrNilDocument
	}

	options := resolveFormatOptions(opts)

	var buffer bytes.Buffer
	buffer.Grow(estimateFormatSize(document, options))
	if err := Write(&buffer, document, &options); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// WriteFile serializes document and writes the result to file path.
func WriteFile(path string, document *Document, opts *FormatOptions) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() {
		closeErr := file.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("close file: %w", closeErr)
		}
	}()

	if err := Write(file, document, opts); err != nil {
		return fmt.Errorf("write document: %w", err)
	}

	return nil
}

// resolveFormatOptions applies defaults for optional format options.
func resolveFormatOptions(opts *FormatOptions) FormatOptions {
	if opts == nil {
		return FormatOptions{Indent: defaultIndent}
	}

	options := *opts
	if options.Indent == "" {
		options.Indent = defaultIndent
	}

	return options
}

// estimateFormatSize estimates output size to reduce buffer growth.
func estimateFormatSize(document *Document, opts FormatOptions) int {
	indentLen := len(opts.Indent)
	if indentLen == 0 {
		indentLen = len(defaultIndent)
	}

	size := 64
	size += len(document.Name) + 32
	size += 32

	for _, texture := range document.Textures {
		size += (indentLen * 2) + 64 + len(texture.Path)
	}
	for _, item := range document.Images {
		size += estimateImageSize(item, indentLen)
	}
	for _, group := range document.Groups {
		size += estimateGroupSize(group, indentLen)
	}

	return size
}

// estimateGroupSize returns estimated encoded size for one group block.
func estimateGroupSize(group Group, indentLen int) int {
	size := (indentLen * 2) + 64 + len(group.Name)
	for _, item := range group.Images {
		size += estimateImageSize(item, indentLen)
	}

	return size
}

// estimateImageSize returns estimated encoded size for one image block.
func estimateImageSize(item Image, indentLen int) int {
	size := (indentLen * 2) + 96 + len(item.Name)
	size += 96
	return size
}

// newWriterState builds buffered writer state for formatting.
func newWriterState(writer io.Writer, opts FormatOptions) *writerState {
	// File output benefits from fewer syscalls, keep buffered path here.
	if _, ok := writer.(*os.File); ok {
		buffer := bufio.NewWriterSize(writer, 32*1024)
		return &writerState{
			writer:       buffer,
			stringWriter: buffer,
			byteWriter:   buffer,
			flusher:      buffer,
			opts:         opts,
		}
	}

	if sw, ok := writer.(io.StringWriter); ok {
		bw, _ := writer.(io.ByteWriter)
		return &writerState{
			writer:       writer,
			stringWriter: sw,
			byteWriter:   bw,
			opts:         opts,
		}
	}

	buffer := bufio.NewWriterSize(writer, 32*1024)
	return &writerState{
		writer:       buffer,
		stringWriter: buffer,
		byteWriter:   buffer,
		flusher:      buffer,
		opts:         opts,
	}
}

// writeDocument writes root Document block with nested content.
func writeDocument(state *writerState, document *Document, level int) error {
	if err := state.writeBlockOpen(level, "ImageSetClass"); err != nil {
		return err
	}

	name := NormalizeName(document.Name, state.opts.UseCamelCaseNames)
	if err := state.writeKeyQuoted(level+1, "Name", name); err != nil {
		return err
	}
	if err := state.writeKeyIntPair(
		level+1,
		"RefSize",
		document.RefSize.Width,
		document.RefSize.Height,
	); err != nil {
		return err
	}

	if len(document.Textures) > 0 {
		if err := state.writeBlockOpen(level+1, "Textures"); err != nil {
			return err
		}

		for _, texture := range document.Textures {
			if err := writeTexture(state, texture, level+2); err != nil {
				return err
			}
		}

		if err := state.writeBlockClose(level + 1); err != nil {
			return err
		}
	}

	if len(document.Images) > 0 {
		if err := state.writeBlockOpen(level+1, "Images"); err != nil {
			return err
		}

		for _, item := range document.Images {
			if err := writeImage(state, item, level+2); err != nil {
				return err
			}
		}

		if err := state.writeBlockClose(level + 1); err != nil {
			return err
		}
	}

	if len(document.Groups) > 0 {
		if err := state.writeBlockOpen(level+1, "Groups"); err != nil {
			return err
		}

		for _, group := range document.Groups {
			if err := writeGroup(state, group, level+2); err != nil {
				return err
			}
		}

		if err := state.writeBlockClose(level + 1); err != nil {
			return err
		}
	}

	if err := state.writeBlockClose(level); err != nil {
		return err
	}

	return nil
}

// writeTexture writes ImageSetTextureClass block.
func writeTexture(state *writerState, texture Texture, level int) error {
	if err := state.writeBlockOpen(level, "ImageSetTextureClass"); err != nil {
		return err
	}
	if err := state.writeKeyInt(level+1, "mpix", texture.Mpix); err != nil {
		return err
	}
	if err := state.writeKeyQuoted(level+1, "path", texture.Path); err != nil {
		return err
	}

	if err := state.writeBlockClose(level); err != nil {
		return err
	}

	return nil
}

// writeImage writes ImageSetDefClass block.
func writeImage(state *writerState, item Image, level int) error {
	normalizedName := NormalizeName(item.Name, state.opts.UseCamelCaseNames)
	className := normalizedName
	if normalizedName == "" {
		className = "default"
	}

	if err := state.writeBlockOpenWithName(level, "ImageSetDefClass", className); err != nil {
		return err
	}

	if err := state.writeKeyQuoted(level+1, "Name", normalizedName); err != nil {
		return err
	}
	if err := state.writeKeyIntPair(level+1, "Pos", item.Pos.X, item.Pos.Y); err != nil {
		return err
	}
	if err := state.writeKeyIntPair(
		level+1,
		"Size",
		item.Size.Width,
		item.Size.Height,
	); err != nil {
		return err
	}

	// Keep DayZ-like textual form for single flags and numeric form for combined.
	if err := state.writeKeyString(level+1, "Flags", formatFlagsForWrite(item.Flags)); err != nil {
		return err
	}

	if err := state.writeBlockClose(level); err != nil {
		return err
	}

	return nil
}

// formatFlagsForWrite returns canonical DayZ-like representation for flags.
func formatFlagsForWrite(flags Flags) string {
	switch flags {
	case 0:
		return "0"
	case FlagHorizontalTile:
		return "ISHorizontalTile"
	case FlagVerticalTile:
		return "ISVerticalTile"
	case FlagHorizontalTile | FlagVerticalTile:
		return "3"
	default:
		return strconv.Itoa(int(flags))
	}
}

// writeGroup writes ImageSetGroupClass block.
func writeGroup(state *writerState, group Group, level int) error {
	normalizedName := NormalizeName(group.Name, state.opts.UseCamelCaseNames)
	className := normalizedName
	if normalizedName == "" {
		className = "default"
	}

	if err := state.writeBlockOpenWithName(level, "ImageSetGroupClass", className); err != nil {
		return err
	}
	if err := state.writeKeyQuoted(level+1, "Name", normalizedName); err != nil {
		return err
	}

	if len(group.Images) > 0 {
		if err := state.writeBlockOpen(level+1, "Images"); err != nil {
			return err
		}

		for _, item := range group.Images {
			if err := writeImage(state, item, level+2); err != nil {
				return err
			}
		}

		if err := state.writeBlockClose(level + 1); err != nil {
			return err
		}
	}

	if err := state.writeBlockClose(level); err != nil {
		return err
	}

	return nil
}

// writeBlockOpen writes "<name> {" line at indentation level.
func (state *writerState) writeBlockOpen(level int, name string) error {
	if err := state.writeIndent(level); err != nil {
		return err
	}
	if err := state.writeString(name); err != nil {
		return err
	}
	if err := state.writeString(" {\n"); err != nil {
		return err
	}

	return nil
}

// writeBlockOpenWithName writes "<kind> <name> {" line at indentation level.
func (state *writerState) writeBlockOpenWithName(level int, kind, name string) error {
	if err := state.writeIndent(level); err != nil {
		return err
	}
	if err := state.writeString(kind); err != nil {
		return err
	}
	if err := state.writeString(" "); err != nil {
		return err
	}
	if err := state.writeString(name); err != nil {
		return err
	}
	if err := state.writeString(" {\n"); err != nil {
		return err
	}

	return nil
}

// writeBlockClose writes "}" line at indentation level.
func (state *writerState) writeBlockClose(level int) error {
	if err := state.writeIndent(level); err != nil {
		return err
	}
	if err := state.writeString("}\n"); err != nil {
		return err
	}

	return nil
}

// writeKeyQuoted writes a quoted key-value line at indentation level.
func (state *writerState) writeKeyQuoted(level int, key, value string) error {
	if err := state.writeIndent(level); err != nil {
		return err
	}
	if err := state.writeString(key); err != nil {
		return err
	}
	if err := state.writeString(" "); err != nil {
		return err
	}
	if err := state.writeQuoted(value); err != nil {
		return err
	}
	if err := state.writeString("\n"); err != nil {
		return err
	}

	return nil
}

// writeKeyString writes a non-quoted key-value line at indentation level.
func (state *writerState) writeKeyString(level int, key, value string) error {
	if err := state.writeIndent(level); err != nil {
		return err
	}
	if err := state.writeString(key); err != nil {
		return err
	}
	if err := state.writeString(" "); err != nil {
		return err
	}
	if err := state.writeString(value); err != nil {
		return err
	}
	if err := state.writeString("\n"); err != nil {
		return err
	}

	return nil
}

// writeKeyInt writes an integer key-value line at indentation level.
func (state *writerState) writeKeyInt(level int, key string, value int) error {
	if err := state.writeIndent(level); err != nil {
		return err
	}
	if err := state.writeString(key); err != nil {
		return err
	}
	if err := state.writeString(" "); err != nil {
		return err
	}
	if err := state.writeInt(value); err != nil {
		return err
	}
	if err := state.writeString("\n"); err != nil {
		return err
	}

	return nil
}

// writeKeyIntPair writes a two-integer key-value line at indentation level.
func (state *writerState) writeKeyIntPair(level int, key string, left, right int) error {
	if err := state.writeIndent(level); err != nil {
		return err
	}
	if err := state.writeString(key); err != nil {
		return err
	}
	if err := state.writeString(" "); err != nil {
		return err
	}
	if err := state.writeInt(left); err != nil {
		return err
	}
	if err := state.writeString(" "); err != nil {
		return err
	}
	if err := state.writeInt(right); err != nil {
		return err
	}
	if err := state.writeString("\n"); err != nil {
		return err
	}

	return nil
}

// writeIndent writes indentation prefix for one line.
func (state *writerState) writeIndent(level int) error {
	for range level {
		if err := state.writeString(state.opts.Indent); err != nil {
			return err
		}
	}

	return nil
}

// writeString writes plain string without extra allocation.
func (state *writerState) writeString(value string) error {
	_, err := state.stringWriter.WriteString(value)
	return err
}

// writeInt writes decimal integer without fmt allocation overhead.
func (state *writerState) writeInt(value int) error {
	out := strconv.AppendInt(state.intScratch[:0], int64(value), 10)

	if state.byteWriter != nil {
		for _, ch := range out {
			if err := state.byteWriter.WriteByte(ch); err != nil {
				return err
			}
		}
		return nil
	}

	_, err := state.writer.Write(out)
	return err
}

// writeQuoted writes Go-quoted string representation.
func (state *writerState) writeQuoted(value string) error {
	if isSimpleQuotedASCII(value) {
		if err := state.writeString("\""); err != nil {
			return err
		}
		if err := state.writeString(value); err != nil {
			return err
		}
		if err := state.writeString("\""); err != nil {
			return err
		}

		return nil
	}

	if err := state.writeString(strconv.Quote(value)); err != nil {
		return err
	}

	return nil
}

// isSimpleQuotedASCII reports whether value can be written in quotes as-is.
func isSimpleQuotedASCII(value string) bool {
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if ch < 0x20 || ch > 0x7e || ch == '"' || ch == '\\' {
			return false
		}
	}

	return true
}
