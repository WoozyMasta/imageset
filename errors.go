// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"errors"
	"fmt"

	"github.com/woozymasta/lintkit/lint"
)

var (
	// ErrInvalidSyntax means input text does not match expected .imageset syntax.
	ErrInvalidSyntax = errors.New("imageset: invalid syntax")

	// ErrNilDocument means API was called with nil *Document.
	ErrNilDocument = errors.New("imageset: nil document")

	// ErrUnknownFlag means a flag token is not recognized.
	ErrUnknownFlag = errors.New("imageset: unknown flag")

	// ErrNilLintRuleRegistrar indicates nil lint rule registrar in registration.
	ErrNilLintRuleRegistrar = lint.ErrNilRuleRegistrar
)

// ParseError reports parse failure with line context.
type ParseError struct {
	Cause   error // Underlying parse error.
	Message string
	Line    int // 1-based line number.
}

// Error formats the parse error.
func (e *ParseError) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch {
	case e.Line > 0 && e.Message != "":
		return fmt.Sprintf("line %d: %s", e.Line, e.Message)
	case e.Line > 0 && e.Cause != nil:
		return fmt.Sprintf("line %d: %v", e.Line, e.Cause)
	case e.Message != "":
		return e.Message
	case e.Cause != nil:
		return e.Cause.Error()
	default:
		return ErrInvalidSyntax.Error()
	}
}

// Unwrap returns the underlying parse cause.
func (e *ParseError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}
