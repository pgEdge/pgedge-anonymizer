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

// PassportGenerator generates passport numbers.
type PassportGenerator struct {
	BaseGenerator
}

// NewPassportGenerator creates a new passport generator.
func NewPassportGenerator() *PassportGenerator {
	return &PassportGenerator{
		BaseGenerator: BaseGenerator{name: "PASSPORT"},
	}
}

// Generate produces a passport number.
// Most passport numbers are 9 alphanumeric characters.
func (g *PassportGenerator) Generate(input string) string {
	// Detect if input has letters or is purely numeric
	hasLetters := false
	digitCount := 0
	for _, c := range input {
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
			hasLetters = true
		}
		if c >= '0' && c <= '9' {
			digitCount++
		}
	}

	// Default to 9 characters
	length := 9
	if digitCount > 9 {
		length = digitCount
	}

	if hasLetters {
		// Generate alphanumeric (letter + digits format like UK: AB123456)
		letters := "ABCDEFGHJKLMNPRSTUVWXYZ" // Excluding I, O, Q
		result := make([]byte, length)
		// First 1-2 characters are letters
		result[0] = letters[randomInt(len(letters))]
		if length > 8 {
			result[1] = letters[randomInt(len(letters))]
			for i := 2; i < length; i++ {
				result[i] = randomDigit()
			}
		} else {
			for i := 1; i < length; i++ {
				result[i] = randomDigit()
			}
		}
		return string(result)
	}

	// Pure numeric passport number
	return generateDigits(length)
}
