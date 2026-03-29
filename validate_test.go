// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"errors"
	"testing"

	"github.com/woozymasta/lintkit/lint"
)

func TestValidateNilDocument(t *testing.T) {
	t.Parallel()

	err := Validate(nil)
	if !errors.Is(err, ErrNilDocument) {
		t.Fatalf("Validate(nil) error = %v, want ErrNilDocument", err)
	}
}

func TestValidateSuccess(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:    "ui",
		RefSize: Size{Width: 128, Height: 128},
		Textures: []Texture{
			{Path: "mod/data/ui.edds", Mpix: 1},
		},
		Images: []Image{
			{
				Name: "ok",
				Pos:  Point{X: 1, Y: 2},
				Size: Size{Width: 32, Height: 32},
			},
		},
	}

	if err := Validate(document); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateDiagnostics(t *testing.T) {
	t.Parallel()

	document := &Document{
		RefSize: Size{Width: 64, Height: 64},
		Textures: []Texture{
			{Path: "mod/data/ui.edds", Mpix: 1},
		},
		Images: []Image{
			{
				Name: "",
				Pos:  Point{X: 33, Y: 33},
				Size: Size{Width: 40, Height: 40},
			},
			{
				Name:  "dup",
				Pos:   Point{X: -1, Y: 0},
				Size:  Size{Width: 1, Height: 1},
				Flags: -1,
			},
			{
				Name: "dup",
				Pos:  Point{X: 0, Y: 0},
				Size: Size{Width: 1, Height: 0},
			},
		},
		Groups: []Group{
			{
				Name: "",
				Images: []Image{
					{
						Name: "dup",
						Pos:  Point{X: 0, Y: 0},
						Size: Size{Width: 1, Height: 1},
					},
				},
			},
		},
	}

	err := Validate(document)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if len(validationErr.Diagnostics) == 0 {
		t.Fatal("expected diagnostics")
	}

	var (
		hasError bool
		hasWarn  bool
		hasInfo  bool
	)

	for index := range validationErr.Diagnostics {
		diagnostic := validationErr.Diagnostics[index]
		if diagnostic.Code == 0 {
			t.Fatalf("Diagnostics[%d].Code is empty", index)
		}

		switch diagnostic.Severity {
		case lint.SeverityError:
			hasError = true
		case lint.SeverityWarning:
			hasWarn = true
		case lint.SeverityInfo:
			hasInfo = true
		default:
			t.Fatalf("Diagnostics[%d].Severity=%q is unsupported", index, diagnostic.Severity)
		}
	}

	if !hasError {
		t.Fatal("expected at least one error diagnostic")
	}
	if !hasWarn {
		t.Fatal("expected at least one warning diagnostic")
	}
	if !hasInfo {
		t.Fatal("expected at least one info diagnostic")
	}
}

func TestValidateEmptyTexturesAndImages(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:    "ui",
		RefSize: Size{Width: 128, Height: 128},
	}

	err := Validate(document)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	hasTexturesError := false
	hasImagesWarning := false
	for _, diagnostic := range validationErr.Diagnostics {
		if diagnostic.Code == CodeValidateTexturesEmpty &&
			diagnostic.Severity == lint.SeverityError {
			hasTexturesError = true
		}

		if diagnostic.Code == CodeValidateImagesEmpty &&
			diagnostic.Severity == lint.SeverityWarning {
			hasImagesWarning = true
		}
	}

	if !hasTexturesError {
		t.Fatal("expected CodeValidateTexturesEmpty error diagnostic")
	}
	if !hasImagesWarning {
		t.Fatal("expected CodeValidateImagesEmpty warning diagnostic")
	}
}

func TestValidateDuplicateGroupName(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:     "ui",
		RefSize:  Size{Width: 128, Height: 128},
		Textures: []Texture{{Path: "mod/data/ui.edds", Mpix: 1}},
		Images: []Image{
			{Name: "root", Pos: Point{X: 0, Y: 0}, Size: Size{Width: 8, Height: 8}},
		},
		Groups: []Group{
			{Name: "HUD"},
			{Name: "HUD"},
		},
	}

	err := Validate(document)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	for _, diagnostic := range validationErr.Diagnostics {
		if diagnostic.Code == CodeValidateGroupNameDuplicate &&
			diagnostic.Severity == lint.SeverityError {
			return
		}
	}

	t.Fatal("expected CodeValidateGroupNameDuplicate error diagnostic")
}

