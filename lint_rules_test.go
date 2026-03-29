// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"errors"
	"testing"

	"github.com/woozymasta/lintkit/lint"
	"github.com/woozymasta/lintkit/linttest"
)

type lintRuleTestRegistrar struct {
	rules []lint.RuleRunner
}

// Register stores runners in registrar test double.
func (registrar *lintRuleTestRegistrar) Register(
	runners ...lint.RuleRunner,
) error {
	registrar.rules = append(registrar.rules, runners...)

	return nil
}

func TestRegisterLintRulesNilRegistrar(t *testing.T) {
	t.Parallel()

	if err := RegisterLintRules(nil); !errors.Is(err, ErrNilLintRuleRegistrar) {
		t.Fatalf(
			"RegisterLintRules(nil) error=%v, want ErrNilLintRuleRegistrar",
			err,
		)
	}
}

func TestRegisterLintRules(t *testing.T) {
	t.Parallel()

	registrar := &lintRuleTestRegistrar{
		rules: make([]lint.RuleRunner, 0),
	}
	if err := RegisterLintRules(registrar); err != nil {
		t.Fatalf("RegisterLintRules() error: %v", err)
	}
	if len(registrar.rules) != len(DiagnosticCatalog()) {
		t.Fatalf(
			"registered rules=%d, want %d",
			len(registrar.rules),
			len(DiagnosticCatalog()),
		)
	}
}

func TestRegisterLintRulesByScope(t *testing.T) {
	t.Parallel()

	registrar := &lintRuleTestRegistrar{
		rules: make([]lint.RuleRunner, 0),
	}
	if err := RegisterLintRulesByScope(registrar, string(StageValidate)); err != nil {
		t.Fatalf("RegisterLintRulesByScope() error: %v", err)
	}
	if len(registrar.rules) != len(DiagnosticCatalog()) {
		t.Fatalf(
			"scope-registered rules=%d, want %d",
			len(registrar.rules),
			len(DiagnosticCatalog()),
		)
	}
}

func TestRegisterLintRulesByStage(t *testing.T) {
	t.Parallel()

	registrar := &lintRuleTestRegistrar{
		rules: make([]lint.RuleRunner, 0),
	}
	if err := RegisterLintRulesByStage(registrar, StageValidate); err != nil {
		t.Fatalf("RegisterLintRulesByStage() error: %v", err)
	}
	if len(registrar.rules) != len(DiagnosticCatalog()) {
		t.Fatalf(
			"stage-registered rules=%d, want %d",
			len(registrar.rules),
			len(DiagnosticCatalog()),
		)
	}
}

func TestLintRuleSpecsMatchCatalog(t *testing.T) {
	t.Parallel()

	linttest.AssertCatalogContract(
		t,
		LintModule,
		DiagnosticCatalog(),
		LintRuleSpecs(),
		LintRuleID,
	)
}

func TestAttachLintDiagnostics(t *testing.T) {
	t.Parallel()

	run := lint.RunContext{}
	diagnostics := []Diagnostic{
		errorDiagnostic(
			CodeValidateImageNameDuplicate,
			"images[1].name",
			"duplicate name",
		),
	}

	AttachLintDiagnostics(&run, diagnostics)

	grouped, ok := lint.GetIndexedByCode[Diagnostic, lint.Code](
		&run,
		lintRunValueByCodeKey,
	)
	if !ok {
		t.Fatal("GetIndexedByCode() ok=false, want true")
	}
	if len(grouped[CodeValidateImageNameDuplicate]) != 1 {
		t.Fatalf(
			"grouped[%q] len=%d, want 1",
			CodeValidateImageNameDuplicate,
			len(grouped[CodeValidateImageNameDuplicate]),
		)
	}
}

func TestDiagnosticLintDiagnostic(t *testing.T) {
	t.Parallel()

	diagnostic := Diagnostic{
		Code:     CodeValidateImageOutOfBoundsWidth,
		Severity: lint.SeverityError,
		Path:     "images[2]",
		Message:  "out of bounds by width against ref_size",
	}

	normalized := diagnostic.LintDiagnostic()
	if normalized.RuleID != LintRuleID(CodeValidateImageOutOfBoundsWidth) {
		t.Fatalf(
			"RuleID=%q, want %q",
			normalized.RuleID,
			LintRuleID(CodeValidateImageOutOfBoundsWidth),
		)
	}
	if normalized.Severity != lint.SeverityError {
		t.Fatalf("Severity=%q, want %q", normalized.Severity, lint.SeverityError)
	}
	if normalized.Path != diagnostic.Path {
		t.Fatalf("Path=%q, want %q", normalized.Path, diagnostic.Path)
	}
	if normalized.Message != diagnostic.Message {
		t.Fatalf("Message=%q, want %q", normalized.Message, diagnostic.Message)
	}
}
