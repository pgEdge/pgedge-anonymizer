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

// USPhoneGenerator generates US phone numbers.
type USPhoneGenerator struct {
	BaseGenerator
}

// NewUSPhoneGenerator creates a new US phone generator.
func NewUSPhoneGenerator() *USPhoneGenerator {
	return &USPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "US_PHONE"},
	}
}

// Generate produces a US phone number preserving the input format.
// Uses 555 exchange which is reserved for fictional use in North America.
func (g *USPhoneGenerator) Generate(input string) string {
	format := detectPhoneFormat(input)

	// Generate area code (200-999, avoiding special codes)
	areaCode := fmt.Sprintf("%d%s", 2+randomInt(8), generateDigits(2))

	// Use 555 exchange - reserved for fictional use
	exchange := "555"

	// Generate subscriber number (0100-0199 range is specifically fictional)
	subscriber := fmt.Sprintf("01%02d", randomInt(100))

	digits := areaCode + exchange + subscriber
	return formatPhone(digits, format)
}

// UKPhoneGenerator generates UK phone numbers.
type UKPhoneGenerator struct {
	BaseGenerator
}

// NewUKPhoneGenerator creates a new UK phone generator.
func NewUKPhoneGenerator() *UKPhoneGenerator {
	return &UKPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "UK_PHONE"},
	}
}

// ukFictionalPrefix represents an Ofcom-reserved fictional phone prefix.
type ukFictionalPrefix struct {
	areaCode string // Area code without leading 0
	exchange string // Exchange/local prefix
	isMobile bool   // Whether this is a mobile number
}

// Ofcom-reserved ranges for dramatic use (TV, radio, etc.)
var ukFictionalPrefixes = []ukFictionalPrefix{
	{"20", "7946 0", false}, // London
	{"117", "496 0", false}, // Bristol
	{"131", "496 0", false}, // Edinburgh
	{"161", "496 0", false}, // Manchester
	{"7700", "900", true},   // Mobile
}

// Generate produces a UK phone number using Ofcom-reserved fictional ranges.
func (g *UKPhoneGenerator) Generate(input string) string {
	// Detect if input has +44 prefix
	hasCountryCode := strings.Contains(input, "+44")

	// Detect if input looks like a mobile (starts with 07)
	isMobile := strings.Contains(input, "07") || strings.Contains(input, "+447")

	// Select appropriate fictional prefix
	var prefix ukFictionalPrefix
	if isMobile {
		prefix = ukFictionalPrefixes[4] // Mobile prefix
	} else {
		// Pick a random landline prefix
		prefix = ukFictionalPrefixes[randomInt(4)]
	}

	// Generate subscriber number (3 digits for the 0xxx part)
	subscriber := fmt.Sprintf("%03d", randomInt(1000))

	if hasCountryCode {
		return fmt.Sprintf("+44 %s %s%s", prefix.areaCode, prefix.exchange, subscriber)
	}
	return fmt.Sprintf("0%s %s%s", prefix.areaCode, prefix.exchange, subscriber)
}

// InternationalPhoneGenerator generates international phone numbers.
type InternationalPhoneGenerator struct {
	BaseGenerator
}

// NewInternationalPhoneGenerator creates a new international phone generator.
func NewInternationalPhoneGenerator() *InternationalPhoneGenerator {
	return &InternationalPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "INTERNATIONAL_PHONE"},
	}
}

// Generate produces an international phone number with country code.
func (g *InternationalPhoneGenerator) Generate(input string) string {
	// Generate country code (1-3 digits)
	countryCode := fmt.Sprintf("%d", 1+randomInt(99))

	// Generate area code
	areaCode := generateDigits(3)

	// Generate local number
	localNumber := generateDigits(7)

	return fmt.Sprintf("+%s %s %s",
		countryCode, areaCode, localNumber)
}

// WorldwidePhoneGenerator generates phone numbers in various formats.
type WorldwidePhoneGenerator struct {
	BaseGenerator
}

// NewWorldwidePhoneGenerator creates a new worldwide phone generator.
func NewWorldwidePhoneGenerator() *WorldwidePhoneGenerator {
	return &WorldwidePhoneGenerator{
		BaseGenerator: BaseGenerator{name: "WORLDWIDE_PHONE"},
	}
}

// Generate produces a phone number matching the input length.
func (g *WorldwidePhoneGenerator) Generate(input string) string {
	// Count digits in input
	digitCount := 0
	for _, c := range input {
		if c >= '0' && c <= '9' {
			digitCount++
		}
	}

	if digitCount < 7 {
		digitCount = 10 // Default to 10 digits
	}

	// Generate matching number of digits
	return generateDigits(digitCount)
}
