/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package jsonpath provides JSON path extraction and replacement operations
// for anonymizing values within JSON/JSONB columns.
package jsonpath

import (
	"fmt"
	"log"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

// PathMatch represents a value found at a JSON path along with its location.
type PathMatch struct {
	Path  string // The concrete path to this value (e.g., "$.users[0].email")
	Value string // The extracted string value
}

// Processor handles JSON path operations for anonymization.
type Processor struct {
	quiet bool // suppress warnings
}

// NewProcessor creates a new JSON path processor.
func NewProcessor(quiet bool) *Processor {
	return &Processor{quiet: quiet}
}

// Extract finds all string values matching a JSON path expression.
// For paths with wildcards (e.g., $.users[*].email), returns all matches.
// Non-string values (objects, arrays, null) are skipped with a warning.
func (p *Processor) Extract(jsonData []byte, pathExpr string) ([]PathMatch, error) {
	// Parse the JSON
	data, err := oj.Parse(jsonData)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Parse the JSON path expression
	path, err := jp.ParseString(pathExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON path %q: %w", pathExpr, err)
	}

	// Get all matching values
	results := path.Get(data)
	if len(results) == 0 {
		return nil, nil // No matches, not an error
	}

	var matches []PathMatch
	for i, result := range results {
		// Only process string values
		switch v := result.(type) {
		case string:
			// Build the concrete path for this match
			concretePath := buildConcretePath(pathExpr, i, len(results))
			matches = append(matches, PathMatch{
				Path:  concretePath,
				Value: v,
			})
		case nil:
			// Skip null values silently
			continue
		default:
			// Log warning for non-string types
			if !p.quiet {
				log.Printf("Warning: path %s[%d] contains %T, expected string, skipping",
					pathExpr, i, result)
			}
		}
	}

	return matches, nil
}

// Replace substitutes values in JSON data based on a replacement map.
// The map keys are concrete paths (e.g., "$.users[0].email") and values
// are the replacement strings.
func (p *Processor) Replace(jsonData []byte, replacements map[string]string) ([]byte, error) {
	if len(replacements) == 0 {
		return jsonData, nil
	}

	// Parse the JSON
	data, err := oj.Parse(jsonData)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Apply each replacement
	for pathExpr, newValue := range replacements {
		path, err := jp.ParseString(pathExpr)
		if err != nil {
			return nil, fmt.Errorf("invalid replacement path %q: %w", pathExpr, err)
		}

		// Set the new value
		if err := path.Set(data, newValue); err != nil {
			// Log warning but continue - the path might not exist in this row
			if !p.quiet {
				log.Printf("Warning: failed to set path %s: %v", pathExpr, err)
			}
		}
	}

	// Serialize back to JSON
	return oj.Marshal(data)
}

// ExtractAndCollect extracts values from multiple paths and returns them
// grouped by path expression. This is useful for processing multiple
// json_paths on a single JSON value.
func (p *Processor) ExtractAndCollect(jsonData []byte, pathExprs []string) (map[string][]PathMatch, error) {
	result := make(map[string][]PathMatch)

	for _, pathExpr := range pathExprs {
		matches, err := p.Extract(jsonData, pathExpr)
		if err != nil {
			return nil, err
		}
		if len(matches) > 0 {
			result[pathExpr] = matches
		}
	}

	return result, nil
}

// buildConcretePath converts a wildcard path to a concrete path with an index.
// For example, "$.users[*].email" with index 2 becomes "$.users[2].email"
func buildConcretePath(pathExpr string, index int, total int) string {
	// If there's only one result or no wildcard, return the original path
	if total == 1 {
		return pathExpr
	}

	// Replace the first [*] with the concrete index
	// This is a simplified approach - for deeply nested wildcards,
	// we might need more sophisticated tracking
	result := make([]byte, 0, len(pathExpr)+10)
	replaced := false

	for i := 0; i < len(pathExpr); i++ {
		if !replaced && i+2 < len(pathExpr) && pathExpr[i:i+3] == "[*]" {
			result = append(result, fmt.Sprintf("[%d]", index)...)
			i += 2 // Skip past [*]
			replaced = true
		} else {
			result = append(result, pathExpr[i])
		}
	}

	return string(result)
}
