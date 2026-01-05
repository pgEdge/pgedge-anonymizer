/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package data provides embedded data files for realistic anonymization.
package data

import (
	_ "embed"
	"strings"
)

//go:embed first_names.txt
var firstNamesRaw string

//go:embed last_names.txt
var lastNamesRaw string

//go:embed street_names.txt
var streetNamesRaw string

//go:embed cities.txt
var citiesRaw string

//go:embed domains.txt
var domainsRaw string

//go:embed lorem_words.txt
var loremWordsRaw string

// DataSet provides access to parsed data lists.
type DataSet struct {
	FirstNames  []string
	LastNames   []string
	StreetNames []string
	Cities      []string
	Domains     []string
	LoremWords  []string
}

// parseLines splits raw text into lines, filtering empty lines.
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

// Load parses all embedded data files and returns a DataSet.
func Load() *DataSet {
	return &DataSet{
		FirstNames:  parseLines(firstNamesRaw),
		LastNames:   parseLines(lastNamesRaw),
		StreetNames: parseLines(streetNamesRaw),
		Cities:      parseLines(citiesRaw),
		Domains:     parseLines(domainsRaw),
		LoremWords:  parseLines(loremWordsRaw),
	}
}
