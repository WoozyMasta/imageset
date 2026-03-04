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
	"strings"
)

const (
	maxIntValue = int(^uint(0) >> 1)
	minIntValue = -maxIntValue - 1
)

// Parse decodes .imageset text from reader into Document.
func Parse(reader io.Reader) (*Document, error) {
	return parseFromReader(reader)
}

// ParseBytes decodes .imageset text bytes into Document.
func ParseBytes(data []byte) (*Document, error) {
	return Parse(bytes.NewReader(data))
}

// ParseString decodes .imageset text from a string.
func ParseString(data string) (*Document, error) {
	return Parse(strings.NewReader(data))
}

// ParseFile decodes .imageset file from disk.
func ParseFile(path string) (*Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	doc, err := Parse(file)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	return doc, nil
}

// parseFromReader decodes .imageset text from reader.
func parseFromReader(reader io.Reader) (*Document, error) {
	document := &Document{}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var (
		seenRoot     bool
		inRoot       bool
		inTextures   bool
		inImages     bool
		inGroups     bool
		inGroupImage bool

		currentTexture *Texture
		currentGroup   *Group
		currentImage   *Image
	)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := trimASCIISpace(stripLineCommentBytes(scanner.Bytes()))
		if len(line) == 0 {
			continue
		}

		switch {
		case hasBlockOpenBytes(line, "ImageSetClass"):
			seenRoot = true
			inRoot = true
			continue

		case hasBlockOpenBytes(line, "Textures"):
			inTextures = true
			continue

		case hasBlockOpenBytes(line, "Groups"):
			inGroups = true
			continue

		case hasBlockOpenBytes(line, "Images"):
			if currentGroup != nil {
				inGroupImage = true
			} else {
				inImages = true
			}
			continue

		case hasBlockOpenBytes(line, "ImageSetTextureClass"):
			currentTexture = &Texture{Mpix: 1}
			continue

		case hasBlockOpenBytes(line, "ImageSetGroupClass"):
			currentGroup = &Group{}
			if className := parseClassNameBytes(line); className != "" {
				currentGroup.Name = className
			}
			inGroupImage = false
			continue

		case hasBlockOpenBytes(line, "ImageSetDefClass"):
			currentImage = &Image{}
			continue
		}

		if len(line) == 1 && line[0] == '{' {
			continue
		}

		if len(line) == 1 && line[0] == '}' {
			if currentImage != nil {
				if currentGroup != nil {
					currentGroup.Images = append(currentGroup.Images, *currentImage)
				} else {
					document.Images = append(document.Images, *currentImage)
				}
				currentImage = nil
				continue
			}

			if currentTexture != nil {
				document.Textures = append(document.Textures, *currentTexture)
				currentTexture = nil
				continue
			}

			if inGroupImage {
				inGroupImage = false
				continue
			}

			if currentGroup != nil && inGroups {
				document.Groups = append(document.Groups, *currentGroup)
				currentGroup = nil
				continue
			}

			if inTextures {
				inTextures = false
				continue
			}

			if inImages {
				inImages = false
				continue
			}

			if inGroups {
				inGroups = false
				continue
			}

			if inRoot {
				inRoot = false
				continue
			}

			return nil, &ParseError{
				Line:    lineNumber,
				Cause:   ErrInvalidSyntax,
				Message: "unexpected block close",
			}
		}

		key, rawValue := splitKeyAndRestBytes(line)
		if len(key) == 0 {
			continue
		}

		if !inRoot {
			return nil, &ParseError{
				Line:    lineNumber,
				Cause:   ErrInvalidSyntax,
				Message: "unexpected content outside ImageSetClass",
			}
		}

		switch {
		case equalToken(key, "Name"):
			value := parseStringValueBytes(rawValue)
			switch {
			case currentImage != nil:
				currentImage.Name = value

			case currentGroup != nil:
				currentGroup.Name = value

			default:
				document.Name = value
			}

		case equalToken(key, "RefSize"):
			if len(rawValue) == 0 {
				return nil, newParseError(lineNumber, "invalid RefSize", ErrInvalidSyntax)
			}

			size, err := parseSizeValueBytes(rawValue)
			if err != nil {
				return nil, newParseError(lineNumber, "invalid RefSize values", err)
			}
			document.RefSize = size

		case equalToken(key, "Pos"):
			if currentImage == nil || len(rawValue) == 0 {
				return nil, newParseError(lineNumber, "invalid Pos", ErrInvalidSyntax)
			}

			point, err := parsePointValueBytes(rawValue)
			if err != nil {
				return nil, newParseError(lineNumber, "invalid Pos values", err)
			}
			currentImage.Pos = point

		case equalToken(key, "Size"):
			if currentImage == nil || len(rawValue) == 0 {
				return nil, newParseError(lineNumber, "invalid Size", ErrInvalidSyntax)
			}

			size, err := parseSizeValueBytes(rawValue)
			if err != nil {
				return nil, newParseError(lineNumber, "invalid Size values", err)
			}
			currentImage.Size = size

		case equalToken(key, "Flags"):
			if currentImage == nil || len(rawValue) == 0 {
				return nil, newParseError(lineNumber, "invalid Flags", ErrInvalidSyntax)
			}

			flags, err := parseFlagsExprBytes(rawValue)
			if err != nil {
				return nil, newParseError(lineNumber, "invalid Flags value", err)
			}
			currentImage.Flags = flags

		case equalToken(key, "mpix"), equalToken(key, "Mpix"):
			if currentTexture == nil || len(rawValue) == 0 {
				return nil, newParseError(lineNumber, "invalid mpix", ErrInvalidSyntax)
			}

			value, err := parseIntTokenBytes(rawValue)
			if err != nil {
				return nil, newParseError(lineNumber, "invalid mpix value", err)
			}
			currentTexture.Mpix = value

		case equalToken(key, "path"), equalToken(key, "Path"):
			if currentTexture == nil || len(rawValue) == 0 {
				return nil, newParseError(lineNumber, "invalid path", ErrInvalidSyntax)
			}
			currentTexture.Path = parseStringValueBytes(rawValue)

		default:
			// keep parser tolerant for unknown fields.
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read input: %w", err)
	}

	if currentImage != nil || currentTexture != nil || currentGroup != nil ||
		inTextures || inImages || inGroups || inGroupImage || inRoot {
		return nil, &ParseError{
			Line:    lineNumber,
			Cause:   ErrInvalidSyntax,
			Message: "unclosed block at end of file",
		}
	}
	if !seenRoot {
		return nil, &ParseError{
			Line:    lineNumber,
			Cause:   ErrInvalidSyntax,
			Message: "missing ImageSetClass root block",
		}
	}

	return document, nil
}

