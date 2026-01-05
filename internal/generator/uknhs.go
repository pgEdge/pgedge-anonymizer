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
	"fmt"
	"strings"
)

// UK NHS numbers are 10 digits with a modulus 11 check digit.
// Format is typically XXX XXX XXXX

// UKNHSGenerator generates UK NHS numbers.
type UKNHSGenerator struct {
	BaseGenerator
}

// NewUKNHSGenerator creates a new UK NHS generator.
func NewUKNHSGenerator() *UKNHSGenerator {
	return &UKNHSGenerator{
		BaseGenerator: BaseGenerator{name: "UK_NHS"},
	}
}

// Generate produces a UK NHS number with valid check digit.
func (g *UKNHSGenerator) Generate(input string) string {
	// Generate first 9 digits
	digits := make([]int, 10)
	for i := 0; i < 9; i++ {
		digits[i] = randomInt(10)
	}

	// Calculate modulus 11 check digit
	// Weights: 10, 9, 8, 7, 6, 5, 4, 3, 2
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	remainder := sum % 11
	checkDigit := 11 - remainder
	if checkDigit == 11 {
		checkDigit = 0
	}
	if checkDigit == 10 {
		// Invalid check digit, regenerate first digit
		digits[0] = (digits[0] + 1) % 10
		// Recalculate
		sum = 0
		for i := 0; i < 9; i++ {
			sum += digits[i] * (10 - i)
		}
		remainder = sum % 11
		checkDigit = 11 - remainder
		if checkDigit == 11 {
			checkDigit = 0
		}
	}
	digits[9] = checkDigit

	// Format as string
	result := ""
	for _, d := range digits {
		result += fmt.Sprintf("%d", d)
	}

	// Detect format from input
	hasSpaces := strings.Contains(input, " ")
	if hasSpaces {
		return result[0:3] + " " + result[3:6] + " " + result[6:10]
	}

	return result
}
