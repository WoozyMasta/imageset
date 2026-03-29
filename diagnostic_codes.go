// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"github.com/woozymasta/lintkit/lint"
)

const (
	// LintModule is stable lint module namespace for imageset rules.
	LintModule = "imageset"
)

const (
	// StageValidate marks semantic validation diagnostics.
	StageValidate lint.Stage = "validate"
)

const (
	// CodeValidateRefSizeWidthNonPositive reports non-positive atlas width.
	CodeValidateRefSizeWidthNonPositive lint.Code = 2001

	// CodeValidateRefSizeHeightNonPositive reports non-positive atlas height.
	CodeValidateRefSizeHeightNonPositive lint.Code = 2002

	// CodeValidateRefSizeNonPowerOfTwo reports non-power-of-two ref_size side.
	CodeValidateRefSizeNonPowerOfTwo lint.Code = 2003

	// CodeValidateTexturesEmpty reports missing Textures section entries.
	CodeValidateTexturesEmpty lint.Code = 2004

	// CodeValidateTexturePathEmpty reports empty texture path.
	CodeValidateTexturePathEmpty lint.Code = 2005

	// CodeValidateTextureMpixNegative reports negative texture mpix.
	CodeValidateTextureMpixNegative lint.Code = 2006

	// CodeValidateImagesEmpty reports missing root Images section entries.
	CodeValidateImagesEmpty lint.Code = 2007

	// CodeValidateGroupNameEmpty reports empty group name.
	CodeValidateGroupNameEmpty lint.Code = 2008

	// CodeValidateGroupNameDuplicate reports duplicate group name.
	CodeValidateGroupNameDuplicate lint.Code = 2009

	// CodeValidateGroupImagesEmpty reports empty images list inside group.
	CodeValidateGroupImagesEmpty lint.Code = 2010

	// CodeValidateImageNameEmpty reports empty image name.
	CodeValidateImageNameEmpty lint.Code = 2011

	// CodeValidateImageNameDuplicate reports duplicate image name.
	CodeValidateImageNameDuplicate lint.Code = 2012

	// CodeValidateImageNameDuplicateGlobal reports duplicate image name globally.
	CodeValidateImageNameDuplicateGlobal lint.Code = 2013

	// CodeValidateImagePosXNegative reports negative image x coordinate.
	CodeValidateImagePosXNegative lint.Code = 2014

	// CodeValidateImagePosYNegative reports negative image y coordinate.
	CodeValidateImagePosYNegative lint.Code = 2015

	// CodeValidateImageWidthNonPositive reports non-positive image width.
	CodeValidateImageWidthNonPositive lint.Code = 2016

	// CodeValidateImageHeightNonPositive reports non-positive image height.
	CodeValidateImageHeightNonPositive lint.Code = 2017

	// CodeValidateImageOutOfBoundsWidth reports atlas width overflow.
	CodeValidateImageOutOfBoundsWidth lint.Code = 2018

	// CodeValidateImageOutOfBoundsHeight reports atlas height overflow.
	CodeValidateImageOutOfBoundsHeight lint.Code = 2019

	// CodeValidateImageFlagsUnsupportedMask reports unsupported image flags mask.
	CodeValidateImageFlagsUnsupportedMask lint.Code = 2020

	// CodeValidateImageOverlap reports overlapping image rectangles.
	CodeValidateImageOverlap lint.Code = 2021

	// CodeValidateImagePaddingTooSmall reports too small gap between images.
	CodeValidateImagePaddingTooSmall lint.Code = 2022
)

var diagnosticCodeCatalogHandle = lint.NewCodeCatalogHandle(
	lint.CodeCatalogConfig{
		Module:            LintModule,
		CodePrefix:        "IMGSET",
		ModuleName:        "ImageSet",
		ModuleDescription: "Lint rules for image set atlas and entry validation.",
		ScopeDescriptions: map[lint.Stage]string{
			StageValidate: "Image set semantic validation diagnostics.",
		},
	},
	diagnosticCatalog,
)

// getDiagnosticCodeCatalog returns lazy-initialized diagnostics catalog.
func getDiagnosticCodeCatalog() (lint.CodeCatalog, error) {
	return diagnosticCodeCatalogHandle.Catalog()
}

// DiagnosticRuleSpec converts one diagnostic spec into lint rule metadata.
func DiagnosticRuleSpec(spec lint.CodeSpec) (lint.RuleSpec, error) {
	return diagnosticCodeCatalogHandle.RuleSpec(spec)
}

// LintRuleID returns lint rule ID mapped from stable imageset diagnostic code.
func LintRuleID(code lint.Code) string {
	return diagnosticCodeCatalogHandle.RuleIDOrUnknown(code)
}

// DiagnosticCatalog returns stable diagnostics metadata list.
func DiagnosticCatalog() []lint.CodeSpec {
	return diagnosticCodeCatalogHandle.CodeSpecs()
}

// DiagnosticByCode returns diagnostic metadata for code.
func DiagnosticByCode(code lint.Code) (lint.CodeSpec, bool) {
	return diagnosticCodeCatalogHandle.ByCode(code)
}

// LintRuleSpecs returns deterministic lint rule specs from diagnostics catalog.
func LintRuleSpecs() []lint.RuleSpec {
	return diagnosticCodeCatalogHandle.RuleSpecs()
}
