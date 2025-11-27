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

	"github.com/pgedge/pgedge-anonymizer/internal/generator/data"
)

// NameGenerator generates person names.
type NameGenerator struct {
	BaseGenerator
	data *data.DataSet
}

// NewNameGenerator creates a new name generator.
func NewNameGenerator(d *data.DataSet) *NameGenerator {
	return &NameGenerator{
		BaseGenerator: BaseGenerator{name: "PERSON_NAME"},
		data:          d,
	}
}

// Generate produces a person name, attempting to match the input format.
func (g *NameGenerator) Generate(input string) string {
	firstName := randomString(g.data.FirstNames)
	lastName := randomString(g.data.LastNames)

	// Detect format: "Last, First" vs "First Last"
	if strings.Contains(input, ",") {
		return lastName + ", " + firstName
	}

	// Check if input appears to be all caps
	if input == strings.ToUpper(input) && len(input) > 1 {
		return strings.ToUpper(firstName + " " + lastName)
	}

	// Check if input appears to be all lower
	if input == strings.ToLower(input) && len(input) > 1 {
		return strings.ToLower(firstName + " " + lastName)
	}

	return firstName + " " + lastName
}

// FirstNameGenerator generates first names only.
type FirstNameGenerator struct {
	BaseGenerator
	data *data.DataSet
}

// NewFirstNameGenerator creates a new first name generator.
func NewFirstNameGenerator(d *data.DataSet) *FirstNameGenerator {
	return &FirstNameGenerator{
		BaseGenerator: BaseGenerator{name: "PERSON_FIRST_NAME"},
		data:          d,
	}
}

// Generate produces a first name.
func (g *FirstNameGenerator) Generate(input string) string {
	firstName := randomString(g.data.FirstNames)

	// Match case of input
	if input == strings.ToUpper(input) && len(input) > 1 {
		return strings.ToUpper(firstName)
	}
	if input == strings.ToLower(input) && len(input) > 1 {
		return strings.ToLower(firstName)
	}

	return firstName
}

// LastNameGenerator generates last names only.
type LastNameGenerator struct {
	BaseGenerator
	data *data.DataSet
}

// NewLastNameGenerator creates a new last name generator.
func NewLastNameGenerator(d *data.DataSet) *LastNameGenerator {
	return &LastNameGenerator{
		BaseGenerator: BaseGenerator{name: "PERSON_LAST_NAME"},
		data:          d,
	}
}

// Generate produces a last name.
func (g *LastNameGenerator) Generate(input string) string {
	lastName := randomString(g.data.LastNames)

	// Match case of input
	if input == strings.ToUpper(input) && len(input) > 1 {
		return strings.ToUpper(lastName)
	}
	if input == strings.ToLower(input) && len(input) > 1 {
		return strings.ToLower(lastName)
	}

	return lastName
}
