// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

/*
Package imageset provides parser, formatter, and validation helpers for
DayZ .imageset files.

The package centers around Document:

  - Parse/ParseBytes/ParseFile read .imageset text into Document
  - Write/Format serialize Document back to canonical text form
  - Validate checks semantic constraints with stable diagnostic codes
  - ValidateWithOptions enables optional checks such as padding

Common flow:

	doc, err := imageset.ParseFile("ui.imageset")
	if err != nil {
		// handle parse error
	}
	if err := imageset.Validate(doc); err != nil {
		// handle validation diagnostics
	}

lintkit integration is included to expose stable rule metadata and register
imageset diagnostics in shared lint pipelines.
*/
package imageset
