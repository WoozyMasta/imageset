// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"testing"

	"github.com/woozymasta/lintkit/lint"
)

func TestDiagnosticCatalogLookup(t *testing.T) {
	t.Parallel()

	spec, ok := DiagnosticByCode(CodeValidateImageNameDuplicate)
	if !ok {
		t.Fatalf(
			"DiagnosticByCode(%q) ok=false, want true",
			CodeValidateImageNameDuplicate,
		)
	}
	if spec.Code != CodeValidateImageNameDuplicate {
		t.Fatalf(
			"spec.Code=%q, want %q",
			spec.Code,
			CodeValidateImageNameDuplicate,
		)
	}
}

func TestLintRuleID(t *testing.T) {
	t.Parallel()

	spec, ok := DiagnosticByCode(CodeValidateImageNameDuplicate)
	if !ok {
		t.Fatalf(
			"DiagnosticByCode(%q) ok=false, want true",
			CodeValidateImageNameDuplicate,
		)
	}
	want := lint.BuildRuleID(
		LintModule,
		spec.Stage,
		spec.Message,
		CodeValidateImageNameDuplicate,
	)
	if got := LintRuleID(CodeValidateImageNameDuplicate); got != want {
		t.Fatalf("LintRuleID()=%q, want %q", got, want)
	}
	if got := LintRuleID(0); got != "imageset.unknown" {
		t.Fatalf("LintRuleID(empty)=%q, want %q", got, "imageset.unknown")
	}
}

func TestDiagnosticRuleSpec(t *testing.T) {
	t.Parallel()

	spec, ok := DiagnosticByCode(CodeValidateTexturesEmpty)
	if !ok {
		t.Fatalf("DiagnosticByCode(%q) ok=false, want true", CodeValidateTexturesEmpty)
	}

	ruleSpec, err := DiagnosticRuleSpec(spec)
	if err != nil {
		t.Fatalf("DiagnosticRuleSpec() error: %v", err)
	}

	if ruleSpec.Code != lint.ApplyCodePrefix("IMGSET", CodeValidateTexturesEmpty) {
		t.Fatalf("rule code=%q is unexpected", ruleSpec.Code)
	}
}
