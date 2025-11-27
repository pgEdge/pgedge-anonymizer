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

// FormatType indicates the type of format pattern.
type FormatType string

const (
	FormatTypeDate   FormatType = "date"
	FormatTypeMask   FormatType = "mask"
	FormatTypeNumber FormatType = "number"
)

// FormatConfig holds configuration for a format-based generator.
type FormatConfig struct {
	Format  string     // The format string
	Type    FormatType // Type of format (date, mask, number)
	Min     int64      // Minimum value for number type
	Max     int64      // Maximum value for number type
	MinYear int        // Minimum year for date type
	MaxYear int        // Maximum year for date type
}

// FormatGenerator generates values based on format strings.
type FormatGenerator struct {
	BaseGenerator
	config FormatConfig
}

// NewFormatGenerator creates a new format-based generator.
func NewFormatGenerator(name string, config FormatConfig) *FormatGenerator {
	// Set defaults
	if config.MinYear == 0 {
		config.MinYear = 1950
	}
	if config.MaxYear == 0 {
		config.MaxYear = time.Now().Year()
	}
	if config.Max == 0 && config.Type == FormatTypeNumber {
		config.Max = 999999999
	}

	return &FormatGenerator{
		BaseGenerator: BaseGenerator{name: name},
		config:        config,
	}
}

// Generate produces a value matching the format.
func (g *FormatGenerator) Generate(input string) string {
	switch g.config.Type {
	case FormatTypeDate:
		return g.generateDate()
	case FormatTypeMask:
		return g.generateMask()
	case FormatTypeNumber:
		return g.generateNumber()
	default:
		// Auto-detect based on format string
		if containsDateCodes(g.config.Format) {
			return g.generateDate()
		}
		if containsNumberCodes(g.config.Format) {
			return g.generateNumber()
		}
		return g.generateMask()
	}
}

// generateDate generates a random date in the specified format.
// Supports strftime-like format codes.
func (g *FormatGenerator) generateDate() string {
	// Generate random date components
	year := g.config.MinYear + randomInt(g.config.MaxYear-g.config.MinYear+1)
	month := 1 + randomInt(12)
	day := 1 + randomInt(28) // Safe for all months
	hour := randomInt(24)
	minute := randomInt(60)
	second := randomInt(60)

	// Replace format codes
	result := g.config.Format
	result = strings.ReplaceAll(result, "%Y", fmt.Sprintf("%04d", year))
	result = strings.ReplaceAll(result, "%y", fmt.Sprintf("%02d", year%100))
	result = strings.ReplaceAll(result, "%m", fmt.Sprintf("%02d", month))
	result = strings.ReplaceAll(result, "%d", fmt.Sprintf("%02d", day))
	result = strings.ReplaceAll(result, "%H", fmt.Sprintf("%02d", hour))
	result = strings.ReplaceAll(result, "%M", fmt.Sprintf("%02d", minute))
	result = strings.ReplaceAll(result, "%S", fmt.Sprintf("%02d", second))
	result = strings.ReplaceAll(result, "%I", fmt.Sprintf("%02d", (hour%12)+1))

	// Month and day names
	monthNames := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}
	monthAbbr := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday",
		"Thursday", "Friday", "Saturday"}
	dayAbbr := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	result = strings.ReplaceAll(result, "%B", monthNames[month-1])
	result = strings.ReplaceAll(result, "%b", monthAbbr[month-1])
	result = strings.ReplaceAll(result, "%A", dayNames[randomInt(7)])
	result = strings.ReplaceAll(result, "%a", dayAbbr[randomInt(7)])

	// AM/PM
	if hour < 12 {
		result = strings.ReplaceAll(result, "%p", "AM")
		result = strings.ReplaceAll(result, "%P", "am")
	} else {
		result = strings.ReplaceAll(result, "%p", "PM")
		result = strings.ReplaceAll(result, "%P", "pm")
	}

	return result
}

// generateMask generates a value matching a mask pattern.
// Placeholders:
//
//	# - random digit (0-9)
//	9 - random digit (0-9) - alternative
//	A - random uppercase letter
//	a - random lowercase letter
//	X - random uppercase alphanumeric
//	x - random lowercase alphanumeric
//	* - random character (letter or digit)
//	\ - escape next character (use literal)
//
// All other characters are literals.
func (g *FormatGenerator) generateMask() string {
	var result strings.Builder
	format := g.config.Format
	escaped := false

	for i := 0; i < len(format); i++ {
		c := format[i]

		if escaped {
			result.WriteByte(c)
			escaped = false
			continue
		}

		switch c {
		case '\\':
			escaped = true
		case '#', '9':
			result.WriteByte(randomDigit())
		case 'A':
			result.WriteByte(randomUpperLetter())
		case 'a':
			result.WriteByte(randomLowerLetter())
		case 'X':
			if randomInt(2) == 0 {
				result.WriteByte(randomDigit())
			} else {
				result.WriteByte(randomUpperLetter())
			}
		case 'x':
			if randomInt(2) == 0 {
				result.WriteByte(randomDigit())
			} else {
				result.WriteByte(randomLowerLetter())
			}
		case '*':
			r := randomInt(3)
			if r == 0 {
				result.WriteByte(randomDigit())
			} else if r == 1 {
				result.WriteByte(randomUpperLetter())
			} else {
				result.WriteByte(randomLowerLetter())
			}
		default:
			result.WriteByte(c)
		}
	}

	return result.String()
}

// generateNumber generates a random number in the specified format.
// Supports printf-like format codes for integers.
func (g *FormatGenerator) generateNumber() string {
	min := g.config.Min
	max := g.config.Max
	if max <= min {
		max = min + 1000000
	}

	value := min + int64(randomInt(int(max-min+1)))
	return fmt.Sprintf(g.config.Format, value)
}

// randomUpperLetter returns a random uppercase letter A-Z.
func randomUpperLetter() byte {
	return byte('A' + randomInt(26))
}

// randomLowerLetter returns a random lowercase letter a-z.
func randomLowerLetter() byte {
	return byte('a' + randomInt(26))
}

// containsDateCodes checks if a format string contains date/time codes.
func containsDateCodes(format string) bool {
	dateCodes := []string{"%Y", "%y", "%m", "%d", "%H", "%M", "%S", "%I",
		"%B", "%b", "%A", "%a", "%p", "%P"}
	for _, code := range dateCodes {
		if strings.Contains(format, code) {
			return true
		}
	}
	return false
}

// containsNumberCodes checks if a format string contains printf number codes.
func containsNumberCodes(format string) bool {
	// Look for printf-style integer/float format codes
	// e.g., %d, %5d, %05d, %f, %.2f
	for i := 0; i < len(format)-1; i++ {
		if format[i] == '%' {
			// Skip escaped %%
			if format[i+1] == '%' {
				i++
				continue
			}
			// Look for d, f, or digits followed by d/f
			j := i + 1
			for j < len(format) && (format[j] >= '0' && format[j] <= '9' || format[j] == '.') {
				j++
			}
			if j < len(format) && (format[j] == 'd' || format[j] == 'f') {
				return true
			}
		}
	}
	return false
}

// DetectFormatType attempts to detect the format type from a format string.
func DetectFormatType(format string) FormatType {
	if containsDateCodes(format) {
		return FormatTypeDate
	}
	if containsNumberCodes(format) {
		return FormatTypeNumber
	}
	return FormatTypeMask
}
