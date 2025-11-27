/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
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

// ForeignKey represents a foreign key relationship.
type ForeignKey struct {
	ConstraintName string
	ParentSchema   string
	ParentTable    string
	ParentColumn   string
	ChildSchema    string
	ChildTable     string
	ChildColumn    string
	OnUpdate       string // CASCADE, SET NULL, NO ACTION, etc.
	OnDelete       string
}

// FKAnalyzer analyzes foreign key relationships.
type FKAnalyzer struct {
	db *sql.DB
}

// NewFKAnalyzer creates a new foreign key analyzer.
func NewFKAnalyzer(db *sql.DB) *FKAnalyzer {
	return &FKAnalyzer{db: db}
}

// Analyze retrieves all foreign key relationships involving the given columns.
func (a *FKAnalyzer) Analyze(ctx context.Context,
	columns []errors.ColumnRef) ([]ForeignKey, error) {

	if len(columns) == 0 {
		return nil, nil
	}

	// Build a set of schema.table pairs to check
	tables := make(map[string]bool)
	for _, col := range columns {
		key := fmt.Sprintf("%s.%s", col.Schema, col.Table)
		tables[key] = true
	}

	// Query pg_constraint for foreign key relationships
	query := `
        SELECT
            c.conname AS constraint_name,
            pn.nspname AS parent_schema,
            pc.relname AS parent_table,
            pa.attname AS parent_column,
            cn.nspname AS child_schema,
            cc.relname AS child_table,
            ca.attname AS child_column,
            CASE c.confupdtype
                WHEN 'a' THEN 'NO ACTION'
                WHEN 'r' THEN 'RESTRICT'
                WHEN 'c' THEN 'CASCADE'
                WHEN 'n' THEN 'SET NULL'
                WHEN 'd' THEN 'SET DEFAULT'
                ELSE 'UNKNOWN'
            END AS on_update,
            CASE c.confdeltype
                WHEN 'a' THEN 'NO ACTION'
                WHEN 'r' THEN 'RESTRICT'
                WHEN 'c' THEN 'CASCADE'
                WHEN 'n' THEN 'SET NULL'
                WHEN 'd' THEN 'SET DEFAULT'
                ELSE 'UNKNOWN'
            END AS on_delete
        FROM pg_constraint c
        JOIN pg_class pc ON pc.oid = c.confrelid
        JOIN pg_namespace pn ON pn.oid = pc.relnamespace
        JOIN pg_class cc ON cc.oid = c.conrelid
        JOIN pg_namespace cn ON cn.oid = cc.relnamespace
        JOIN pg_attribute pa ON pa.attrelid = c.confrelid
            AND pa.attnum = ANY(c.confkey)
        JOIN pg_attribute ca ON ca.attrelid = c.conrelid
            AND ca.attnum = ANY(c.conkey)
        WHERE c.contype = 'f'
    `

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.NewDatabaseError("fk_analyze",
			fmt.Sprintf("failed to query foreign keys: %v", err), err)
	}
	defer rows.Close()

	var fks []ForeignKey
	for rows.Next() {
		var fk ForeignKey
		err := rows.Scan(
			&fk.ConstraintName,
			&fk.ParentSchema, &fk.ParentTable, &fk.ParentColumn,
			&fk.ChildSchema, &fk.ChildTable, &fk.ChildColumn,
			&fk.OnUpdate, &fk.OnDelete,
		)
		if err != nil {
			return nil, errors.NewDatabaseError("fk_analyze",
				fmt.Sprintf("failed to scan foreign key: %v", err), err)
		}

		// Only include FKs where at least one table is in our list
		parentKey := fmt.Sprintf("%s.%s", fk.ParentSchema, fk.ParentTable)
		childKey := fmt.Sprintf("%s.%s", fk.ChildSchema, fk.ChildTable)
		if tables[parentKey] || tables[childKey] {
			fks = append(fks, fk)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseError("fk_analyze",
			fmt.Sprintf("error iterating foreign keys: %v", err), err)
	}

	return fks, nil
}

// GetCascadeTargets returns columns that will be updated via CASCADE
// when their parent column is updated.
func (a *FKAnalyzer) GetCascadeTargets(ctx context.Context,
	columns []errors.ColumnRef) ([]errors.ColumnRef, error) {

	fks, err := a.Analyze(ctx, columns)
	if err != nil {
		return nil, err
	}

	// Build set of columns we're updating
	updating := make(map[string]bool)
	for _, col := range columns {
		updating[col.String()] = true
	}

	// Find child columns that will be cascaded
	var cascaded []errors.ColumnRef
	for _, fk := range fks {
		if fk.OnUpdate != "CASCADE" {
			continue
		}

		// If we're updating the parent column, the child is cascaded
		parentRef := errors.ColumnRef{
			Schema: fk.ParentSchema,
			Table:  fk.ParentTable,
			Column: fk.ParentColumn,
		}

		if updating[parentRef.String()] {
			cascaded = append(cascaded, errors.ColumnRef{
				Schema: fk.ChildSchema,
				Table:  fk.ChildTable,
				Column: fk.ChildColumn,
			})
		}
	}

	return cascaded, nil
}

// GetProcessingOrder returns the columns in an order that respects
// foreign key dependencies (parent before child for CASCADE).
func (a *FKAnalyzer) GetProcessingOrder(ctx context.Context,
	columns []errors.ColumnRef) ([]errors.ColumnRef, error) {

	fks, err := a.Analyze(ctx, columns)
	if err != nil {
		return nil, err
	}

	// Build dependency graph
	// For CASCADE: parent must be processed before child
	deps := make(map[string][]string) // child -> []parent
	colSet := make(map[string]errors.ColumnRef)

	for _, col := range columns {
		colSet[col.String()] = col
	}

	for _, fk := range fks {
		if fk.OnUpdate != "CASCADE" {
			continue
		}

		parentRef := errors.ColumnRef{
			Schema: fk.ParentSchema,
			Table:  fk.ParentTable,
			Column: fk.ParentColumn,
		}
		childRef := errors.ColumnRef{
			Schema: fk.ChildSchema,
			Table:  fk.ChildTable,
			Column: fk.ChildColumn,
		}

		// Only track if both are in our list
		if _, ok := colSet[parentRef.String()]; !ok {
			continue
		}
		if _, ok := colSet[childRef.String()]; !ok {
			continue
		}

		deps[childRef.String()] = append(deps[childRef.String()],
			parentRef.String())
	}

	// Topological sort
	var result []errors.ColumnRef
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	var visit func(key string) error
	visit = func(key string) error {
		if temp[key] {
			return fmt.Errorf("circular dependency detected at %s", key)
		}
		if visited[key] {
			return nil
		}

		temp[key] = true
		for _, dep := range deps[key] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		temp[key] = false
		visited[key] = true
		result = append(result, colSet[key])
		return nil
	}

	for key := range colSet {
		if err := visit(key); err != nil {
			return nil, errors.NewDatabaseError("ordering",
				err.Error(), nil)
		}
	}

	return result, nil
}
