// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestWriteNilDocument(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	err := Write(&buffer, nil, nil)
	if !errors.Is(err, ErrNilDocument) {
		t.Fatalf("Write(nil) error = %v, want ErrNilDocument", err)
	}
}

func TestWriteSmoke(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:    "my ui",
		RefSize: Size{Width: 512, Height: 256},
		Textures: []Texture{
			{Mpix: 1, Path: "mod/data/ui.edds"},
		},
		Images: []Image{
			{
				Name: "root_icon",
				Pos:  Point{X: 1, Y: 2},
				Size: Size{Width: 3, Height: 4},
			},
		},
		Groups: []Group{
			{
				Name: "main hud",
				Images: []Image{
					{
						Name: "group_icon",
						Pos:  Point{X: 5, Y: 6},
						Size: Size{Width: 7, Height: 8},
					},
				},
			},
		},
	}

	var buffer bytes.Buffer
	if err := Write(&buffer, document, nil); err != nil {
		t.Fatalf("Write: %v", err)
	}

	output := buffer.String()
	contains := []string{
		"ImageSetClass {",
		`Name "my_ui"`,
		"RefSize 512 256",
		"Textures {",
		`path "mod/data/ui.edds"`,
		"Images {",
		`Name "root_icon"`,
		"Groups {",
		"ImageSetGroupClass main_hud {",
		`Name "main_hud"`,
	}
	for _, item := range contains {
		if !strings.Contains(output, item) {
			t.Fatalf("output does not contain %q\n%s", item, output)
		}
	}
}

func TestFormatAndRoundTrip(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:    "ui_icons",
		RefSize: Size{Width: 128, Height: 64},
		Images: []Image{
			{
				Name:  "ok",
				Pos:   Point{X: 1, Y: 2},
				Size:  Size{Width: 16, Height: 16},
				Flags: FlagHorizontalTile,
			},
		},
	}

	data, err := Format(document, &FormatOptions{
		UseCamelCaseNames: true,
		Indent:            "  ",
	})
	if err != nil {
		t.Fatalf("Format: %v", err)
	}

	parsed, err := ParseBytes(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if parsed.RefSize != document.RefSize {
		t.Fatalf("RefSize mismatch: got %+v want %+v", parsed.RefSize, document.RefSize)
	}
	if len(parsed.Images) != 1 || parsed.Images[0].Name != "Ok" {
		t.Fatalf("unexpected parsed images: %+v", parsed.Images)
	}
}

func TestFormatFlagsCanonicalStyle(t *testing.T) {
	t.Parallel()

	document := &Document{
		Name:    "flags",
		RefSize: Size{Width: 64, Height: 64},
		Images: []Image{
			{Name: "a", Pos: Point{X: 0, Y: 0}, Size: Size{Width: 8, Height: 8}, Flags: 0},
			{Name: "b", Pos: Point{X: 8, Y: 0}, Size: Size{Width: 8, Height: 8}, Flags: FlagHorizontalTile},
			{Name: "c", Pos: Point{X: 16, Y: 0}, Size: Size{Width: 8, Height: 8}, Flags: FlagVerticalTile},
			{Name: "d", Pos: Point{X: 24, Y: 0}, Size: Size{Width: 8, Height: 8}, Flags: FlagHorizontalTile | FlagVerticalTile},
		},
	}

	out, err := Format(document, nil)
	if err != nil {
		t.Fatalf("Format: %v", err)
	}

	text := string(out)
	if !strings.Contains(text, "Flags 0") {
		t.Fatalf("expected Flags 0 in output:\n%s", text)
	}
	if !strings.Contains(text, "Flags ISHorizontalTile") {
		t.Fatalf("expected Flags ISHorizontalTile in output:\n%s", text)
	}
	if !strings.Contains(text, "Flags ISVerticalTile") {
		t.Fatalf("expected Flags ISVerticalTile in output:\n%s", text)
	}
	if !strings.Contains(text, "Flags 3") {
		t.Fatalf("expected Flags 3 in output:\n%s", text)
	}
	if strings.Contains(text, "Flags ISHorizontalTile + ISVerticalTile") {
		t.Fatalf("unexpected symbolic combined flags in output:\n%s", text)
	}
	if strings.Contains(text, "Flags ISHorizontalTile | ISVerticalTile") {
		t.Fatalf("unexpected pipe combined flags in output:\n%s", text)
	}
}
