// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/woozymasta/lintkit/lint"
)

// Diagnostic describes one validation issue.
type Diagnostic struct {
	Severity lint.Severity `json:"severity" yaml:"severity"` // Issue level.
	Path     string        `json:"path" yaml:"path"`         // Field path.
	Message  string        `json:"message" yaml:"message"`   // Human message.
	Code     lint.Code     `json:"code" yaml:"code"`         // Stable code.
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

// knownImageFlagMask stores known DayZ .imageset flag bits.
const knownImageFlagMask = FlagHorizontalTile | FlagVerticalTile

const (
	// defaultMinImagePadding stores default minimum gap for optional padding check.
	defaultMinImagePadding = 4
)

// ValidateOptions configures optional semantic checks.
type ValidateOptions struct {
	// EnablePaddingCheck enables optional minimum image-gap check.
	EnablePaddingCheck bool `json:"enable_padding_check,omitempty" yaml:"enable_padding_check,omitempty"`

	// MinPadding stores minimum allowed gap between image rectangles.
	// Used only when EnablePaddingCheck is true.
	MinPadding int `json:"min_padding,omitempty" yaml:"min_padding,omitempty"`
}

// ImagePaddingRuleOptions documents lint rule options for padding check.
type ImagePaddingRuleOptions struct {
	// MinPadding stores minimum allowed gap between image rectangles.
	MinPadding int `json:"min_padding" yaml:"min_padding"`
}

// imageRectRef stores image rectangle metadata for pair checks.
type imageRectRef struct {
	groupIndex int
	imageIndex int
	x          int
	y          int
	w          int
	h          int
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
	return ValidateWithOptions(document, nil)
}

// ValidateWithOptions checks semantic constraints and returns aggregated error.
func ValidateWithOptions(document *Document, options *ValidateOptions) error {
	if document == nil {
		return ErrNilDocument
	}

	resolvedOptions := resolveValidateOptions(options)
	diagnostics := collectDiagnostics(document, resolvedOptions)
	if len(diagnostics) == 0 {
		return nil
	}

	return &ValidationError{Diagnostics: diagnostics}
}

// resolveValidateOptions applies defaults for optional validation settings.
func resolveValidateOptions(options *ValidateOptions) ValidateOptions {
	if options == nil {
		return ValidateOptions{
			EnablePaddingCheck: false,
			MinPadding:         defaultMinImagePadding,
		}
	}

	resolved := *options
	if resolved.MinPadding <= 0 {
		resolved.MinPadding = defaultMinImagePadding
	}

	return resolved
}

// errorDiagnostic builds one error-level diagnostic.
func errorDiagnostic(code lint.Code, path string, message string) Diagnostic {
	return diagnostic(code, lint.SeverityError, path, message)
}

// warningDiagnostic builds one warning-level diagnostic.
func warningDiagnostic(code lint.Code, path string, message string) Diagnostic {
	return diagnostic(code, lint.SeverityWarning, path, message)
}

// infoDiagnostic builds one info-level diagnostic.
func infoDiagnostic(code lint.Code, path string, message string) Diagnostic {
	return diagnostic(code, lint.SeverityInfo, path, message)
}

// diagnostic builds one diagnostic with explicit severity.
func diagnostic(
	code lint.Code,
	severity lint.Severity,
	path string,
	message string,
) Diagnostic {
	return Diagnostic{
		Code:     code,
		Severity: severity,
		Path:     path,
		Message:  message,
	}
}

// collectDiagnostics performs semantic checks and collects issues.
func collectDiagnostics(document *Document, options ValidateOptions) []Diagnostic {
	diagnostics := make([]Diagnostic, 0, 8)

	if document.RefSize.Width <= 0 {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateRefSizeWidthNonPositive,
				"ref_size.width",
				"must be > 0",
			),
		)
	}
	if document.RefSize.Height <= 0 {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateRefSizeHeightNonPositive,
				"ref_size.height",
				"must be > 0",
			),
		)
	}
	if document.RefSize.Width > 0 && !isPowerOfTwo(document.RefSize.Width) {
		diagnostics = append(
			diagnostics,
			warningDiagnostic(
				CodeValidateRefSizeNonPowerOfTwo,
				"ref_size.width",
				"should be a power of two",
			),
		)
	}
	if document.RefSize.Height > 0 && !isPowerOfTwo(document.RefSize.Height) {
		diagnostics = append(
			diagnostics,
			warningDiagnostic(
				CodeValidateRefSizeNonPowerOfTwo,
				"ref_size.height",
				"should be a power of two",
			),
		)
	}

	if len(document.Textures) == 0 {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateTexturesEmpty,
				"textures",
				"must contain at least one texture",
			),
		)
	}

	if len(document.Images) == 0 {
		diagnostics = append(
			diagnostics,
			warningDiagnostic(
				CodeValidateImagesEmpty,
				"images",
				"root images section is empty",
			),
		)
	}

	groupNames := make(map[string]int, len(document.Groups))
	globalImageNames := make(map[string]imageNameRef, len(document.Images))

	for textureIndex, texture := range document.Textures {
		if strings.TrimSpace(texture.Path) == "" {
			diagnostics = append(
				diagnostics,
				errorDiagnostic(
					CodeValidateTexturePathEmpty,
					textureFieldPath(textureIndex, "path"),
					"must be non-empty",
				),
			)
		}
		if texture.Mpix < 0 {
			diagnostics = append(
				diagnostics,
				errorDiagnostic(
					CodeValidateTextureMpixNegative,
					textureFieldPath(textureIndex, "mpix"),
					"must be >= 0",
				),
			)
		}
	}

	rootImageNames := make(map[string]imageNameRef, len(document.Images))
	for imageIndex, item := range document.Images {
		diagnostics = validateImage(
			diagnostics,
			item,
			-1,
			imageIndex,
			document.RefSize,
			rootImageNames,
			globalImageNames,
		)
	}

	for groupIndex, group := range document.Groups {
		groupName := strings.TrimSpace(group.Name)
		if groupName == "" {
			diagnostics = append(
				diagnostics,
				errorDiagnostic(
					CodeValidateGroupNameEmpty,
					groupFieldPath(groupIndex, "name"),
					"must be non-empty",
				),
			)
		} else if firstIndex, exists := groupNames[groupName]; exists {
			diagnostics = append(
				diagnostics,
				errorDiagnostic(
					CodeValidateGroupNameDuplicate,
					groupFieldPath(groupIndex, "name"),
					"duplicate name, first seen at "+
						groupFieldPath(firstIndex, "name"),
				),
			)
		} else {
			groupNames[groupName] = groupIndex
		}

		if len(group.Images) == 0 {
			diagnostics = append(
				diagnostics,
				warningDiagnostic(
					CodeValidateGroupImagesEmpty,
					groupFieldPath(groupIndex, "images"),
					"group images section is empty",
				),
			)
		}

		groupImageNames := make(map[string]imageNameRef, len(group.Images))
		for imageIndex, item := range group.Images {
			diagnostics = validateImage(
				diagnostics,
				item,
				groupIndex,
				imageIndex,
				document.RefSize,
				groupImageNames,
				globalImageNames,
			)
		}
	}

	diagnostics = append(
		diagnostics,
		collectImagePairDiagnostics(document, options)...,
	)

	return diagnostics
}

