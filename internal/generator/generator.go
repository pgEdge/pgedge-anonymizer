/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package generator provides data generators for anonymization patterns.
package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
)

// Generator defines the interface for all anonymization generators.
type Generator interface {
	// Name returns the pattern name (e.g., "US_PHONE", "EMAIL")
	Name() string

	// Generate produces an anonymized value for the given input.
	// The implementation should preserve format characteristics where
	// appropriate (e.g., phone number format, text length).
	Generate(input string) string
}

// Registry holds all registered generators indexed by name.
type Registry struct {
	generators map[string]Generator
	mu         sync.RWMutex
}

// NewRegistry creates an empty generator registry.
func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]Generator),
	}
}

// Register adds a generator to the registry.
func (r *Registry) Register(g Generator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.generators[g.Name()] = g
}

// Get retrieves a generator by name.
func (r *Registry) Get(name string) (Generator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.generators[name]
	return g, ok
}

// List returns all registered generator names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.generators))
	for name := range r.generators {
		names = append(names, name)
	}
	return names
}

// randomInt returns a cryptographically secure random integer in [0, max).
func randomInt(max int) int {
	if max <= 0 {
		return 0
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		// Fall back to a simple value on error (should never happen)
		return 0
	}
	return int(n.Int64())
}

// randomDigit returns a random digit '0'-'9'.
func randomDigit() byte {
	return byte('0' + randomInt(10))
}

// randomDigitNonZero returns a random digit '1'-'9'.
func randomDigitNonZero() byte {
	return byte('1' + randomInt(9))
}

// randomString selects a random string from a slice.
func randomString(choices []string) string {
	if len(choices) == 0 {
		return ""
	}
	return choices[randomInt(len(choices))]
}

// generateDigits generates a string of n random digits.
func generateDigits(n int) string {
	result := make([]byte, n)
	for i := range result {
		result[i] = randomDigit()
	}
	return string(result)
}

// luhnCheckDigit calculates the Luhn check digit for a sequence of digits.
func luhnCheckDigit(digits string) byte {
	sum := 0
	// Process from right to left, doubling every second digit
	for i := len(digits) - 1; i >= 0; i-- {
		d := int(digits[i] - '0')
		if (len(digits)-1-i)%2 == 0 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	checkDigit := (10 - (sum % 10)) % 10
	return byte('0' + checkDigit)
}

// BaseGenerator provides common functionality for generators.
type BaseGenerator struct {
	name string
}

// Name returns the generator name.
func (b *BaseGenerator) Name() string {
	return b.name
}

// detectPhoneFormat detects the format of a phone number string.
// Returns format indicators for formatting output.
type phoneFormat struct {
	hasParens  bool
	separator  byte // '-', '.', ' ', or 0 for none
	hasCountry bool
}

func detectPhoneFormat(input string) phoneFormat {
	var pf phoneFormat

	for _, c := range input {
		switch c {
		case '(':
			pf.hasParens = true
		case '-':
			if pf.separator == 0 {
				pf.separator = '-'
			}
		case '.':
			if pf.separator == 0 {
				pf.separator = '.'
			}
		case ' ':
			if pf.separator == 0 {
				pf.separator = ' '
			}
		case '+':
			pf.hasCountry = true
		}
	}

	return pf
}

// formatPhone formats digits according to the detected format.
func formatPhone(digits string, format phoneFormat) string {
	if len(digits) < 10 {
		return digits
	}

	sep := string(format.separator)
	if format.separator == 0 {
		sep = ""
	}

	if format.hasParens {
		return fmt.Sprintf("(%s) %s%s%s",
			digits[0:3], digits[3:6], sep, digits[6:10])
	}

	if format.separator != 0 {
		return fmt.Sprintf("%s%s%s%s%s",
			digits[0:3], sep, digits[3:6], sep, digits[6:10])
	}

	return digits
}
