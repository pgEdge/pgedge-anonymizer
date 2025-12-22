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
	"log"

	"github.com/pgedge/pgedge-anonymizer/internal/config"
	"github.com/pgedge/pgedge-anonymizer/internal/database"
	"github.com/pgedge/pgedge-anonymizer/internal/errors"
	"github.com/pgedge/pgedge-anonymizer/internal/generator"
	"github.com/pgedge/pgedge-anonymizer/internal/jsonpath"
)

// JSONColumnProcessor processes a JSON/JSONB column for anonymization.
// It extracts values at specified JSON paths, anonymizes them, and
// updates the JSON with the anonymized values.
type JSONColumnProcessor struct {
	tx         *sql.Tx
	column     errors.ColumnRef
	dataType   string
	jsonPaths  []config.JSONPathConfig
	generators map[string]generator.Generator // path -> generator
	dictionary *Dictionary
	batchSize  int
	processor  *jsonpath.Processor
	quiet      bool
}

// NewJSONColumnProcessor creates a new JSON column processor.
func NewJSONColumnProcessor(
	tx *sql.Tx,
	column errors.ColumnRef,
	dataType string,
	jsonPaths []config.JSONPathConfig,
	generators map[string]generator.Generator,
	dict *Dictionary,
	batchSize int,
	quiet bool,
) *JSONColumnProcessor {
	return &JSONColumnProcessor{
		tx:         tx,
		column:     column,
		dataType:   dataType,
		jsonPaths:  jsonPaths,
		generators: generators,
		dictionary: dict,
		batchSize:  batchSize,
		processor:  jsonpath.NewProcessor(quiet),
		quiet:      quiet,
	}
}

// Process anonymizes all JSON values in the column at the specified paths.
func (p *JSONColumnProcessor) Process(ctx context.Context,
	progress func(processed int64)) (*ProcessResult, error) {

	batch := database.NewBatchProcessor(p.tx, p.column, p.dataType, p.batchSize)

	// Open cursor - for JSON columns we fetch the full JSON value
	if err := batch.OpenCursor(ctx); err != nil {
		return nil, err
	}
	defer func() { _ = batch.CloseCursor(ctx) }()

	result := &ProcessResult{}

	// Collect all path expressions for batch extraction
	pathExprs := make([]string, len(p.jsonPaths))
	for i, jp := range p.jsonPaths {
		pathExprs[i] = jp.Path
	}

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

			// Process this JSON value
			modifiedJSON, valuesAnonymized, err := p.processJSONValue(
				row.CTID, []byte(row.Value), pathExprs)
			if err != nil {
				// Log error but continue processing other rows
				if !p.quiet {
					log.Printf("Warning: failed to process JSON at %s (ctid=%s): %v",
						p.column, row.CTID, err)
				}
				continue
			}

			if valuesAnonymized > 0 {
				updates[row.CTID] = string(modifiedJSON)
				result.ValuesAnonymized += int64(valuesAnonymized)
			}
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

// processJSONValue extracts values at all paths, anonymizes them, and returns
// the modified JSON. Returns the modified JSON bytes and count of values anonymized.
func (p *JSONColumnProcessor) processJSONValue(
	ctid string,
	jsonData []byte,
	pathExprs []string,
) ([]byte, int, error) {

	// Extract all values at all paths
	allMatches, err := p.processor.ExtractAndCollect(jsonData, pathExprs)
	if err != nil {
		return nil, 0, err
	}

	if len(allMatches) == 0 {
		return jsonData, 0, nil // No matching paths in this JSON
	}

	// Build replacement map: concrete path -> anonymized value
	replacements := make(map[string]string)
	valuesAnonymized := 0

	for pathExpr, matches := range allMatches {
		gen, ok := p.generators[pathExpr]
		if !ok {
			continue // No generator for this path (shouldn't happen)
		}

		for _, match := range matches {
			// Check dictionary for existing mapping
			anonymized, exists := p.dictionary.Get(match.Value)
			if !exists {
				// Generate new anonymized value
				anonymized = gen.Generate(match.Value)
				p.dictionary.Set(match.Value, anonymized)
			}

			replacements[match.Path] = anonymized
			valuesAnonymized++
		}
	}

	if len(replacements) == 0 {
		return jsonData, 0, nil
	}

	// Apply all replacements to the JSON
	modifiedJSON, err := p.processor.Replace(jsonData, replacements)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to replace values: %w", err)
	}

	return modifiedJSON, valuesAnonymized, nil
}
