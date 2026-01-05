/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package errors defines custom error types for pgedge-anonymizer.
package errors

import (
	"fmt"
	"strings"
)

// ConfigError represents configuration-related errors.
type ConfigError struct {
	Path    string
	Message string
	Cause   error
}

func (e *ConfigError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("config error (%s): %s", e.Path, e.Message)
	}
	return fmt.Sprintf("config error: %s", e.Message)
}

func (e *ConfigError) Unwrap() error {
	return e.Cause
}

// NewConfigError creates a new ConfigError.
func NewConfigError(path, message string, cause error) *ConfigError {
	return &ConfigError{Path: path, Message: message, Cause: cause}
}

// PatternError represents pattern-related errors.
type PatternError struct {
	PatternName string
	Message     string
	Cause       error
}

func (e *PatternError) Error() string {
	if e.PatternName != "" {
		return fmt.Sprintf("pattern error (%s): %s", e.PatternName, e.Message)
	}
	return fmt.Sprintf("pattern error: %s", e.Message)
}

func (e *PatternError) Unwrap() error {
	return e.Cause
}

// NewPatternError creates a new PatternError.
func NewPatternError(name, message string, cause error) *PatternError {
	return &PatternError{PatternName: name, Message: message, Cause: cause}
}

// ColumnRef represents a fully-qualified column reference.
type ColumnRef struct {
	Schema string
	Table  string
	Column string
}

func (c ColumnRef) String() string {
	return fmt.Sprintf("%s.%s.%s", c.Schema, c.Table, c.Column)
}

// ParseColumnRef parses a schema.table.column string into a ColumnRef.
func ParseColumnRef(s string) (ColumnRef, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return ColumnRef{}, fmt.Errorf(
			"invalid column reference %q: expected schema.table.column format", s)
	}
	return ColumnRef{
		Schema: parts[0],
		Table:  parts[1],
		Column: parts[2],
	}, nil
}

// ValidationError represents validation failures for columns.
type ValidationError struct {
	Columns []ColumnRef
	Message string
}

func (e *ValidationError) Error() string {
	if len(e.Columns) == 0 {
		return fmt.Sprintf("validation error: %s", e.Message)
	}

	cols := make([]string, len(e.Columns))
	for i, c := range e.Columns {
		cols[i] = c.String()
	}
	return fmt.Sprintf("validation error: %s: %s", e.Message, strings.Join(cols, ", "))
}

// NewValidationError creates a new ValidationError.
func NewValidationError(message string, columns []ColumnRef) *ValidationError {
	return &ValidationError{Message: message, Columns: columns}
}

// DatabaseError represents database-related errors.
type DatabaseError struct {
	Operation string
	Column    *ColumnRef
	Message   string
	Cause     error
}

func (e *DatabaseError) Error() string {
	var sb strings.Builder
	sb.WriteString("database error")
	if e.Operation != "" {
		sb.WriteString(" during ")
		sb.WriteString(e.Operation)
	}
	if e.Column != nil {
		sb.WriteString(" on ")
		sb.WriteString(e.Column.String())
	}
	sb.WriteString(": ")
	sb.WriteString(e.Message)
	return sb.String()
}

func (e *DatabaseError) Unwrap() error {
	return e.Cause
}

// NewDatabaseError creates a new DatabaseError.
func NewDatabaseError(operation, message string, cause error) *DatabaseError {
	return &DatabaseError{Operation: operation, Message: message, Cause: cause}
}

// NewDatabaseErrorWithColumn creates a new DatabaseError with column context.
func NewDatabaseErrorWithColumn(operation string, col ColumnRef,
	message string, cause error) *DatabaseError {

	return &DatabaseError{
		Operation: operation,
		Column:    &col,
		Message:   message,
		Cause:     cause,
	}
}

// AnonymizationError represents errors during the anonymization process.
type AnonymizationError struct {
	Column  ColumnRef
	Row     int64
	Value   string
	Message string
	Cause   error
}

func (e *AnonymizationError) Error() string {
	var sb strings.Builder
	sb.WriteString("anonymization error")
	if e.Column.Schema != "" {
		sb.WriteString(" on ")
		sb.WriteString(e.Column.String())
	}
	if e.Row > 0 {
		sb.WriteString(fmt.Sprintf(" at row %d", e.Row))
	}
	sb.WriteString(": ")
	sb.WriteString(e.Message)
	return sb.String()
}

func (e *AnonymizationError) Unwrap() error {
	return e.Cause
}

// NewAnonymizationError creates a new AnonymizationError.
func NewAnonymizationError(col ColumnRef, row int64, value, message string,
	cause error) *AnonymizationError {

	return &AnonymizationError{
		Column:  col,
		Row:     row,
		Value:   value,
		Message: message,
		Cause:   cause,
	}
}
