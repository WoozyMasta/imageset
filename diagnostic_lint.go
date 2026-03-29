// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import "github.com/woozymasta/lintkit/lint"

// LintDiagnostic converts one imageset diagnostic into shared lint model.
func (diagnostic Diagnostic) LintDiagnostic() lint.Diagnostic {
	return lint.Diagnostic{
		RuleID:   LintRuleID(diagnostic.Code),
		Severity: diagnostic.Severity,
		Message:  diagnostic.Message,
		Path:     diagnostic.Path,
	}
}
