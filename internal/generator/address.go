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

	"github.com/pgedge/pgedge-anonymizer/internal/generator/data/countries"
)

// AddressGenerator generates street addresses from worldwide data.
// This generator now uses diverse data from all supported countries.
type AddressGenerator struct {
	BaseGenerator
	worldwideGen *WorldwideAddressGenerator
}

// NewAddressGenerator creates a new address generator using worldwide data.
func NewAddressGenerator(countryData *countries.CountryDataSet) *AddressGenerator {
	return &AddressGenerator{
		BaseGenerator: BaseGenerator{name: "ADDRESS"},
		worldwideGen:  NewWorldwideAddressGenerator(countryData),
	}
}

// Generate produces a street address from a randomly selected country.
func (g *AddressGenerator) Generate(input string) string {
	return g.worldwideGen.Generate(input)
}

// USZipGenerator generates US ZIP codes.
type USZipGenerator struct {
	BaseGenerator
}

// NewUSZipGenerator creates a new US ZIP code generator.
func NewUSZipGenerator() *USZipGenerator {
	return &USZipGenerator{
		BaseGenerator: BaseGenerator{name: "US_ZIP"},
	}
}

// Generate produces a US ZIP code.
// It detects the format (5-digit or ZIP+4) and generates a matching format.
func (g *USZipGenerator) Generate(input string) string {
	// Check if input uses ZIP+4 format (12345-6789)
	if strings.Contains(input, "-") && len(input) >= 10 {
		return generateDigits(5) + "-" + generateDigits(4)
	}
	return generateDigits(5)
}

// CityGenerator generates city names from worldwide data.
// This generator now uses diverse data from all supported countries.
type CityGenerator struct {
	BaseGenerator
	allCities []string
}

// NewCityGenerator creates a new city name generator using worldwide data.
func NewCityGenerator(countryData *countries.CountryDataSet) *CityGenerator {
	return &CityGenerator{
		BaseGenerator: BaseGenerator{name: "CITY"},
		allCities:     countryData.AllCities(),
	}
}

// Generate produces a city name from any country.
// It preserves uppercase/lowercase formatting.
func (g *CityGenerator) Generate(input string) string {
	city := randomString(g.allCities)

	// Check for case preservation
	if strings.ToUpper(input) == input && len(input) > 1 {
		return strings.ToUpper(city)
	}
	if strings.ToLower(input) == input && len(input) > 1 {
		return strings.ToLower(city)
	}

	return city
}

// UKPostcodeGenerator generates UK postcodes.
type UKPostcodeGenerator struct {
	BaseGenerator
}

// NewUKPostcodeGenerator creates a new UK postcode generator.
func NewUKPostcodeGenerator() *UKPostcodeGenerator {
	return &UKPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "UK_POSTCODE"},
	}
}

// Generate produces a UK postcode.
// UK postcodes have formats like: SW1A 1AA, M1 1AE, B33 8TH, EC1A 1BB
func (g *UKPostcodeGenerator) Generate(input string) string {
	// UK postcode format: outward code + space + inward code
	// Outward: 2-4 chars (1-2 letters + 1-2 digits, optionally ending with letter)
	// Inward: 3 chars (digit + 2 letters)

	// Valid outward code letters (first position)
	firstLetters := "ABCDEFGHIJKLMNOPRSTUWYZ"
	// Valid letters for other positions (excluding CIKMOV)
	otherLetters := "ABDEFGHJLNPQRSTUWXYZ"

	// Generate outward code - use common formats
	var outward string
	format := randomInt(4)
	switch format {
	case 0: // A9 format (e.g., M1, B1)
		outward = string(firstLetters[randomInt(len(firstLetters))]) +
			string('1'+byte(randomInt(9)))
	case 1: // A99 format (e.g., M11, B33)
		outward = string(firstLetters[randomInt(len(firstLetters))]) +
			string('1'+byte(randomInt(9))) +
			string('0'+byte(randomInt(10)))
	case 2: // AA9 format (e.g., SW1, EC1)
		outward = string(firstLetters[randomInt(len(firstLetters))]) +
			string(otherLetters[randomInt(len(otherLetters))]) +
			string('1'+byte(randomInt(9)))
	default: // AA99 format (e.g., SW19, EC1A)
		outward = string(firstLetters[randomInt(len(firstLetters))]) +
			string(otherLetters[randomInt(len(otherLetters))]) +
			string('1'+byte(randomInt(9))) +
			string(otherLetters[randomInt(len(otherLetters))])
	}

	// Generate inward code (digit + 2 letters)
	inward := string('0'+byte(randomInt(10))) +
		string(otherLetters[randomInt(len(otherLetters))]) +
		string(otherLetters[randomInt(len(otherLetters))])

	// Check if input has space
	if strings.Contains(input, " ") {
		return outward + " " + inward
	}
	return outward + inward
}

