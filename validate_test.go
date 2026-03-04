// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"errors"
	"testing"
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
		Images: []Image{
			{
				Name: "",
				Pos:  Point{X: 33, Y: 33},
				Size: Size{Width: 40, Height: 40},
			},
			{
				Name: "dup",
				Pos:  Point{X: -1, Y: 0},
				Size: Size{Width: 1, Height: 1},
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
}
