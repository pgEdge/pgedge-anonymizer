/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package countries provides country-specific data for anonymization.
package countries

import (
	_ "embed"
	"strings"
)

// Country code constants
const (
	AU = "AU" // Australia
	CA = "CA" // Canada
	DE = "DE" // Germany
	ES = "ES" // Spain
	FI = "FI" // Finland
	FR = "FR" // France
	IE = "IE" // Ireland
	IN = "IN" // India
	IT = "IT" // Italy
	JP = "JP" // Japan
	KR = "KR" // South Korea
	MX = "MX" // Mexico
	NO = "NO" // Norway
	NZ = "NZ" // New Zealand
	PK = "PK" // Pakistan
	SE = "SE" // Sweden
	SG = "SG" // Singapore
	UK = "UK" // United Kingdom
	US = "US" // United States
)

// AllCountries lists all supported country codes
var AllCountries = []string{AU, CA, DE, ES, FI, FR, IE, IN, IT, JP, KR, MX, NO, NZ, PK, SE, SG, UK, US}

// Embedded data files for each country
//
//go:embed au_first_names.txt
var auFirstNamesRaw string

//go:embed au_last_names.txt
var auLastNamesRaw string

//go:embed au_cities.txt
var auCitiesRaw string

//go:embed ca_first_names.txt
var caFirstNamesRaw string

//go:embed ca_last_names.txt
var caLastNamesRaw string

//go:embed ca_cities.txt
var caCitiesRaw string

//go:embed de_first_names.txt
var deFirstNamesRaw string

//go:embed de_last_names.txt
var deLastNamesRaw string

//go:embed de_cities.txt
var deCitiesRaw string

//go:embed es_first_names.txt
var esFirstNamesRaw string

//go:embed es_last_names.txt
var esLastNamesRaw string

//go:embed es_cities.txt
var esCitiesRaw string

//go:embed fi_first_names.txt
var fiFirstNamesRaw string

//go:embed fi_last_names.txt
var fiLastNamesRaw string

//go:embed fi_cities.txt
var fiCitiesRaw string

//go:embed fr_first_names.txt
var frFirstNamesRaw string

//go:embed fr_last_names.txt
var frLastNamesRaw string

//go:embed fr_cities.txt
var frCitiesRaw string

//go:embed ie_first_names.txt
var ieFirstNamesRaw string

//go:embed ie_last_names.txt
var ieLastNamesRaw string

//go:embed ie_cities.txt
var ieCitiesRaw string

//go:embed in_first_names.txt
var inFirstNamesRaw string

//go:embed in_last_names.txt
var inLastNamesRaw string

//go:embed in_cities.txt
var inCitiesRaw string

//go:embed it_first_names.txt
var itFirstNamesRaw string

//go:embed it_last_names.txt
var itLastNamesRaw string

//go:embed it_cities.txt
var itCitiesRaw string

//go:embed jp_first_names.txt
var jpFirstNamesRaw string

//go:embed jp_last_names.txt
var jpLastNamesRaw string

//go:embed jp_cities.txt
var jpCitiesRaw string

//go:embed kr_first_names.txt
var krFirstNamesRaw string

//go:embed kr_last_names.txt
var krLastNamesRaw string

//go:embed kr_cities.txt
var krCitiesRaw string

//go:embed mx_first_names.txt
var mxFirstNamesRaw string

//go:embed mx_last_names.txt
var mxLastNamesRaw string

//go:embed mx_cities.txt
var mxCitiesRaw string

//go:embed no_first_names.txt
var noFirstNamesRaw string

//go:embed no_last_names.txt
var noLastNamesRaw string

//go:embed no_cities.txt
var noCitiesRaw string

//go:embed nz_first_names.txt
var nzFirstNamesRaw string

//go:embed nz_last_names.txt
var nzLastNamesRaw string

//go:embed nz_cities.txt
var nzCitiesRaw string

//go:embed pk_first_names.txt
var pkFirstNamesRaw string

//go:embed pk_last_names.txt
var pkLastNamesRaw string

//go:embed pk_cities.txt
var pkCitiesRaw string

//go:embed se_first_names.txt
var seFirstNamesRaw string

//go:embed se_last_names.txt
var seLastNamesRaw string

//go:embed se_cities.txt
var seCitiesRaw string

//go:embed sg_first_names.txt
var sgFirstNamesRaw string

//go:embed sg_last_names.txt
var sgLastNamesRaw string

//go:embed sg_cities.txt
var sgCitiesRaw string

//go:embed uk_first_names.txt
var ukFirstNamesRaw string

//go:embed uk_last_names.txt
var ukLastNamesRaw string

//go:embed uk_cities.txt
var ukCitiesRaw string

//go:embed us_first_names.txt
var usFirstNamesRaw string

//go:embed us_last_names.txt
var usLastNamesRaw string

//go:embed us_cities.txt
var usCitiesRaw string

// CountryData holds country-specific name and location data
type CountryData struct {
	FirstNames []string
	LastNames  []string
	Cities     []string
}

// CountryDataSet holds data for all countries
type CountryDataSet struct {
	Countries map[string]*CountryData
}

// parseLines splits raw text into lines, filtering empty lines and comments
func parseLines(raw string) []string {
	lines := strings.Split(raw, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			result = append(result, line)
		}
	}
	return result
}

