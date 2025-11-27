/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package generator

import (
	"fmt"
	"strings"
)

// SSNGenerator generates US Social Security Numbers.
type SSNGenerator struct {
	BaseGenerator
}

// NewSSNGenerator creates a new SSN generator.
func NewSSNGenerator() *SSNGenerator {
	return &SSNGenerator{
		BaseGenerator: BaseGenerator{name: "US_SSN"},
	}
}

// Generate produces a US Social Security Number.
func (g *SSNGenerator) Generate(input string) string {
	// Generate area number (001-665, 667-899)
	// Avoid 000, 666, and 900-999
	area := g.generateValidArea()

	// Generate group number (01-99)
	group := 1 + randomInt(99)

	// Generate serial number (0001-9999)
	serial := 1 + randomInt(9999)

	// Detect format (with or without dashes)
	if strings.Contains(input, "-") {
		return fmt.Sprintf("%03d-%02d-%04d", area, group, serial)
	}

	if strings.Contains(input, " ") {
		return fmt.Sprintf("%03d %02d %04d", area, group, serial)
	}

	// No separator
	return fmt.Sprintf("%03d%02d%04d", area, group, serial)
}

// generateValidArea generates a valid SSN area number.
func (g *SSNGenerator) generateValidArea() int {
	for {
		area := 1 + randomInt(899)
		// Avoid 000, 666, and 900-999
		if area != 666 && area < 900 {
			return area
		}
	}
}
