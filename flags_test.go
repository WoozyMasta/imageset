// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"errors"
	"testing"
)

func TestParseFlagsExpr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		expr    string
		want    Flags
		wantErr bool
	}{
		{name: "empty", expr: "", want: 0},
		{name: "numeric", expr: "3", want: Flags(3)},
		{name: "named-single", expr: "ISHorizontalTile", want: FlagHorizontalTile},
		{
			name: "named-combined-plus",
			expr: "ISHorizontalTile + ISVerticalTile",
			want: FlagHorizontalTile | FlagVerticalTile,
		},
		{
			name: "named-combined-pipe",
			expr: "ISHorizontalTile | ISVerticalTile",
			want: FlagHorizontalTile | FlagVerticalTile,
		},
		{name: "invalid", expr: "UNKNOWN_FLAG", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseFlagsExpr(tc.expr)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseFlagsExpr(%q): expected error", tc.expr)
				}
				if !errors.Is(err, ErrUnknownFlag) {
					t.Fatalf("ParseFlagsExpr(%q): wrong error %v", tc.expr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseFlagsExpr(%q): %v", tc.expr, err)
			}
			if got != tc.want {
				t.Fatalf("ParseFlagsExpr(%q) = %d, want %d", tc.expr, got, tc.want)
			}
		})
	}
}

func TestFlagsString(t *testing.T) {
	t.Parallel()

	if got := Flags(0).String(); got != "0" {
		t.Fatalf("Flags(0).String() = %q, want 0", got)
	}

	if got := (FlagHorizontalTile | FlagVerticalTile).String(); got != "ISHorizontalTile + ISVerticalTile" {
		t.Fatalf("combined String() = %q", got)
	}
}
