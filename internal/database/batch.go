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
	"strings"

	"github.com/pgedge/pgedge-anonymizer/internal/errors"
)

// DefaultBatchSize is the default number of rows to process in a batch.
const DefaultBatchSize = 10000

// RowData represents a row fetched for processing.
type RowData struct {
	CTID  string // PostgreSQL physical row ID
	Value string // The column value to anonymize
}

// BatchProcessor handles batch reading and writing for a table column.
type BatchProcessor struct {
	tx        *sql.Tx
	column    errors.ColumnRef
	dataType  string
	batchSize int

	// Cursor state
	cursorName string
	cursorOpen bool
}

// NewBatchProcessor creates a new batch processor.
func NewBatchProcessor(tx *sql.Tx, col errors.ColumnRef, dataType string,
	batchSize int) *BatchProcessor {
	if batchSize <= 0 {
		batchSize = DefaultBatchSize
	}

	return &BatchProcessor{
		tx:         tx,
		column:     col,
		dataType:   dataType,
		batchSize:  batchSize,
		cursorName: fmt.Sprintf("anon_%s_%s_%s", col.Schema, col.Table, col.Column),
	}
}

// OpenCursor declares a server-side cursor for reading rows.
func (p *BatchProcessor) OpenCursor(ctx context.Context) error {
	// Use ctid for efficient updates
	query := fmt.Sprintf(
		`DECLARE %s CURSOR FOR
         SELECT ctid::text, %s::text
         FROM %s.%s
         WHERE %s IS NOT NULL`,
		p.cursorName,
		quoteIdent(p.column.Column),
		quoteIdent(p.column.Schema),
		quoteIdent(p.column.Table),
		quoteIdent(p.column.Column),
	)

	_, err := p.tx.ExecContext(ctx, query)
	if err != nil {
		return errors.NewDatabaseErrorWithColumn("cursor_open", p.column,
			fmt.Sprintf("failed to declare cursor: %v", err), err)
	}

	p.cursorOpen = true
	return nil
}

// FetchBatch fetches the next batch of rows from the cursor.
func (p *BatchProcessor) FetchBatch(ctx context.Context) ([]RowData, error) {
	if !p.cursorOpen {
		return nil, errors.NewDatabaseErrorWithColumn("fetch", p.column,
			"cursor not open", nil)
	}

	query := fmt.Sprintf("FETCH %d FROM %s", p.batchSize, p.cursorName)
	rows, err := p.tx.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.NewDatabaseErrorWithColumn("fetch", p.column,
			fmt.Sprintf("failed to fetch from cursor: %v", err), err)
	}
	defer rows.Close()

	var batch []RowData
	for rows.Next() {
		var rd RowData
		if err := rows.Scan(&rd.CTID, &rd.Value); err != nil {
			return nil, errors.NewDatabaseErrorWithColumn("fetch", p.column,
				fmt.Sprintf("failed to scan row: %v", err), err)
		}
		batch = append(batch, rd)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewDatabaseErrorWithColumn("fetch", p.column,
			fmt.Sprintf("error iterating rows: %v", err), err)
	}

	return batch, nil
}

// CloseCursor closes the server-side cursor.
func (p *BatchProcessor) CloseCursor(ctx context.Context) error {
	if !p.cursorOpen {
		return nil
	}

	_, err := p.tx.ExecContext(ctx, fmt.Sprintf("CLOSE %s", p.cursorName))
	if err != nil {
		return errors.NewDatabaseErrorWithColumn("cursor_close", p.column,
			fmt.Sprintf("failed to close cursor: %v", err), err)
	}

	p.cursorOpen = false
	return nil
}

// UpdateRow updates a single row by CTID.
func (p *BatchProcessor) UpdateRow(ctx context.Context, ctid, newValue string) error {
	query := fmt.Sprintf(
		`UPDATE %s.%s SET %s = $1 WHERE ctid = $2::tid`,
		quoteIdent(p.column.Schema),
		quoteIdent(p.column.Table),
		quoteIdent(p.column.Column),
	)

	_, err := p.tx.ExecContext(ctx, query, newValue, ctid)
	if err != nil {
		return errors.NewDatabaseErrorWithColumn("update", p.column,
			fmt.Sprintf("failed to update row %s: %v", ctid, err), err)
	}

	return nil
}

// UpdateBatch updates multiple rows by CTID in a single statement.
func (p *BatchProcessor) UpdateBatch(ctx context.Context,
	updates map[string]string) error {

	if len(updates) == 0 {
		return nil
	}

	// Build arrays for unnest
	ctids := make([]string, 0, len(updates))
	values := make([]string, 0, len(updates))
	for ctid, value := range updates {
		ctids = append(ctids, ctid)
		values = append(values, value)
	}

	// Build the value expression with appropriate type cast
	valueExpr := "u.new_value"
	if p.dataType != "" && p.dataType != "text" &&
		p.dataType != "character varying" && p.dataType != "character" {
		// Cast to the column's actual type for non-text columns
		valueExpr = fmt.Sprintf("u.new_value::%s", p.dataType)
	}

	// Use UPDATE FROM with unnest for efficient batch updates
	query := fmt.Sprintf(`
        UPDATE %s.%s t
        SET %s = %s
        FROM (
            SELECT unnest($1::tid[]) AS ctid, unnest($2::text[]) AS new_value
        ) u
        WHERE t.ctid = u.ctid`,
		quoteIdent(p.column.Schema),
		quoteIdent(p.column.Table),
		quoteIdent(p.column.Column),
		valueExpr,
	)

	_, err := p.tx.ExecContext(ctx, query, ctids, values)
	if err != nil {
		return errors.NewDatabaseErrorWithColumn("batch_update", p.column,
			fmt.Sprintf("failed to batch update: %v", err), err)
	}

	return nil
}

// quoteIdent quotes a PostgreSQL identifier to prevent SQL injection.
func quoteIdent(s string) string {
	// Replace any double quotes with two double quotes
	escaped := strings.ReplaceAll(s, `"`, `""`)
	return `"` + escaped + `"`
}