// Load parses all embedded country data files and returns a CountryDataSet
func Load() *CountryDataSet {
	ds := &CountryDataSet{
		Countries: make(map[string]*CountryData),
	}

	ds.Countries[AU] = &CountryData{
		FirstNames: parseLines(auFirstNamesRaw),
		LastNames:  parseLines(auLastNamesRaw),
		Cities:     parseLines(auCitiesRaw),
	}

	ds.Countries[CA] = &CountryData{
		FirstNames: parseLines(caFirstNamesRaw),
		LastNames:  parseLines(caLastNamesRaw),
		Cities:     parseLines(caCitiesRaw),
	}

	ds.Countries[DE] = &CountryData{
		FirstNames: parseLines(deFirstNamesRaw),
		LastNames:  parseLines(deLastNamesRaw),
		Cities:     parseLines(deCitiesRaw),
	}

	ds.Countries[ES] = &CountryData{
		FirstNames: parseLines(esFirstNamesRaw),
		LastNames:  parseLines(esLastNamesRaw),
		Cities:     parseLines(esCitiesRaw),
	}

	ds.Countries[FI] = &CountryData{
		FirstNames: parseLines(fiFirstNamesRaw),
		LastNames:  parseLines(fiLastNamesRaw),
		Cities:     parseLines(fiCitiesRaw),
	}

	ds.Countries[FR] = &CountryData{
		FirstNames: parseLines(frFirstNamesRaw),
		LastNames:  parseLines(frLastNamesRaw),
		Cities:     parseLines(frCitiesRaw),
	}

	ds.Countries[IE] = &CountryData{
		FirstNames: parseLines(ieFirstNamesRaw),
		LastNames:  parseLines(ieLastNamesRaw),
		Cities:     parseLines(ieCitiesRaw),
	}

	ds.Countries[IN] = &CountryData{
		FirstNames: parseLines(inFirstNamesRaw),
		LastNames:  parseLines(inLastNamesRaw),
		Cities:     parseLines(inCitiesRaw),
	}

	ds.Countries[IT] = &CountryData{
		FirstNames: parseLines(itFirstNamesRaw),
		LastNames:  parseLines(itLastNamesRaw),
		Cities:     parseLines(itCitiesRaw),
	}

	ds.Countries[JP] = &CountryData{
		FirstNames: parseLines(jpFirstNamesRaw),
		LastNames:  parseLines(jpLastNamesRaw),
		Cities:     parseLines(jpCitiesRaw),
	}

	ds.Countries[KR] = &CountryData{
		FirstNames: parseLines(krFirstNamesRaw),
		LastNames:  parseLines(krLastNamesRaw),
		Cities:     parseLines(krCitiesRaw),
	}

	ds.Countries[MX] = &CountryData{
		FirstNames: parseLines(mxFirstNamesRaw),
		LastNames:  parseLines(mxLastNamesRaw),
		Cities:     parseLines(mxCitiesRaw),
	}

	ds.Countries[NO] = &CountryData{
		FirstNames: parseLines(noFirstNamesRaw),
		LastNames:  parseLines(noLastNamesRaw),
		Cities:     parseLines(noCitiesRaw),
	}

	ds.Countries[NZ] = &CountryData{
		FirstNames: parseLines(nzFirstNamesRaw),
		LastNames:  parseLines(nzLastNamesRaw),
		Cities:     parseLines(nzCitiesRaw),
	}

	ds.Countries[PK] = &CountryData{
		FirstNames: parseLines(pkFirstNamesRaw),
		LastNames:  parseLines(pkLastNamesRaw),
		Cities:     parseLines(pkCitiesRaw),
	}

	ds.Countries[SE] = &CountryData{
		FirstNames: parseLines(seFirstNamesRaw),
		LastNames:  parseLines(seLastNamesRaw),
		Cities:     parseLines(seCitiesRaw),
	}

	ds.Countries[SG] = &CountryData{
		FirstNames: parseLines(sgFirstNamesRaw),
		LastNames:  parseLines(sgLastNamesRaw),
		Cities:     parseLines(sgCitiesRaw),
	}

	ds.Countries[UK] = &CountryData{
		FirstNames: parseLines(ukFirstNamesRaw),
		LastNames:  parseLines(ukLastNamesRaw),
		Cities:     parseLines(ukCitiesRaw),
	}

	ds.Countries[US] = &CountryData{
		FirstNames: parseLines(usFirstNamesRaw),
		LastNames:  parseLines(usLastNamesRaw),
		Cities:     parseLines(usCitiesRaw),
	}

	return ds
}

// Get returns the data for a specific country
func (ds *CountryDataSet) Get(country string) *CountryData {
	return ds.Countries[country]
}

// AllFirstNames returns all first names from all countries combined
func (ds *CountryDataSet) AllFirstNames() []string {
	var all []string
	for _, cd := range ds.Countries {
		all = append(all, cd.FirstNames...)
	}
	return all
}

// AllLastNames returns all last names from all countries combined
func (ds *CountryDataSet) AllLastNames() []string {
	var all []string
	for _, cd := range ds.Countries {
		all = append(all, cd.LastNames...)
	}
	return all
}

// AllCities returns all cities from all countries combined
func (ds *CountryDataSet) AllCities() []string {
	var all []string
	for _, cd := range ds.Countries {
		all = append(all, cd.Cities...)
	}
	return all
}
