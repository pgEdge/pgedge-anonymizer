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

	"github.com/pgedge/pgedge-anonymizer/internal/generator/data/countries"
)

// CountryFirstNameGenerator generates first names for a specific country.
type CountryFirstNameGenerator struct {
	BaseGenerator
	names []string
}

// NewCountryFirstNameGenerator creates a new country-specific first name generator.
func NewCountryFirstNameGenerator(country string, data *countries.CountryData) *CountryFirstNameGenerator {
	return &CountryFirstNameGenerator{
		BaseGenerator: BaseGenerator{name: country + "_FIRST_NAME"},
		names:         data.FirstNames,
	}
}

// Generate produces a first name for the country.
func (g *CountryFirstNameGenerator) Generate(input string) string {
	name := randomString(g.names)

	// Preserve case
	if strings.ToUpper(input) == input && len(input) > 1 {
		return strings.ToUpper(name)
	}
	if strings.ToLower(input) == input && len(input) > 1 {
		return strings.ToLower(name)
	}
	return name
}

// CountryLastNameGenerator generates last names for a specific country.
type CountryLastNameGenerator struct {
	BaseGenerator
	names []string
}

// NewCountryLastNameGenerator creates a new country-specific last name generator.
func NewCountryLastNameGenerator(country string, data *countries.CountryData) *CountryLastNameGenerator {
	return &CountryLastNameGenerator{
		BaseGenerator: BaseGenerator{name: country + "_LAST_NAME"},
		names:         data.LastNames,
	}
}

// Generate produces a last name for the country.
func (g *CountryLastNameGenerator) Generate(input string) string {
	name := randomString(g.names)

	// Preserve case
	if strings.ToUpper(input) == input && len(input) > 1 {
		return strings.ToUpper(name)
	}
	if strings.ToLower(input) == input && len(input) > 1 {
		return strings.ToLower(name)
	}
	return name
}

// CountryFullNameGenerator generates full names for a specific country.
type CountryFullNameGenerator struct {
	BaseGenerator
	firstNames []string
	lastNames  []string
}

// NewCountryFullNameGenerator creates a new country-specific full name generator.
func NewCountryFullNameGenerator(country string, data *countries.CountryData) *CountryFullNameGenerator {
	return &CountryFullNameGenerator{
		BaseGenerator: BaseGenerator{name: country + "_NAME"},
		firstNames:    data.FirstNames,
		lastNames:     data.LastNames,
	}
}

// Generate produces a full name for the country.
func (g *CountryFullNameGenerator) Generate(input string) string {
	firstName := randomString(g.firstNames)
	lastName := randomString(g.lastNames)

	// Detect comma-separated format (Last, First)
	if strings.Contains(input, ",") {
		result := lastName + ", " + firstName
		if strings.ToUpper(input) == input {
			return strings.ToUpper(result)
		}
		if strings.ToLower(input) == input {
			return strings.ToLower(result)
		}
		return result
	}

	result := firstName + " " + lastName
	if strings.ToUpper(input) == input {
		return strings.ToUpper(result)
	}
	if strings.ToLower(input) == input {
		return strings.ToLower(result)
	}
	return result
}

// CountryCityGenerator generates city names for a specific country.
type CountryCityGenerator struct {
	BaseGenerator
	cities []string
}

// NewCountryCityGenerator creates a new country-specific city generator.
func NewCountryCityGenerator(country string, data *countries.CountryData) *CountryCityGenerator {
	return &CountryCityGenerator{
		BaseGenerator: BaseGenerator{name: country + "_CITY"},
		cities:        data.Cities,
	}
}

// Generate produces a city name for the country.
func (g *CountryCityGenerator) Generate(input string) string {
	city := randomString(g.cities)

	// Preserve case
	if strings.ToUpper(input) == input && len(input) > 1 {
		return strings.ToUpper(city)
	}
	if strings.ToLower(input) == input && len(input) > 1 {
		return strings.ToLower(city)
	}
	return city
}

