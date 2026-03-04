// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import "testing"

func TestNormalizeName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		useCamelCase bool
		expect       string
	}{
		{
			name:         "snake-basic",
			input:        "My Icon Name",
			useCamelCase: false,
			expect:       "my_icon_name",
		},
		{
			name:         "snake-symbols",
			input:        "UI@@Button##OK",
			useCamelCase: false,
			expect:       "ui_button_ok",
		},
		{
			name:         "camel-basic",
			input:        "my icon name",
			useCamelCase: true,
			expect:       "MyIconName",
		},
		{
			name:         "camel-mixed",
			input:        "HP_Bar-42",
			useCamelCase: true,
			expect:       "HpBar42",
		},
		{
			name:         "empty",
			input:        "___",
			useCamelCase: false,
			expect:       "",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := NormalizeName(tc.input, tc.useCamelCase)
			if got != tc.expect {
				t.Fatalf(
					"NormalizeName(%q, %v) = %q, want %q",
					tc.input,
					tc.useCamelCase,
					got,
					tc.expect,
				)
			}
		})
	}
}
