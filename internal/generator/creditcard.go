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

// CreditCardGenerator generates credit card numbers.
type CreditCardGenerator struct {
	BaseGenerator
}

// NewCreditCardGenerator creates a new credit card generator.
func NewCreditCardGenerator() *CreditCardGenerator {
	return &CreditCardGenerator{
		BaseGenerator: BaseGenerator{name: "CREDIT_CARD"},
	}
}

// Generate produces a credit card number with valid Luhn check digit.
func (g *CreditCardGenerator) Generate(input string) string {
	// Detect card type from first digit(s) of input
	// Visa: 4, MC: 51-55, Amex: 34/37, Discover: 6011
	prefix := "4" // Default to Visa format

	// Detect separator format
	var sep string
	if strings.Contains(input, "-") {
		sep = "-"
	} else if strings.Contains(input, " ") {
		sep = " "
	}

	// Generate 15 digits (16th will be check digit)
	digits := prefix + generateDigits(14)

	// Calculate and append Luhn check digit
	checkDigit := luhnCheckDigit(digits)
	digits += string(checkDigit)

	// Format with separators if detected
	if sep != "" {
		return fmt.Sprintf("%s%s%s%s%s%s%s",
			digits[0:4], sep, digits[4:8], sep,
			digits[8:12], sep, digits[12:16])
	}

	return digits
}

// CreditCardExpiryGenerator generates credit card expiry dates.
type CreditCardExpiryGenerator struct {
	BaseGenerator
}

// NewCreditCardExpiryGenerator creates a new expiry date generator.
func NewCreditCardExpiryGenerator() *CreditCardExpiryGenerator {
	return &CreditCardExpiryGenerator{
		BaseGenerator: BaseGenerator{name: "CREDIT_CARD_EXPIRY"},
	}
}

// Generate produces a credit card expiry date.
func (g *CreditCardExpiryGenerator) Generate(input string) string {
	// Generate month (01-12)
	month := 1 + randomInt(12)

	// Generate year (current + 1-5 years, using 25-30 for simplicity)
	year := 25 + randomInt(6)

	// Detect format (MM/YY vs MM/YYYY)
	if len(input) >= 7 && strings.Contains(input, "/") {
		// MM/YYYY format
		return fmt.Sprintf("%02d/20%02d", month, year)
	}

	// MM/YY format (default)
	return fmt.Sprintf("%02d/%02d", month, year)
}

// CreditCardCVVGenerator generates credit card CVV numbers.
type CreditCardCVVGenerator struct {
	BaseGenerator
}

// NewCreditCardCVVGenerator creates a new CVV generator.
func NewCreditCardCVVGenerator() *CreditCardCVVGenerator {
	return &CreditCardCVVGenerator{
		BaseGenerator: BaseGenerator{name: "CREDIT_CARD_CVV"},
	}
}

// Generate produces a CVV number.
func (g *CreditCardCVVGenerator) Generate(input string) string {
	// Detect length (3 for most cards, 4 for Amex)
	length := 3
	inputDigits := 0
	for _, c := range input {
		if c >= '0' && c <= '9' {
			inputDigits++
		}
	}
	if inputDigits == 4 {
		length = 4
	}

	return generateDigits(length)
}
