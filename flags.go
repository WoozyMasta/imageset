// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// FlagHorizontalTile corresponds to ISHorizontalTile.
	FlagHorizontalTile Flags = 1

	// FlagVerticalTile corresponds to ISVerticalTile.
	FlagVerticalTile Flags = 2
)

// Flags is a DayZ .imageset bitset.
type Flags int

// ParseFlagsExpr parses numeric or symbolic flags expression.
func ParseFlagsExpr(expr string) (Flags, error) {
	return parseFlagsExprString(expr)
}

// Has reports whether a specific flag bit is set.
func (f Flags) Has(flag Flags) bool {
	return f&flag != 0
}

// String returns a stable symbolic representation.
func (f Flags) String() string {
	if f == 0 {
		return "0"
	}

	names := make([]string, 0, len(flagNames))
	for _, item := range flagNames {
		if f.Has(item.flag) {
			names = append(names, item.name)
		}
	}

	if len(names) == 0 {
		return strconv.Itoa(int(f))
	}

	return strings.Join(names, " + ")
}

// flagNames stores canonical symbolic names for known flags.
var flagNames = []struct {
	name string
	flag Flags
}{
	{flag: FlagHorizontalTile, name: "ISHorizontalTile"},
	{flag: FlagVerticalTile, name: "ISVerticalTile"},
}

// parseFlagsExprString parses numeric or symbolic flags expression string.
func parseFlagsExprString(expr string) (Flags, error) {
	return parseFlagsExprBytes([]byte(expr))
}

// parseFlagsExprBytes parses numeric or symbolic flags expression bytes.
func parseFlagsExprBytes(expr []byte) (Flags, error) {
	expr = trimASCIISpace(expr)
	if len(expr) == 0 {
		return 0, nil
	}

	var (
		out      Flags
		hasToken bool
		index    int
	)

	for index < len(expr) {
		for index < len(expr) && isFlagsSeparator(expr[index]) {
			index++
		}

		if index >= len(expr) {
			break
		}

		start := index
		for index < len(expr) && !isFlagsSeparator(expr[index]) {
			index++
		}

		token := trimASCIISpace(expr[start:index])
		if len(token) == 0 {
			continue
		}

		hasToken = true
		if value, ok := parseNamedFlagToken(token); ok {
			out |= value
			continue
		}

		value, ok := parseIntTokenFromBytes(token)
		if ok {
			out |= Flags(value)
			continue
		}

		return 0, fmt.Errorf("%w: %q", ErrUnknownFlag, token)
	}

	if !hasToken {
		return 0, nil
	}

	return out, nil
}

// parseNamedFlagToken resolves symbolic flags without temporary strings.
func parseNamedFlagToken(token []byte) (Flags, bool) {
	switch {
	case equalToken(token, "ISHorizontalTile"):
		return FlagHorizontalTile, true
	case equalToken(token, "ISVerticalTile"):
		return FlagVerticalTile, true
	default:
		return 0, false
	}
}

// parseIntTokenFromBytes parses signed decimal integer token bytes.
func parseIntTokenFromBytes(raw []byte) (int, bool) {
	value, err := parseIntTokenBytes(raw)
	if err != nil {
		return 0, false
	}

	return value, true
}

// isFlagsSeparator reports token separators for flags expressions.
func isFlagsSeparator(ch byte) bool {
	return isASCIISpace(ch) || ch == '|' || ch == '+'
}
