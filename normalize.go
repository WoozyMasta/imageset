// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/imageset

package imageset

import "strings"

// NormalizeName converts a name to a safe imageset identifier.
//
// Output contains only ASCII letters, digits, and underscore.
// If useCamelCase is true, tokens are joined as CamelCase.
// Otherwise tokens are joined as snake_case.
func NormalizeName(input string, useCamelCase bool) string {
	if !useCamelCase && isLowerSnakeASCII(input) {
		return input
	}

	tokens := splitTokens(input)
	if len(tokens) == 0 {
		return ""
	}

	if useCamelCase {
		var builder strings.Builder
		for _, token := range tokens {
			builder.WriteString(toCamelToken(token))
		}

		return builder.String()
	}

	return strings.Join(tokens, "_")
}

// isLowerSnakeASCII reports whether input is already safe snake_case ASCII.
func isLowerSnakeASCII(input string) bool {
	if input == "" {
		return false
	}

	hasAlphaNum := false
	for _, r := range input {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			if r != '_' {
				hasAlphaNum = true
			}
			continue
		}

		return false
	}

	return hasAlphaNum
}

// splitTokens extracts ASCII alpha-numeric tokens from input.
func splitTokens(input string) []string {
	tokens := make([]string, 0, 8)
	buffer := make([]rune, 0, len(input))

	// flush emits one collected token.
	flush := func() {
		if len(buffer) == 0 {
			return
		}

		tokens = append(tokens, strings.ToLower(string(buffer)))
		buffer = buffer[:0]
	}

	for _, r := range input {
		if isASCIIAlphaNum(r) {
			buffer = append(buffer, r)
			continue
		}

		flush()
	}

	flush()

	return tokens
}

// toCamelToken converts one token to CamelCase form.
func toCamelToken(token string) string {
	if token == "" {
		return ""
	}

	var builder strings.Builder
	for i, r := range token {
		if i == 0 {
			builder.WriteRune(toUpperASCII(r))
			continue
		}

		builder.WriteRune(toLowerASCII(r))
	}

	return builder.String()
}

// isASCIIAlphaNum reports whether rune is ASCII letter or digit.
func isASCIIAlphaNum(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

// toLowerASCII converts ASCII uppercase letter to lowercase.
func toLowerASCII(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + ('a' - 'A')
	}

	return r
}

// toUpperASCII converts ASCII lowercase letter to uppercase.
func toUpperASCII(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - ('a' - 'A')
	}

	return r
}
