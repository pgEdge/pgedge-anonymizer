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
	"time"
)

// DOBGenerator generates random dates of birth.
type DOBGenerator struct {
	BaseGenerator
	minAge int
	maxAge int
}

// NewDOBGenerator creates a generator for any age date of birth.
func NewDOBGenerator() *DOBGenerator {
	return &DOBGenerator{
		BaseGenerator: BaseGenerator{name: "DOB"},
		minAge:        0,
		maxAge:        100,
	}
}

// NewDOBOver13Generator creates a generator for dates of birth over 13.
func NewDOBOver13Generator() *DOBGenerator {
	return &DOBGenerator{
		BaseGenerator: BaseGenerator{name: "DOB_OVER_13"},
		minAge:        13,
		maxAge:        100,
	}
}

// NewDOBOver16Generator creates a generator for dates of birth over 16.
func NewDOBOver16Generator() *DOBGenerator {
	return &DOBGenerator{
		BaseGenerator: BaseGenerator{name: "DOB_OVER_16"},
		minAge:        16,
		maxAge:        100,
	}
}

// NewDOBOver18Generator creates a generator for dates of birth over 18.
func NewDOBOver18Generator() *DOBGenerator {
	return &DOBGenerator{
		BaseGenerator: BaseGenerator{name: "DOB_OVER_18"},
		minAge:        18,
		maxAge:        100,
	}
}

// NewDOBOver21Generator creates a generator for dates of birth over 21.
func NewDOBOver21Generator() *DOBGenerator {
	return &DOBGenerator{
		BaseGenerator: BaseGenerator{name: "DOB_OVER_21"},
		minAge:        21,
		maxAge:        100,
	}
}

// Generate produces a date of birth within the configured age range.
func (g *DOBGenerator) Generate(input string) string {
	now := time.Now()

	// Calculate date range
	maxDate := now.AddDate(-g.minAge, 0, 0) // Youngest possible
	minDate := now.AddDate(-g.maxAge, 0, 0) // Oldest possible

	// Random date within range
	dayRange := int(maxDate.Sub(minDate).Hours() / 24)
	if dayRange <= 0 {
		dayRange = 365
	}
	randomDays := randomInt(dayRange)
	dob := minDate.AddDate(0, 0, randomDays)

	// Detect format from input
	format := detectDateFormat(input)
	return formatDate(dob, format)
}

// dateFormat represents detected date format.
type dateFormat int

const (
	formatISO     dateFormat = iota // YYYY-MM-DD
	formatUS                        // MM/DD/YYYY
	formatEU                        // DD/MM/YYYY
	formatUSShort                   // MM/DD/YY
	formatEUShort                   // DD/MM/YY
	formatLong                      // Month DD, YYYY
)

// detectDateFormat attempts to detect the date format from input.
func detectDateFormat(input string) dateFormat {
	// Check for ISO format (YYYY-MM-DD)
	if len(input) >= 10 && input[4] == '-' && input[7] == '-' {
		return formatISO
	}

	// Check for slash-separated formats
	if strings.Contains(input, "/") {
		parts := strings.Split(input, "/")
		if len(parts) == 3 {
			// If first part > 12, likely DD/MM
			if len(parts[0]) == 2 {
				// Could be DD/MM or MM/DD
				// Default to US format for ambiguous cases
				if len(parts[2]) == 4 {
					return formatUS
				}
				return formatUSShort
			}
		}
	}

	// Check for written format (contains month name)
	months := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
		"Jan", "Feb", "Mar", "Apr", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	for _, m := range months {
		if strings.Contains(input, m) {
			return formatLong
		}
	}

	// Default to ISO
	return formatISO
}

// formatDate formats a date according to the detected format.
func formatDate(t time.Time, format dateFormat) string {
	switch format {
	case formatISO:
		return t.Format("2006-01-02")
	case formatUS:
		return t.Format("01/02/2006")
	case formatEU:
		return t.Format("02/01/2006")
	case formatUSShort:
		return t.Format("01/02/06")
	case formatEUShort:
		return t.Format("02/01/06")
	case formatLong:
		return t.Format("January 2, 2006")
	default:
		return t.Format("2006-01-02")
	}
}

// WithAgeRange creates a DOB generator with custom age range.
func WithAgeRange(name string, minAge, maxAge int) *DOBGenerator {
	return &DOBGenerator{
		BaseGenerator: BaseGenerator{name: name},
		minAge:        minAge,
		maxAge:        maxAge,
	}
}

// Unused but kept for completeness
var _ = fmt.Sprintf
