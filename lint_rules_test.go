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
	diagnostics := []lint.Diagnostic{
		diagnostic(
			CodeValidateImageNameDuplicate,
			lint.SeverityError,
			"images[1].name",
			"duplicate name",
		),
	}

	AttachLintDiagnostics(&run, diagnostics)

	grouped, ok := lint.GetIndexedByCode[lint.Diagnostic, lint.Code](
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

func TestLintDiagnosticCode(t *testing.T) {
	t.Parallel()

	item := lint.Diagnostic{Code: publicLintCode(CodeValidateImageOutOfBoundsWidth)}
	got := lintDiagnosticCode(item)
	if got != CodeValidateImageOutOfBoundsWidth {
		t.Fatalf(
			"lintDiagnosticCode()=%d, want %d",
			got,
			CodeValidateImageOutOfBoundsWidth,
		)
	}

	if got := lintDiagnosticCode(lint.Diagnostic{Code: ""}); got != 0 {
		t.Fatalf("lintDiagnosticCode(empty)=%d, want 0", got)
	}
}