func TestValidateRefSizeNonPowerOfTwoWarning(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:     "ui",
		RefSize:  Size{Width: 300, Height: 512},
		Textures: []Texture{{Path: "mod/data/ui.edds", Mpix: 1}},
		Images: []Image{
			{Name: "root", Pos: Point{X: 0, Y: 0}, Size: Size{Width: 8, Height: 8}},
		},
	}

	err := Validate(document)
	if err == nil {
		t.Fatal("expected validation warning")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	for _, diagnostic := range validationErr.Diagnostics {
		if diagnostic.Code == CodeValidateRefSizeNonPowerOfTwo &&
			diagnostic.Severity == lint.SeverityWarning &&
			diagnostic.Path == "ref_size.width" {
			return
		}
	}

	t.Fatal("expected CodeValidateRefSizeNonPowerOfTwo warning for ref_size.width")
}

func TestValidateUnsupportedFlagsMaskError(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:     "ui",
		RefSize:  Size{Width: 512, Height: 512},
		Textures: []Texture{{Path: "mod/data/ui.edds", Mpix: 1}},
		Images: []Image{
			{Name: "root", Pos: Point{X: 0, Y: 0}, Size: Size{Width: 8, Height: 8}, Flags: 4},
		},
	}

	err := Validate(document)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	for _, diagnostic := range validationErr.Diagnostics {
		if diagnostic.Code == CodeValidateImageFlagsUnsupportedMask &&
			diagnostic.Severity == lint.SeverityError &&
			diagnostic.Path == "images[0].flags" {
			return
		}
	}

	t.Fatal("expected CodeValidateImageFlagsUnsupportedMask error")
}

func TestValidateGroupImagesEmptyWarning(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:     "ui",
		RefSize:  Size{Width: 512, Height: 512},
		Textures: []Texture{{Path: "mod/data/ui.edds", Mpix: 1}},
		Images: []Image{
			{Name: "root", Pos: Point{X: 0, Y: 0}, Size: Size{Width: 8, Height: 8}},
		},
		Groups: []Group{
			{Name: "Empty"},
		},
	}

	err := Validate(document)
	if err == nil {
		t.Fatal("expected validation warning")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	for _, diagnostic := range validationErr.Diagnostics {
		if diagnostic.Code == CodeValidateGroupImagesEmpty &&
			diagnostic.Severity == lint.SeverityWarning &&
			diagnostic.Path == "groups[0].images" {
			return
		}
	}

	t.Fatal("expected CodeValidateGroupImagesEmpty warning")
}

func TestValidateImageOverlapWarning(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:     "ui",
		RefSize:  Size{Width: 512, Height: 512},
		Textures: []Texture{{Path: "mod/data/ui.edds", Mpix: 1}},
		Images: []Image{
			{Name: "a", Pos: Point{X: 0, Y: 0}, Size: Size{Width: 16, Height: 16}},
			{Name: "b", Pos: Point{X: 8, Y: 8}, Size: Size{Width: 16, Height: 16}},
		},
	}

	err := Validate(document)
	if err == nil {
		t.Fatal("expected validation warning")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	for _, diagnostic := range validationErr.Diagnostics {
		if diagnostic.Code == CodeValidateImageOverlap &&
			diagnostic.Severity == lint.SeverityWarning {
			return
		}
	}

	t.Fatal("expected CodeValidateImageOverlap warning")
}

func TestValidateImagePaddingWarningDisabledByDefault(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:     "ui",
		RefSize:  Size{Width: 512, Height: 512},
		Textures: []Texture{{Path: "mod/data/ui.edds", Mpix: 1}},
		Images: []Image{
			{Name: "a", Pos: Point{X: 0, Y: 0}, Size: Size{Width: 16, Height: 16}},
			{Name: "b", Pos: Point{X: 18, Y: 0}, Size: Size{Width: 16, Height: 16}},
		},
	}

	err := Validate(document)
	if err != nil {
		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected ValidationError, got %T", err)
		}

		for _, diagnostic := range validationErr.Diagnostics {
			if diagnostic.Code == CodeValidateImagePaddingTooSmall {
				t.Fatal("did not expect CodeValidateImagePaddingTooSmall by default")
			}
		}
	}
}

func TestValidateImagePaddingWarningEnabledWithOptions(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:     "ui",
		RefSize:  Size{Width: 512, Height: 512},
		Textures: []Texture{{Path: "mod/data/ui.edds", Mpix: 1}},
		Images: []Image{
			{Name: "a", Pos: Point{X: 0, Y: 0}, Size: Size{Width: 16, Height: 16}},
			{Name: "b", Pos: Point{X: 18, Y: 0}, Size: Size{Width: 16, Height: 16}},
		},
	}

	err := ValidateWithOptions(document, &ValidateOptions{
		EnablePaddingCheck: true,
		MinPadding:         4,
	})
	if err == nil {
		t.Fatal("expected validation warning")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	for _, diagnostic := range validationErr.Diagnostics {
		if diagnostic.Code == CodeValidateImagePaddingTooSmall &&
			diagnostic.Severity == lint.SeverityWarning {
			return
		}
	}

	t.Fatal("expected CodeValidateImagePaddingTooSmall warning")
}
