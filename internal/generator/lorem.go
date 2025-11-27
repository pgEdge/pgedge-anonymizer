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
	"strings"
	"unicode"

	"github.com/pgedge/pgedge-anonymizer/internal/generator/data"
)

// LoremGenerator generates lorem ipsum text.
type LoremGenerator struct {
	BaseGenerator
	data *data.DataSet
}

// NewLoremGenerator creates a new lorem ipsum generator.
func NewLoremGenerator(d *data.DataSet) *LoremGenerator {
	return &LoremGenerator{
		BaseGenerator: BaseGenerator{name: "LOREMIPSUM"},
		data:          d,
	}
}

// capitalizeFirst capitalizes the first letter of a string.
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Generate produces lorem ipsum text approximately matching the input length.
func (g *LoremGenerator) Generate(input string) string {
	targetLen := len(input)
	if targetLen == 0 {
		targetLen = 50 // Default minimum
	}

	var result strings.Builder
	wordCount := 0

	for result.Len() < targetLen {
		word := randomString(g.data.LoremWords)

		if result.Len() > 0 {
			// Check if adding this word would exceed target
			if result.Len()+1+len(word) > targetLen+10 {
				break
			}
			result.WriteByte(' ')
		}

		// Capitalize first word
		if wordCount == 0 {
			word = capitalizeFirst(word)
		}

		result.WriteString(word)
		wordCount++

		// Add period occasionally for sentence breaks
		if wordCount > 0 && wordCount%8 == 0 && result.Len() < targetLen-10 {
			result.WriteByte('.')
			// Capitalize next word
			if result.Len() < targetLen-5 {
				result.WriteByte(' ')
				nextWord := randomString(g.data.LoremWords)
				result.WriteString(capitalizeFirst(nextWord))
				wordCount++
			}
		}
	}

	// Ensure we end with a period if we have content
	text := result.String()
	if len(text) > 0 && !strings.HasSuffix(text, ".") {
		text += "."
	}

	return text
}
