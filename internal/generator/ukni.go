/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package generator

import (
	"strings"
)

// UK National Insurance number format: AB 12 34 56 C
// Two prefix letters + 6 digits + 1 suffix letter
// Valid prefix letters exclude D, F, I, Q, U, V
// Valid suffix letters are A, B, C, D

// UKNIGenerator generates UK National Insurance numbers.
type UKNIGenerator struct {
	BaseGenerator
}

// NewUKNIGenerator creates a new UK NI generator.
func NewUKNIGenerator() *UKNIGenerator {
	return &UKNIGenerator{
		BaseGenerator: BaseGenerator{name: "UK_NI"},
	}
}

// Generate produces a UK National Insurance number.
func (g *UKNIGenerator) Generate(input string) string {
	// Valid prefix letters (excluding D, F, I, Q, U, V)
	prefixLetters := "ABCEGHJKLMNOPRSTWXYZ"
	// Valid suffix letters
	suffixLetters := "ABCD"

	// Generate two prefix letters
	prefix1 := prefixLetters[randomInt(len(prefixLetters))]
	prefix2 := prefixLetters[randomInt(len(prefixLetters))]

	// Generate 6 digits (3 pairs)
	digits := generateDigits(6)

	// Generate suffix letter
	suffix := suffixLetters[randomInt(len(suffixLetters))]

	// Detect format from input
	hasSpaces := strings.Contains(input, " ")

	if hasSpaces {
		return string(prefix1) + string(prefix2) + " " +
			digits[0:2] + " " + digits[2:4] + " " + digits[4:6] + " " +
			string(suffix)
	}

	return string(prefix1) + string(prefix2) + digits + string(suffix)
}