// collectImagePairDiagnostics emits overlap/padding diagnostics for image pairs.
func collectImagePairDiagnostics(
	document *Document,
	options ValidateOptions,
) []Diagnostic {
	images := collectImageRects(document)
	if len(images) < 2 {
		return nil
	}

	diagnostics := make([]Diagnostic, 0)

	for leftIndex := 0; leftIndex < len(images)-1; leftIndex++ {
		left := images[leftIndex]
		for rightIndex := leftIndex + 1; rightIndex < len(images); rightIndex++ {
			right := images[rightIndex]

			dx, dy := rectDistance(left, right)
			if dx == 0 && dy == 0 {
				diagnostics = append(diagnostics, warningDiagnostic(
					CodeValidateImageOverlap,
					imagePath(right.groupIndex, right.imageIndex),
					"overlaps with "+imagePath(left.groupIndex, left.imageIndex),
				))
				continue
			}

			if !options.EnablePaddingCheck {
				continue
			}

			if maxInt(dx, dy) >= options.MinPadding {
				continue
			}

			diagnostics = append(diagnostics, warningDiagnostic(
				CodeValidateImagePaddingTooSmall,
				imagePath(right.groupIndex, right.imageIndex),
				"padding to "+imagePath(left.groupIndex, left.imageIndex)+" is "+
					strconv.Itoa(maxInt(dx, dy))+
					", expected >= "+strconv.Itoa(options.MinPadding),
			))
		}
	}

	return diagnostics
}

// collectImageRects returns flattened image rectangles from root and groups.
func collectImageRects(document *Document) []imageRectRef {
	total := len(document.Images)
	for _, group := range document.Groups {
		total += len(group.Images)
	}

	out := make([]imageRectRef, 0, total)
	for imageIndex, item := range document.Images {
		out = append(out, imageRectRef{
			groupIndex: -1,
			imageIndex: imageIndex,
			x:          item.Pos.X,
			y:          item.Pos.Y,
			w:          item.Size.Width,
			h:          item.Size.Height,
		})
	}
	for groupIndex, group := range document.Groups {
		for imageIndex, item := range group.Images {
			out = append(out, imageRectRef{
				groupIndex: groupIndex,
				imageIndex: imageIndex,
				x:          item.Pos.X,
				y:          item.Pos.Y,
				w:          item.Size.Width,
				h:          item.Size.Height,
			})
		}
	}

	return out
}

