// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/woozymasta/lintkit/lint"
)

// fixtureCase describes one synthetic testdata file and expected structure.
type fixtureCase struct {
	name         string
	file         string
	wantRoot     int
	wantGroups   int
	wantGrouped  int
	wantTextures int
}

// fixtureCases lists synthetic parser fixtures in testdata.
var fixtureCases = []fixtureCase{
	{
		name:         "flat",
		file:         "flat.imageset",
		wantRoot:     64,
		wantGroups:   0,
		wantGrouped:  0,
		wantTextures: 1,
	},
	{
		name:         "groups",
		file:         "groups.imageset",
		wantRoot:     0,
		wantGroups:   4,
		wantGrouped:  96,
		wantTextures: 1,
	},
	{
		name:         "mixed",
		file:         "mixed.imageset",
		wantRoot:     20,
		wantGroups:   3,
		wantGrouped:  54,
		wantTextures: 2,
	},
}

// TestSyntheticFixturesParse validates shape of generated test fixtures.
func TestSyntheticFixturesParse(t *testing.T) {
	t.Parallel()

	for _, fixture := range fixtureCases {
		fixture := fixture
		t.Run(fixture.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join("testdata", fixture.file)
			doc, err := ParseFile(path)
			if err != nil {
				t.Fatalf("ParseFile(%s): %v", fixture.file, err)
			}

			if len(doc.Images) != fixture.wantRoot {
				t.Fatalf(
					"%s root images = %d, want %d",
					fixture.file,
					len(doc.Images),
					fixture.wantRoot,
				)
			}
			if len(doc.Groups) != fixture.wantGroups {
				t.Fatalf(
					"%s groups = %d, want %d",
					fixture.file,
					len(doc.Groups),
					fixture.wantGroups,
				)
			}

			grouped := 0
			for _, group := range doc.Groups {
				grouped += len(group.Images)
			}
			if grouped != fixture.wantGrouped {
				t.Fatalf(
					"%s grouped images = %d, want %d",
					fixture.file,
					grouped,
					fixture.wantGrouped,
				)
			}
			if len(doc.Textures) != fixture.wantTextures {
				t.Fatalf(
					"%s textures = %d, want %d",
					fixture.file,
					len(doc.Textures),
					fixture.wantTextures,
				)
			}

			if err := Validate(doc); err != nil {
				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("Validate(%s): %v", fixture.file, err)
				}

				if thresholdErr := lint.ErrorFromDiagnostics(
					validationErr.Diagnostics,
					lint.SeverityError,
				); thresholdErr != nil {
					t.Fatalf("Validate(%s): %v", fixture.file, thresholdErr)
				}
			}
		})
	}
}

// readFixtureBytes reads fixture file content from testdata.
func readFixtureBytes(tb testing.TB, file string) []byte {
	tb.Helper()

	path := filepath.Join("testdata", file)
	data, err := os.ReadFile(path)
	if err != nil {
		tb.Fatalf("read fixture %s: %v", path, err)
	}

	return data
}