// hasBlockOpenBytes reports whether line opens a class/section block.
func hasBlockOpenBytes(line []byte, head string) bool {
	line = trimASCIISpace(line)
	if len(line) < len(head)+1 || line[len(line)-1] != '{' {
		return false
	}

	for idx := 0; idx < len(head); idx++ {
		if line[idx] != head[idx] {
			return false
		}
	}

	if len(line) == len(head)+1 {
		return true
	}

	return isASCIISpace(line[len(head)])
}

// parseClassNameBytes extracts class identifier from "ClassName id {" lines.
func parseClassNameBytes(line []byte) string {
	line = trimASCIISpace(line)
	if len(line) == 0 {
		return ""
	}

	if line[len(line)-1] == '{' {
		line = trimASCIISpace(line[:len(line)-1])
	}

	firstEnd := indexASCIISpace(line)
	if firstEnd == -1 {
		return ""
	}

	secondStart := firstEnd
	for secondStart < len(line) && isASCIISpace(line[secondStart]) {
		secondStart++
	}
	if secondStart >= len(line) {
		return ""
	}

	secondEnd := secondStart
	for secondEnd < len(line) && !isASCIISpace(line[secondEnd]) {
		secondEnd++
	}

	name := line[secondStart:secondEnd]
	if len(name) == 0 {
		return ""
	}

	if name[len(name)-1] == '{' {
		name = name[:len(name)-1]
	}

	if len(name) == 0 {
		return ""
	}

	return string(name)
}

// splitKeyAndRestBytes splits "Key value..." bytes into key and value segment.
func splitKeyAndRestBytes(line []byte) (key, rest []byte) {
	line = trimASCIISpace(line)
	if len(line) == 0 {
		return nil, nil
	}

	spaceIndex := indexASCIISpace(line)
	if spaceIndex == -1 {
		return line, nil
	}

	return line[:spaceIndex], trimASCIISpace(line[spaceIndex+1:])
}

// splitPairBytes extracts exactly two tokens from a value segment.
func splitPairBytes(raw []byte) (left, right []byte, ok bool) {
	raw = trimASCIISpace(raw)
	if len(raw) == 0 {
		return nil, nil, false
	}

	firstSpace := indexASCIISpace(raw)
	if firstSpace == -1 {
		return nil, nil, false
	}

	left = trimASCIISpace(raw[:firstSpace])
	if len(left) == 0 {
		return nil, nil, false
	}

	tail := trimASCIISpace(raw[firstSpace+1:])
	if len(tail) == 0 {
		return nil, nil, false
	}

	secondSpace := indexASCIISpace(tail)
	if secondSpace == -1 {
		return left, tail, true
	}

	right = trimASCIISpace(tail[:secondSpace])
	if len(right) == 0 {
		return nil, nil, false
	}

	extra := trimASCIISpace(tail[secondSpace+1:])
	if len(extra) != 0 {
		return nil, nil, false
	}

	return left, right, true
}