// rectDistance returns horizontal and vertical distance between rectangles.
func rectDistance(left imageRectRef, right imageRectRef) (int, int) {
	leftRight := left.x + left.w
	rightRight := right.x + right.w
	leftBottom := left.y + left.h
	rightBottom := right.y + right.h

	dx := 0
	if leftRight < right.x {
		dx = right.x - leftRight
	} else if rightRight < left.x {
		dx = left.x - rightRight
	}

	dy := 0
	if leftBottom < right.y {
		dy = right.y - leftBottom
	} else if rightBottom < left.y {
		dy = left.y - rightBottom
	}

	return dx, dy
}

// maxInt returns larger value of two ints.
func maxInt(left int, right int) int {
	if left >= right {
		return left
	}

	return right
}

// validateImage validates one image and updates diagnostics.
func validateImage(
	diagnostics []Diagnostic,
	item Image,
	groupIndex, imageIndex int,
	refSize Size,
	scopeNames map[string]imageNameRef,
	globalNames map[string]imageNameRef,
) []Diagnostic {
	name := strings.TrimSpace(item.Name)
	currentNamePath := imageFieldPath(groupIndex, imageIndex, "name")

	if name == "" {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateImageNameEmpty,
				currentNamePath,
				"must be non-empty",
			),
		)
	} else {
		firstRef, exists := scopeNames[name]
		if exists {
			diagnostics = append(
				diagnostics,
				errorDiagnostic(
					CodeValidateImageNameDuplicate,
					currentNamePath,
					"duplicate name, first seen at "+
						imageFieldPath(
							firstRef.groupIndex,
							firstRef.imageIndex,
							"name",
						),
				),
			)
		} else {
			scopeNames[name] = imageNameRef{
				groupIndex: groupIndex,
				imageIndex: imageIndex,
			}
		}

		firstGlobalRef, globalExists := globalNames[name]
		if globalExists {
			diagnostics = append(
				diagnostics,
				infoDiagnostic(
					CodeValidateImageNameDuplicateGlobal,
					currentNamePath,
					"duplicate global name, first seen at "+
						imageFieldPath(
							firstGlobalRef.groupIndex,
							firstGlobalRef.imageIndex,
							"name",
						),
				),
			)
		} else {
			globalNames[name] = imageNameRef{
				groupIndex: groupIndex,
				imageIndex: imageIndex,
			}
		}
	}

	if item.Pos.X < 0 {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateImagePosXNegative,
				imageFieldPath(groupIndex, imageIndex, "pos.x"),
				"must be >= 0",
			),
		)
	}
	if item.Pos.Y < 0 {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateImagePosYNegative,
				imageFieldPath(groupIndex, imageIndex, "pos.y"),
				"must be >= 0",
			),
		)
	}

	if item.Size.Width <= 0 {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateImageWidthNonPositive,
				imageFieldPath(groupIndex, imageIndex, "size.width"),
				"must be > 0",
			),
		)
	}
	if item.Size.Height <= 0 {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateImageHeightNonPositive,
				imageFieldPath(groupIndex, imageIndex, "size.height"),
				"must be > 0",
			),
		)
	}

	unknownMask := item.Flags &^ knownImageFlagMask
	if item.Flags < 0 || unknownMask != 0 {
		message := "must use supported mask values (0..3)"
		if unknownMask != 0 {
			message += ", unknown bits: " + strconv.Itoa(int(unknownMask))
		}

		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateImageFlagsUnsupportedMask,
				imageFieldPath(groupIndex, imageIndex, "flags"),
				message,
			),
		)
	}

	if refSize.Width > 0 &&
		item.Size.Width > 0 &&
		item.Pos.X >= 0 &&
		item.Pos.X+item.Size.Width > refSize.Width {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateImageOutOfBoundsWidth,
				imagePath(groupIndex, imageIndex),
				"out of bounds by width against ref_size",
			),
		)
	}

	if refSize.Height > 0 &&
		item.Size.Height > 0 &&
		item.Pos.Y >= 0 &&
		item.Pos.Y+item.Size.Height > refSize.Height {
		diagnostics = append(
			diagnostics,
			errorDiagnostic(
				CodeValidateImageOutOfBoundsHeight,
				imagePath(groupIndex, imageIndex),
				"out of bounds by height against ref_size",
			),
		)
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

// isPowerOfTwo reports whether positive integer value is a power of two.
func isPowerOfTwo(value int) bool {
	return value > 0 && value&(value-1) == 0
}
