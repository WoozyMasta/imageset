// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"fmt"
	"strconv"
	"strings"
)

// Diagnostic describes one validation issue.
type Diagnostic struct {
	Path    string `json:"path" yaml:"path"`       // Field path.
	Message string `json:"message" yaml:"message"` // Human message.
}

// ValidationError aggregates semantic validation diagnostics.
type ValidationError struct {
	Diagnostics []Diagnostic `json:"diagnostics" yaml:"diagnostics"` // Issues list.
}

// imageNameRef stores first-seen image location for duplicate checks.
type imageNameRef struct {
	groupIndex int // Group index, -1 for root images section.
	imageIndex int // Image index inside root or group section.
}

// Error formats validation diagnostics as one sentence.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Diagnostics) == 0 {
		return "imageset: validation failed"
	}

	first := e.Diagnostics[0]
	if len(e.Diagnostics) == 1 {
		return fmt.Sprintf("imageset: %s: %s", first.Path, first.Message)
	}

	return fmt.Sprintf(
		"imageset: %s: %s (and %d more)",
		first.Path,
		first.Message,
		len(e.Diagnostics)-1,
	)
}

// Validate checks semantic constraints and returns aggregated error.
func Validate(document *Document) error {
	if document == nil {
		return ErrNilDocument
	}

	diagnostics := collectDiagnostics(document)
	if len(diagnostics) == 0 {
		return nil
	}

	return &ValidationError{Diagnostics: diagnostics}
}

// collectDiagnostics performs semantic checks and collects issues.
func collectDiagnostics(document *Document) []Diagnostic {
	diagnostics := make([]Diagnostic, 0, 8)

	if document.RefSize.Width <= 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    "ref_size.width",
			Message: "must be > 0",
		})
	}
	if document.RefSize.Height <= 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    "ref_size.height",
			Message: "must be > 0",
		})
	}

	seenNames := make(map[string]imageNameRef, len(document.Images))
	for textureIndex, texture := range document.Textures {
		if strings.TrimSpace(texture.Path) == "" {
			diagnostics = append(diagnostics, Diagnostic{
				Path:    textureFieldPath(textureIndex, "path"),
				Message: "must be non-empty",
			})
		}
		if texture.Mpix < 0 {
			diagnostics = append(diagnostics, Diagnostic{
				Path:    textureFieldPath(textureIndex, "mpix"),
				Message: "must be >= 0",
			})
		}
	}

	for imageIndex, item := range document.Images {
		diagnostics = validateImage(
			diagnostics,
			item,
			-1,
			imageIndex,
			document.RefSize,
			seenNames,
		)
	}

	for groupIndex, group := range document.Groups {
		if strings.TrimSpace(group.Name) == "" {
			diagnostics = append(diagnostics, Diagnostic{
				Path:    groupFieldPath(groupIndex, "name"),
				Message: "must be non-empty",
			})
		}

		for imageIndex, item := range group.Images {
			diagnostics = validateImage(
				diagnostics,
				item,
				groupIndex,
				imageIndex,
				document.RefSize,
				seenNames,
			)
		}
	}

	return diagnostics
}

// validateImage validates one image and updates diagnostics.
func validateImage(
	diagnostics []Diagnostic,
	item Image,
	groupIndex, imageIndex int,
	refSize Size,
	seenNames map[string]imageNameRef,
) []Diagnostic {
	name := strings.TrimSpace(item.Name)
	if name == "" {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    imageFieldPath(groupIndex, imageIndex, "name"),
			Message: "must be non-empty",
		})
	} else {
		firstRef, exists := seenNames[name]
		if exists {
			diagnostics = append(diagnostics, Diagnostic{
				Path: imageFieldPath(groupIndex, imageIndex, "name"),
				Message: "duplicate name, first seen at " +
					imageFieldPath(firstRef.groupIndex, firstRef.imageIndex, "name"),
			})
		} else {
			seenNames[name] = imageNameRef{
				groupIndex: groupIndex,
				imageIndex: imageIndex,
			}
		}
	}

	if item.Pos.X < 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    imageFieldPath(groupIndex, imageIndex, "pos.x"),
			Message: "must be >= 0",
		})
	}
	if item.Pos.Y < 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    imageFieldPath(groupIndex, imageIndex, "pos.y"),
			Message: "must be >= 0",
		})
	}

	if item.Size.Width <= 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    imageFieldPath(groupIndex, imageIndex, "size.width"),
			Message: "must be > 0",
		})
	}
	if item.Size.Height <= 0 {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    imageFieldPath(groupIndex, imageIndex, "size.height"),
			Message: "must be > 0",
		})
	}

	if refSize.Width > 0 &&
		item.Size.Width > 0 &&
		item.Pos.X >= 0 &&
		item.Pos.X+item.Size.Width > refSize.Width {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    imagePath(groupIndex, imageIndex),
			Message: "out of bounds by width against ref_size",
		})
	}

	if refSize.Height > 0 &&
		item.Size.Height > 0 &&
		item.Pos.Y >= 0 &&
		item.Pos.Y+item.Size.Height > refSize.Height {
		diagnostics = append(diagnostics, Diagnostic{
			Path:    imagePath(groupIndex, imageIndex),
			Message: "out of bounds by height against ref_size",
		})
	}

	return diagnostics
}

// textureFieldPath returns "textures[i].field".
func textureFieldPath(index int, field string) string {
	return "textures[" + strconv.Itoa(index) + "]." + field
}

// groupFieldPath returns "groups[i].field".
func groupFieldPath(groupIndex int, field string) string {
	return "groups[" + strconv.Itoa(groupIndex) + "]." + field
}

// imagePath returns image item path in root or group section.
func imagePath(groupIndex, imageIndex int) string {
	if groupIndex < 0 {
		return "images[" + strconv.Itoa(imageIndex) + "]"
	}

	return "groups[" + strconv.Itoa(groupIndex) +
		"].images[" + strconv.Itoa(imageIndex) + "]"
}

// imageFieldPath returns image item field path in root or group section.
func imageFieldPath(groupIndex, imageIndex int, field string) string {
	return imagePath(groupIndex, imageIndex) + "." + field
}