// parseStringValueBytes parses raw value segment, supporting quoted strings.
func parseStringValueBytes(raw []byte) string {
	raw = trimASCIISpace(raw)
	if len(raw) == 0 {
		return ""
	}

	if raw[0] == '"' && raw[len(raw)-1] == '"' {
		quoted := raw[1 : len(raw)-1]
		if bytes.IndexByte(quoted, '\\') == -1 && bytes.IndexByte(quoted, '"') == -1 {
			return string(quoted)
		}

		unquoted, err := strconv.Unquote(string(raw))
		if err == nil {
			return unquoted
		}
	}

	if raw[0] == '"' {
		raw = raw[1:]
	}
	if len(raw) > 0 && raw[len(raw)-1] == '"' {
		raw = raw[:len(raw)-1]
	}

	return string(raw)
}

// parseSizeValueBytes parses "width height" value segment.
func parseSizeValueBytes(raw []byte) (Size, error) {
	width, height, ok := splitPairBytes(raw)
	if !ok {
		return Size{}, ErrInvalidSyntax
	}

	return parseSizeBytes(width, height)
}

// parsePointValueBytes parses "x y" value segment.
func parsePointValueBytes(raw []byte) (Point, error) {
	xValue, yValue, ok := splitPairBytes(raw)
	if !ok {
		return Point{}, ErrInvalidSyntax
	}

	return parsePointBytes(xValue, yValue)
}

// parseSizeBytes parses width and height integer values.
func parseSizeBytes(width, height []byte) (Size, error) {
	w, err := parseIntTokenBytes(width)
	if err != nil {
		return Size{}, err
	}

	h, err := parseIntTokenBytes(height)
	if err != nil {
		return Size{}, err
	}

	return Size{Width: w, Height: h}, nil
}

// parsePointBytes parses x and y integer values.
func parsePointBytes(xValue, yValue []byte) (Point, error) {
	x, err := parseIntTokenBytes(xValue)
	if err != nil {
		return Point{}, err
	}

	y, err := parseIntTokenBytes(yValue)
	if err != nil {
		return Point{}, err
	}

	return Point{X: x, Y: y}, nil
}

// parseIntTokenBytes parses signed decimal integer token.
func parseIntTokenBytes(raw []byte) (int, error) {
	raw = trimASCIISpace(raw)
	if len(raw) == 0 {
		return 0, ErrInvalidSyntax
	}

	sign := 1
	index := 0
	switch raw[0] {
	case '-':
		sign = -1
		index = 1
	case '+':
		index = 1
	}

	if index >= len(raw) {
		return 0, ErrInvalidSyntax
	}

	maxAbs := uint64(maxIntValue)
	if sign < 0 {
		maxAbs++
	}

	var value uint64
	for ; index < len(raw); index++ {
		ch := raw[index]
		if ch < '0' || ch > '9' {
			return 0, ErrInvalidSyntax
		}

		digit := uint64(ch - '0')
		if value > (maxAbs-digit)/10 {
			return 0, ErrInvalidSyntax
		}

		value = value*10 + digit
	}

	if sign < 0 {
		if value == uint64(maxIntValue)+1 {
			return minIntValue, nil
		}
		return -int(value), nil
	}

	return int(value), nil
}

// stripLineCommentBytes removes // comment not enclosed in quotes.
func stripLineCommentBytes(line []byte) []byte {
	inString := false
	escaped := false

	for idx := 0; idx < len(line)-1; idx++ {
		ch := line[idx]
		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' && inString {
			escaped = true
			continue
		}

		if ch == '"' {
			inString = !inString
			continue
		}

		if !inString && ch == '/' && line[idx+1] == '/' {
			return line[:idx]
		}
	}

	return line
}

// equalToken reports whether token bytes match the literal value.
func equalToken(token []byte, literal string) bool {
	if len(token) != len(literal) {
		return false
	}

	for idx := 0; idx < len(literal); idx++ {
		if token[idx] != literal[idx] {
			return false
		}
	}

	return true
}

// indexASCIISpace returns first ASCII space index or -1.
func indexASCIISpace(raw []byte) int {
	for idx, ch := range raw {
		if isASCIISpace(ch) {
			return idx
		}
	}

	return -1
}

// trimASCIISpace trims ASCII spaces from both sides.
func trimASCIISpace(raw []byte) []byte {
	start := 0
	for start < len(raw) && isASCIISpace(raw[start]) {
		start++
	}

	end := len(raw)
	for end > start && isASCIISpace(raw[end-1]) {
		end--
	}

	return raw[start:end]
}

// isASCIISpace reports whether byte is ASCII whitespace.
func isASCIISpace(ch byte) bool {
	switch ch {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

// newParseError builds parse error with line and wrapped cause.
func newParseError(line int, message string, cause error) error {
	return &ParseError{
		Line:    line,
		Cause:   cause,
		Message: message,
	}
}