// WorldwideFirstNameGenerator generates first names from any country.
type WorldwideFirstNameGenerator struct {
	BaseGenerator
	allNames []string
}

// NewWorldwideFirstNameGenerator creates a worldwide first name generator.
func NewWorldwideFirstNameGenerator(data *countries.CountryDataSet) *WorldwideFirstNameGenerator {
	return &WorldwideFirstNameGenerator{
		BaseGenerator: BaseGenerator{name: "WORLDWIDE_FIRST_NAME"},
		allNames:      data.AllFirstNames(),
	}
}

// Generate produces a first name from any country.
func (g *WorldwideFirstNameGenerator) Generate(input string) string {
	name := randomString(g.allNames)

	if strings.ToUpper(input) == input && len(input) > 1 {
		return strings.ToUpper(name)
	}
	if strings.ToLower(input) == input && len(input) > 1 {
		return strings.ToLower(name)
	}
	return name
}

// WorldwideLastNameGenerator generates last names from any country.
type WorldwideLastNameGenerator struct {
	BaseGenerator
	allNames []string
}

// NewWorldwideLastNameGenerator creates a worldwide last name generator.
func NewWorldwideLastNameGenerator(data *countries.CountryDataSet) *WorldwideLastNameGenerator {
	return &WorldwideLastNameGenerator{
		BaseGenerator: BaseGenerator{name: "WORLDWIDE_LAST_NAME"},
		allNames:      data.AllLastNames(),
	}
}

// Generate produces a last name from any country.
func (g *WorldwideLastNameGenerator) Generate(input string) string {
	name := randomString(g.allNames)

	if strings.ToUpper(input) == input && len(input) > 1 {
		return strings.ToUpper(name)
	}
	if strings.ToLower(input) == input && len(input) > 1 {
		return strings.ToLower(name)
	}
	return name
}

// WorldwideNameGenerator generates full names from any country.
type WorldwideNameGenerator struct {
	BaseGenerator
	allFirstNames []string
	allLastNames  []string
}

// NewWorldwideNameGenerator creates a worldwide full name generator.
func NewWorldwideNameGenerator(data *countries.CountryDataSet) *WorldwideNameGenerator {
	return &WorldwideNameGenerator{
		BaseGenerator: BaseGenerator{name: "WORLDWIDE_NAME"},
		allFirstNames: data.AllFirstNames(),
		allLastNames:  data.AllLastNames(),
	}
}

// Generate produces a full name from any country.
func (g *WorldwideNameGenerator) Generate(input string) string {
	firstName := randomString(g.allFirstNames)
	lastName := randomString(g.allLastNames)

	if strings.Contains(input, ",") {
		result := lastName + ", " + firstName
		if strings.ToUpper(input) == input {
			return strings.ToUpper(result)
		}
		if strings.ToLower(input) == input {
			return strings.ToLower(result)
		}
		return result
	}

	result := firstName + " " + lastName
	if strings.ToUpper(input) == input {
		return strings.ToUpper(result)
	}
	if strings.ToLower(input) == input {
		return strings.ToLower(result)
	}
	return result
}

// WorldwideCityGenerator generates city names from any country.
type WorldwideCityGenerator struct {
	BaseGenerator
	allCities []string
}

// NewWorldwideCityGenerator creates a worldwide city generator.
func NewWorldwideCityGenerator(data *countries.CountryDataSet) *WorldwideCityGenerator {
	return &WorldwideCityGenerator{
		BaseGenerator: BaseGenerator{name: "WORLDWIDE_CITY"},
		allCities:     data.AllCities(),
	}
}

// Generate produces a city name from any country.
func (g *WorldwideCityGenerator) Generate(input string) string {
	city := randomString(g.allCities)

	if strings.ToUpper(input) == input && len(input) > 1 {
		return strings.ToUpper(city)
	}
	if strings.ToLower(input) == input && len(input) > 1 {
		return strings.ToLower(city)
	}
	return city
}
