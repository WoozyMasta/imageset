// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"errors"
	"strings"
	"testing"
)

func TestParseClassName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line string
		want string
	}{
		{line: "ImageSetGroupClass GroupOne {", want: "GroupOne"},
		{line: "ImageSetGroupClass {", want: ""},
		{line: "ImageSetGroupClass    GroupTwo{", want: "GroupTwo"},
		{line: "", want: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.line, func(t *testing.T) {
			t.Parallel()

			if got := parseClassNameBytes([]byte(tc.line)); got != tc.want {
				t.Fatalf("parseClassNameBytes(%q) = %q, want %q", tc.line, got, tc.want)
			}
		})
	}
}

func TestParseRefSizeError(t *testing.T) {
	t.Parallel()

	content := "ImageSetClass {\n\tRefSize bad 1\n}\n"
	_, err := ParseString(content)
	if err == nil {
		t.Fatal("expected ParseString error for invalid RefSize")
	}

	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected ParseError, got %T", err)
	}
	if !strings.Contains(err.Error(), "invalid RefSize") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRootGroupsAndTextures(t *testing.T) {
	t.Parallel()

	content := `ImageSetClass {
	Name "ui"
	RefSize 256 256
	Textures {
		ImageSetTextureClass {
			mpix 1
			path "mod/data/ui.edds"
		}
	}
	Images {
		ImageSetDefClass RootIcon {
			Name "root_icon"
			Pos 1 2
			Size 3 4
			Flags 0
		}
	}
	Groups {
		ImageSetGroupClass HUD {
			Name "HUD"
			Images {
				ImageSetDefClass GroupIcon {
					Name "group_icon"
					Pos 10 20
					Size 30 40
					Flags ISHorizontalTile
				}
			}
		}
	}
}`

	doc, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString: %v", err)
	}

	if doc.Name != "ui" {
		t.Fatalf("name = %q, want %q", doc.Name, "ui")
	}
	if doc.RefSize != (Size{Width: 256, Height: 256}) {
		t.Fatalf("refsize = %+v, want 256x256", doc.RefSize)
	}
	if len(doc.Textures) != 1 {
		t.Fatalf("textures len = %d, want 1", len(doc.Textures))
	}
	if doc.Textures[0].Path != "mod/data/ui.edds" {
		t.Fatalf("texture path = %q", doc.Textures[0].Path)
	}
	if len(doc.Images) != 1 || doc.Images[0].Name != "root_icon" {
		t.Fatalf("unexpected root images: %+v", doc.Images)
	}
	if len(doc.Groups) != 1 {
		t.Fatalf("groups len = %d, want 1", len(doc.Groups))
	}
	if doc.Groups[0].Name != "HUD" {
		t.Fatalf("group name = %q, want HUD", doc.Groups[0].Name)
	}
	if len(doc.Groups[0].Images) != 1 || doc.Groups[0].Images[0].Name != "group_icon" {
		t.Fatalf("unexpected group images: %+v", doc.Groups[0].Images)
	}
	if doc.Groups[0].Images[0].Flags != FlagHorizontalTile {
		t.Fatalf("group image flags = %d, want %d", doc.Groups[0].Images[0].Flags, FlagHorizontalTile)
	}
}

func TestParseGroupWithoutImagesSection(t *testing.T) {
	t.Parallel()

	content := `ImageSetClass {
	Name "ui"
	RefSize 256 256
	Groups {
		ImageSetGroupClass HUD {
			Name "HUD"
			ImageSetDefClass GroupIcon {
				Name "group_icon"
				Pos 10 20
				Size 30 40
				Flags 0
			}
		}
	}
}`

	doc, err := ParseString(content)
	if err != nil {
		t.Fatalf("ParseString: %v", err)
	}
	if len(doc.Groups) != 1 {
		t.Fatalf("groups len = %d, want 1", len(doc.Groups))
	}
	if len(doc.Groups[0].Images) != 1 {
		t.Fatalf("group images len = %d, want 1", len(doc.Groups[0].Images))
	}
	if len(doc.Images) != 0 {
		t.Fatalf("root images len = %d, want 0", len(doc.Images))
	}
}

func TestParseUnclosedBlock(t *testing.T) {
	t.Parallel()

	_, err := ParseString("ImageSetClass {\n\tName \"ui\"\n")
	if err == nil {
		t.Fatal("expected unclosed block error")
	}

	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected ParseError, got %T", err)
	}
}

func TestParseMissingRoot(t *testing.T) {
	t.Parallel()

	_, err := ParseString(`Name "ui"`)
	if err == nil {
		t.Fatal("expected missing root error")
	}

	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected ParseError, got %T", err)
	}
}

func TestParseUnknownFieldFails(t *testing.T) {
	t.Parallel()

	_, err := ParseString(`ImageSetClass {
	Name "ui"
	RefSize 64 64
	UnknownField 1
}`)
	if err == nil {
		t.Fatal("expected unknown field parse error")
	}

	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected ParseError, got %T", err)
	}
	if !strings.Contains(parseErr.Error(), "unknown field") {
		t.Fatalf("unexpected error: %v", parseErr)
	}
}
