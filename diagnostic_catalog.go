// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import "github.com/woozymasta/lintkit/lint"

// diagnosticCatalog stores stable diagnostics metadata table.
var diagnosticCatalog = []lint.CodeSpec{
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateRefSizeWidthNonPositive,
			StageValidate,
			"`ref_size.width` must be greater than 0",
		),
		"Root atlas width is required and must be a positive integer.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateRefSizeHeightNonPositive,
			StageValidate,
			"`ref_size.height` must be greater than 0",
		),
		"Root atlas height is required and must be a positive integer.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeValidateRefSizeNonPowerOfTwo,
			StageValidate,
			"`ref_size` side should be a power of two",
		),
		"Power-of-two size improves compatibility with common atlas workflows.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateTexturesEmpty,
			StageValidate,
			"`textures` section must contain at least one texture",
		),
		"DayZ expects texture list to exist and contain at least one entry.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateTexturePathEmpty,
			StageValidate,
			"`texture.path` must be non-empty",
		),
		"Each texture entry must define a path to texture source file.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateTextureMpixNegative,
			StageValidate,
			"`texture.mpix` must be zero or greater",
		),
		"Negative pixels-per-meter value is invalid for texture metadata.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeValidateImagesEmpty,
			StageValidate,
			"`images` section is empty",
		),
		"Currently warning-only to allow group-only files during migration.",
	),
	lint.ErrorCodeSpec(
		CodeValidateGroupNameEmpty,
		StageValidate,
		"`group` name must be non-empty",
	),
	lint.ErrorCodeSpec(
		CodeValidateGroupNameDuplicate,
		StageValidate,
		"`group` name must be unique",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeValidateGroupImagesEmpty,
			StageValidate,
			"`group.images` section is empty",
		),
		"Group exists but does not define any image entries.",
	),
	lint.ErrorCodeSpec(
		CodeValidateImageNameEmpty,
		StageValidate,
		"`image` name must be non-empty",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateImageNameDuplicate,
			StageValidate,
			"`image` name must be unique in current section",
		),
		"Duplicates are checked separately for root images and for each group.",
	),
	withDescription(
		lint.InfoCodeSpec(
			CodeValidateImageNameDuplicateGlobal,
			StageValidate,
			"`image` name is duplicated globally across root/groups",
		),
		"Informational only: global duplicate may still resolve by group prefix.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateImagePosXNegative,
			StageValidate,
			"`image.pos.x` must be zero or greater",
		),
		"Negative coordinates place sprite origin outside atlas bounds.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateImagePosYNegative,
			StageValidate,
			"`image.pos.y` must be zero or greater",
		),
		"Negative coordinates place sprite origin outside atlas bounds.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateImageWidthNonPositive,
			StageValidate,
			"`image.size.width` must be greater than 0",
		),
		"Sprite width must be positive for usable atlas entry.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateImageHeightNonPositive,
			StageValidate,
			"`image.size.height` must be greater than 0",
		),
		"Sprite height must be positive for usable atlas entry.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateImageOutOfBoundsWidth,
			StageValidate,
			"image exceeds atlas bounds on width",
		),
		"Triggered when pos.x + size.width is greater than ref_size.width.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateImageOutOfBoundsHeight,
			StageValidate,
			"image exceeds atlas bounds on height",
		),
		"Triggered when pos.y + size.height is greater than ref_size.height.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeValidateImageFlagsUnsupportedMask,
			StageValidate,
			"`image.flags` must use supported mask values (0..3)",
		),
		"Allowed flags are 0, ISHorizontalTile(1), ISVerticalTile(2), and 3.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeValidateImageOverlap,
			StageValidate,
			"image rectangle overlaps another image",
		),
		"Overlap means two sprite rectangles intersect in atlas space.",
	),
	withDescription(
		lint.WithCodeOptions(
			lint.WithCodeEnabled(
				lint.WarningCodeSpec(
					CodeValidateImagePaddingTooSmall,
					StageValidate,
					"image padding is smaller than configured minimum",
				),
				false,
			),
			ImagePaddingRuleOptions{
				MinPadding: 4,
			},
		),
		"Disabled by default. Enable to enforce minimum gap between image rectangles.",
	),
}

// withDescription attaches detailed docs text to code spec row.
func withDescription(spec lint.CodeSpec, description string) lint.CodeSpec {
	spec.Description = description
	return spec
}
