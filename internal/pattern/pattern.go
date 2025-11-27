/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package pattern handles loading and management of anonymization patterns.
package pattern

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/pgedge/pgedge-anonymizer/internal/errors"
)

// Pattern represents an anonymization pattern definition.
type Pattern struct {
	Name        string `yaml:"name"`
	Replacement string `yaml:"replacement"`
	Note        string `yaml:"note,omitempty"`

	// Format-based pattern fields (optional)
	// When Format is set, a format generator is created instead of
	// using Replacement as a generator name.
	Format  string `yaml:"format,omitempty"`   // Format string (strftime, printf, or mask)
	Type    string `yaml:"type,omitempty"`     // Format type: "date", "number", or "mask"
	Min     int64  `yaml:"min,omitempty"`      // Minimum value for number type
	Max     int64  `yaml:"max,omitempty"`      // Maximum value for number type
	MinYear int    `yaml:"min_year,omitempty"` // Minimum year for date type
	MaxYear int    `yaml:"max_year,omitempty"` // Maximum year for date type
}

// IsFormatPattern returns true if this pattern uses format-based generation.
func (p Pattern) IsFormatPattern() bool {
	return p.Format != ""
}

// PatternFile represents the YAML file structure.
type PatternFile struct {
	Patterns []Pattern `yaml:"patterns"`
}

// Registry holds all loaded patterns indexed by name.
type Registry struct {
	patterns map[string]Pattern
}

// NewRegistry creates an empty pattern registry.
func NewRegistry() *Registry {
	return &Registry{
		patterns: make(map[string]Pattern),
	}
}

// Add adds a pattern to the registry.
func (r *Registry) Add(p Pattern) error {
	name := strings.ToUpper(p.Name)
	if _, exists := r.patterns[name]; exists {
		return errors.NewPatternError(p.Name, "pattern already exists", nil)
	}
	r.patterns[name] = p
	return nil
}

// Get retrieves a pattern by name (case-insensitive).
func (r *Registry) Get(name string) (Pattern, bool) {
	p, ok := r.patterns[strings.ToUpper(name)]
	return p, ok
}

// List returns all pattern names.
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.patterns))
	for name := range r.patterns {
		names = append(names, name)
	}
	return names
}

// Count returns the number of patterns in the registry.
func (r *Registry) Count() int {
	return len(r.patterns)
}

// Loader handles loading and merging pattern files.
type Loader struct{}

// NewLoader creates a new pattern loader.
func NewLoader() *Loader {
	return &Loader{}
}

// LoadFile loads patterns from a YAML file.
func (l *Loader) LoadFile(path string) (*PatternFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.NewPatternError("",
			fmt.Sprintf("failed to read pattern file %s", path), err)
	}

	var pf PatternFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, errors.NewPatternError("",
			fmt.Sprintf("failed to parse pattern file %s", path), err)
	}

	// Validate patterns
	for _, p := range pf.Patterns {
		if p.Name == "" {
			return nil, errors.NewPatternError("",
				fmt.Sprintf("pattern in %s has empty name", path), nil)
		}
		// Either Replacement OR Format must be specified
		if p.Replacement == "" && p.Format == "" {
			return nil, errors.NewPatternError(p.Name,
				"pattern must have either 'replacement' or 'format' field", nil)
		}
	}

	return &pf, nil
}

// LoadToRegistry loads patterns from a file into a registry.
func (l *Loader) LoadToRegistry(path string, registry *Registry) error {
	pf, err := l.LoadFile(path)
	if err != nil {
		return err
	}

	for _, p := range pf.Patterns {
		if err := registry.Add(p); err != nil {
			return err
		}
	}

	return nil
}

// MergeToRegistry merges patterns from a file into an existing registry,
// returning an error if any pattern names conflict.
func (l *Loader) MergeToRegistry(path string, registry *Registry) error {
	pf, err := l.LoadFile(path)
	if err != nil {
		return err
	}

	// Check for conflicts first
	var conflicts []string
	for _, p := range pf.Patterns {
		if _, exists := registry.Get(p.Name); exists {
			conflicts = append(conflicts, p.Name)
		}
	}

	if len(conflicts) > 0 {
		return errors.NewPatternError("",
			fmt.Sprintf("user patterns conflict with default patterns: %s",
				strings.Join(conflicts, ", ")), nil)
	}

	// No conflicts, add all patterns
	for _, p := range pf.Patterns {
		if err := registry.Add(p); err != nil {
			return err
		}
	}

	return nil
}

// LoadPatterns loads default and user patterns based on configuration.
func LoadPatterns(defaultPath, userPath string, disableDefaults bool) (
	*Registry, error) {

	registry := NewRegistry()
	loader := NewLoader()

	// Load default patterns unless disabled
	if !disableDefaults && defaultPath != "" {
		if err := loader.LoadToRegistry(defaultPath, registry); err != nil {
			return nil, err
		}
	}

	// Load user patterns if specified
	if userPath != "" {
		if disableDefaults {
			// No defaults loaded, so just load user patterns directly
			if err := loader.LoadToRegistry(userPath, registry); err != nil {
				return nil, err
			}
		} else {
			// Merge with conflict detection
			if err := loader.MergeToRegistry(userPath, registry); err != nil {
				return nil, err
			}
		}
	}

	return registry, nil
}
