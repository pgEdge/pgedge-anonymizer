/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pgedge/pgedge-anonymizer/internal/errors"
)

// SchemaValidator validates column existence in the database.
type SchemaValidator struct {
	db *sql.DB
}

// NewSchemaValidator creates a new schema validator.
func NewSchemaValidator(db *sql.DB) *SchemaValidator {
	return &SchemaValidator{db: db}
}

// ValidateColumns checks that all specified columns exist in the database.
// Returns a list of columns that do NOT exist.
func (v *SchemaValidator) ValidateColumns(ctx context.Context,
	columns []errors.ColumnRef) ([]errors.ColumnRef, error) {

	if len(columns) == 0 {
		return nil, nil
	}

	// Query information_schema to check column existence
	query := `
        SELECT table_schema, table_name, column_name
        FROM information_schema.columns
        WHERE (table_schema, table_name, column_name) IN (
            SELECT * FROM unnest($1::text[], $2::text[], $3::text[])
        )
    `

	// Prepare arrays for the query
	schemas := make([]string, len(columns))
	tables := make([]string, len(columns))
	cols := make([]string, len(columns))

	for i, col := range columns {
		schemas[i] = col.Schema
		tables[i] = col.Table
		cols[i] = col.Column
	}

	rows, err := v.db.QueryContext(ctx, query, schemas, tables, cols)
	if err != nil {
		return nil, errors.NewDatabaseError("validate",
			fmt.Sprintf("failed to query columns: %v", err), err)
	}
	defer rows.Close()

	// Track which columns exist
	exists := make(map[string]bool)
	for rows.Next() {
		var schema, table, column string
		if err := rows.Scan(&schema, &table, &column); err != nil {
			return nil, errors.NewDatabaseError("validate",
				fmt.Sprintf("failed to scan column: %v", err), err)
		}
		key := fmt.Sprintf("%s.%s.%s", schema, table, column)
		exists[key] = true
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError("validate",
			fmt.Sprintf("error iterating columns: %v", err), err)
	}

	// Find columns that don't exist
	var missing []errors.ColumnRef
	for _, col := range columns {
		key := col.String()
		if !exists[key] {
			missing = append(missing, col)
		}
	}

	return missing, nil
}

// GetColumnDataType returns the data type of a column.
func (v *SchemaValidator) GetColumnDataType(ctx context.Context,
	col errors.ColumnRef) (string, error) {

	query := `
        SELECT data_type
        FROM information_schema.columns
        WHERE table_schema = $1
          AND table_name = $2
          AND column_name = $3
    `

	var dataType string
	err := v.db.QueryRowContext(ctx, query,
		col.Schema, col.Table, col.Column).Scan(&dataType)

	if err == sql.ErrNoRows {
		return "", errors.NewDatabaseError("get_type",
			fmt.Sprintf("column %s not found", col.String()), nil)
	}
	if err != nil {
		return "", errors.NewDatabaseError("get_type",
			fmt.Sprintf("failed to get column type: %v", err), err)
	}

	return dataType, nil
}

// GetTableRowEstimate returns an estimated row count for a table.
// This uses pg_class.reltuples for fast estimation without scanning.
func (v *SchemaValidator) GetTableRowEstimate(ctx context.Context,
	schema, table string) (int64, error) {

	query := `
        SELECT COALESCE(c.reltuples, 0)::bigint
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = $1 AND c.relname = $2
    `

	var estimate int64
	err := v.db.QueryRowContext(ctx, query, schema, table).Scan(&estimate)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, errors.NewDatabaseError("estimate",
			fmt.Sprintf("failed to get row estimate: %v", err), err)
	}

	return estimate, nil
}

// HasUniqueConstraint checks if a column has a unique constraint or is part
// of a unique index.
func (v *SchemaValidator) HasUniqueConstraint(ctx context.Context,
	col errors.ColumnRef) (bool, error) {

	// Check for unique constraints (including primary keys)
	query := `
        SELECT COUNT(*) > 0
        FROM pg_constraint c
        JOIN pg_class t ON t.oid = c.conrelid
        JOIN pg_namespace n ON n.oid = t.relnamespace
        JOIN pg_attribute a ON a.attrelid = t.oid
        WHERE n.nspname = $1
          AND t.relname = $2
          AND a.attname = $3
          AND c.contype IN ('u', 'p')
          AND a.attnum = ANY(c.conkey)
    `

	var hasConstraint bool
	err := v.db.QueryRowContext(ctx, query,
		col.Schema, col.Table, col.Column).Scan(&hasConstraint)
	if err != nil {
		return false, errors.NewDatabaseError("check_unique",
			fmt.Sprintf("failed to check unique constraint: %v", err), err)
	}

	if hasConstraint {
		return true, nil
	}

	// Also check for unique indexes not created via constraints
	indexQuery := `
        SELECT COUNT(*) > 0
        FROM pg_index i
        JOIN pg_class t ON t.oid = i.indrelid
        JOIN pg_class ix ON ix.oid = i.indexrelid
        JOIN pg_namespace n ON n.oid = t.relnamespace
        JOIN pg_attribute a ON a.attrelid = t.oid
        WHERE n.nspname = $1
          AND t.relname = $2
          AND a.attname = $3
          AND i.indisunique = true
          AND a.attnum = ANY(i.indkey)
    `

	err = v.db.QueryRowContext(ctx, indexQuery,
		col.Schema, col.Table, col.Column).Scan(&hasConstraint)
	if err != nil {
		return false, errors.NewDatabaseError("check_unique_index",
			fmt.Sprintf("failed to check unique index: %v", err), err)
	}

	return hasConstraint, nil
}

// GetDistinctValues returns all distinct non-null values from a column.
// Used to pre-load existing values for uniqueness checking.
func (v *SchemaValidator) GetDistinctValues(ctx context.Context,
	col errors.ColumnRef) ([]string, error) {

	query := fmt.Sprintf(`
        SELECT DISTINCT %s::text
        FROM %s.%s
        WHERE %s IS NOT NULL
    `,
		quoteIdentForSchema(col.Column),
		quoteIdentForSchema(col.Schema),
		quoteIdentForSchema(col.Table),
		quoteIdentForSchema(col.Column),
	)

	rows, err := v.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.NewDatabaseError("get_distinct",
			fmt.Sprintf("failed to get distinct values: %v", err), err)
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			return nil, errors.NewDatabaseError("get_distinct",
				fmt.Sprintf("failed to scan value: %v", err), err)
		}
		values = append(values, val)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError("get_distinct",
			fmt.Sprintf("error iterating values: %v", err), err)
	}

	return values, nil
}

// quoteIdentForSchema quotes an identifier for use in SQL.
func quoteIdentForSchema(s string) string {
	return `"` + s + `"`
}
