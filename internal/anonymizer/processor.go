/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package anonymizer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pgedge/pgedge-anonymizer/internal/database"
	"github.com/pgedge/pgedge-anonymizer/internal/errors"
	"github.com/pgedge/pgedge-anonymizer/internal/generator"
)

// maxCollisionRetries is the maximum number of times to retry generating
// a unique value when collisions occur.
const maxCollisionRetries = 100

// addUniqueSuffix adds a numeric suffix to make a value unique.
// For emails, it inserts the suffix before the @ sign.
// For other values, it appends the suffix at the end.
func addUniqueSuffix(value string, suffix int) string {
	if idx := strings.LastIndex(value, "@"); idx > 0 {
		// Email address - insert suffix before @
		return fmt.Sprintf("%s%d%s", value[:idx], suffix, value[idx:])
	}
	// Other values - append suffix
	return fmt.Sprintf("%s%d", value, suffix)
}

// ColumnProcessor processes a single column for anonymization.
type ColumnProcessor struct {
	tx                  *sql.Tx
	column              errors.ColumnRef
	dataType            string
	generator           generator.Generator
	dictionary          *Dictionary
	batchSize           int
	hasUniqueConstraint bool
}

// NewColumnProcessor creates a new column processor.
func NewColumnProcessor(
	tx *sql.Tx,
	column errors.ColumnRef,
	dataType string,
	gen generator.Generator,
	dict *Dictionary,
	batchSize int,
	hasUniqueConstraint bool,
) *ColumnProcessor {
	return &ColumnProcessor{
		tx:                  tx,
		column:              column,
		dataType:            dataType,
		generator:           gen,
		dictionary:          dict,
		batchSize:           batchSize,
		hasUniqueConstraint: hasUniqueConstraint,
	}
}

// ProcessResult contains statistics about column processing.
type ProcessResult struct {
	RowsProcessed    int64
	ValuesAnonymized int64
	UniqueValues     int64
}

// Process anonymizes all values in the column.
func (p *ColumnProcessor) Process(ctx context.Context,
	progress func(processed int64)) (*ProcessResult, error) {

	batch := database.NewBatchProcessor(p.tx, p.column, p.dataType, p.batchSize)

	// Open cursor
	if err := batch.OpenCursor(ctx); err != nil {
		return nil, err
	}
	defer func() { _ = batch.CloseCursor(ctx) }()

	result := &ProcessResult{}
	uniqueSeen := make(map[string]bool)

	for {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Fetch next batch
		rows, err := batch.FetchBatch(ctx)
		if err != nil {
			return nil, err
		}

		if len(rows) == 0 {
			break // No more rows
		}

		// Process batch
		updates := make(map[string]string)

		for _, row := range rows {
			// Skip empty values
			if row.Value == "" {
				continue
			}

			// Check dictionary for existing mapping
			anonymized, exists := p.dictionary.Get(row.Value)
			if !exists {
				// Generate new anonymized value
				anonymized = p.generator.Generate(row.Value)

				// For columns with unique constraints, use uniqueness checking
				// to avoid constraint violations. For other columns, just store
				// directly since duplicates are allowed.
				if p.hasUniqueConstraint {
					// Try to set with uniqueness check, retry with suffix if needed
					if !p.dictionary.SetUnique(row.Value, anonymized) {
						// Collision detected - retry with numeric suffix
						base := anonymized
						found := false
						for i := 1; i <= maxCollisionRetries; i++ {
							anonymized = addUniqueSuffix(base, i)
							if p.dictionary.SetUnique(row.Value, anonymized) {
								found = true
								break
							}
						}
						if !found {
							return nil, fmt.Errorf(
								"failed to generate unique value after %d attempts",
								maxCollisionRetries)
						}
					}
				} else {
					// No unique constraint: just store without uniqueness check
					p.dictionary.Set(row.Value, anonymized)
				}
				result.UniqueValues++
			}

			// Track unique values seen
			if !uniqueSeen[row.Value] {
				uniqueSeen[row.Value] = true
			}

			// Queue update
			updates[row.CTID] = anonymized
			result.ValuesAnonymized++
		}

		// Apply batch updates
		if len(updates) > 0 {
			if err := batch.UpdateBatch(ctx, updates); err != nil {
				return nil, err
			}
		}

		result.RowsProcessed += int64(len(rows))

		// Report progress
		if progress != nil {
			progress(result.RowsProcessed)
		}
	}

	return result, nil
}