// CAPostcodeGenerator generates Canadian postcodes.
type CAPostcodeGenerator struct {
	BaseGenerator
}

// NewCAPostcodeGenerator creates a new Canadian postcode generator.
func NewCAPostcodeGenerator() *CAPostcodeGenerator {
	return &CAPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "CA_POSTCODE"},
	}
}

// Generate produces a Canadian postcode.
// Canadian postcodes have the format: A9A 9A9 (letter-digit-letter space digit-letter-digit)
func (g *CAPostcodeGenerator) Generate(input string) string {
	// Valid letters for Canadian postcodes (excludes D, F, I, O, Q, U, W, Z in first position)
	// and D, F, I, O, Q, U in other positions
	firstLetters := "ABCEGHJKLMNPRSTVXY"
	otherLetters := "ABCEGHJKLMNPRSTVWXYZ"

	// Generate FSA (Forward Sortation Area) - first 3 characters
	fsa := string(firstLetters[randomInt(len(firstLetters))]) +
		string('0'+byte(randomInt(10))) +
		string(otherLetters[randomInt(len(otherLetters))])

	// Generate LDU (Local Delivery Unit) - last 3 characters
	ldu := string('0'+byte(randomInt(10))) +
		string(otherLetters[randomInt(len(otherLetters))]) +
		string('0'+byte(randomInt(10)))

	// Check if input has space
	if strings.Contains(input, " ") {
		return fsa + " " + ldu
	}
	return fsa + ldu
}

// WorldwidePostcodeGenerator generates postcodes in various international formats.
type WorldwidePostcodeGenerator struct {
	BaseGenerator
	usGen *USZipGenerator
	ukGen *UKPostcodeGenerator
	caGen *CAPostcodeGenerator
}

// NewWorldwidePostcodeGenerator creates a new worldwide postcode generator.
func NewWorldwidePostcodeGenerator() *WorldwidePostcodeGenerator {
	return &WorldwidePostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "WORLDWIDE_POSTCODE"},
		usGen:         NewUSZipGenerator(),
		ukGen:         NewUKPostcodeGenerator(),
		caGen:         NewCAPostcodeGenerator(),
	}
}

// Generate produces a postcode in a randomly selected international format.
func (g *WorldwidePostcodeGenerator) Generate(input string) string {
	// Try to detect the format from input
	inputLen := len(strings.ReplaceAll(input, " ", ""))

	// US ZIP: 5 or 9 digits
	if isAllDigits(input) {
		return g.usGen.Generate(input)
	}

	// Canadian postcode: exactly 6 alphanumeric (A9A9A9 or A9A 9A9)
	// Check this before UK as Canadian format is more specific
	if inputLen == 6 && hasAlternatingPattern(input) {
		return g.caGen.Generate(input)
	}

	// UK postcode: 5-8 alphanumeric, typically has letter at start
	if inputLen >= 5 && inputLen <= 8 && hasLetterAtStart(input) && hasDigitInMiddle(input) {
		return g.ukGen.Generate(input)
	}

	// Default: randomly select a format
	switch randomInt(3) {
	case 0:
		return g.usGen.Generate(input)
	case 1:
		return g.ukGen.Generate(input)
	default:
		return g.caGen.Generate(input)
	}
}

// Helper functions for format detection
func isAllDigits(s string) bool {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func hasLetterAtStart(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return false
	}
	c := s[0]
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func hasDigitInMiddle(s string) bool {
	s = strings.ReplaceAll(s, " ", "")
	if len(s) < 3 {
		return false
	}
	for i := 1; i < len(s)-1; i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return true
		}
	}
	return false
}

func hasAlternatingPattern(s string) bool {
	// Check for Canadian A9A9A9 pattern
	s = strings.ReplaceAll(s, " ", "")
	if len(s) != 6 {
		return false
	}
	for i, c := range s {
		isLetter := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		isDigit := c >= '0' && c <= '9'
		if i%2 == 0 && !isLetter {
			return false
		}
		if i%2 == 1 && !isDigit {
			return false
		}
	}
	return true
}
