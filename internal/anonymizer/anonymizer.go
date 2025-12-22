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
	"time"

	"github.com/pgedge/pgedge-anonymizer/internal/config"
	"github.com/pgedge/pgedge-anonymizer/internal/database"
	"github.com/pgedge/pgedge-anonymizer/internal/errors"
	"github.com/pgedge/pgedge-anonymizer/internal/generator"
	"github.com/pgedge/pgedge-anonymizer/internal/pattern"
	"github.com/pgedge/pgedge-anonymizer/internal/stats"
)

// Anonymizer orchestrates the complete anonymization process.
type Anonymizer struct {
	config     *config.Config
	patterns   *pattern.Registry
	generators *generator.Manager
	connector  *database.Connector
	dictionary *Dictionary
	quiet      bool
}

// Options configures the anonymizer.
type Options struct {
	Config       *config.Config
	Patterns     *pattern.Registry
	Quiet        bool
	BatchSize    int
	CacheSize    int
	DefaultsPath string
	UserPath     string
}

// New creates a new anonymizer with the given options.
func New(opts Options) (*Anonymizer, error) {
	// Create dictionary
	dict, err := NewDictionary(opts.CacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create dictionary: %w", err)
	}

	// Create generator manager
	genManager := generator.NewManager()

	// Register format patterns from the pattern registry
	if opts.Patterns != nil {
		if err := registerFormatPatterns(genManager, opts.Patterns); err != nil {
			return nil, fmt.Errorf("failed to register format patterns: %w", err)
		}
	}

	return &Anonymizer{
		config:     opts.Config,
		patterns:   opts.Patterns,
		generators: genManager,
		connector:  database.NewConnector(&opts.Config.Database),
		dictionary: dict,
		quiet:      opts.Quiet,
	}, nil
}

// registerFormatPatterns registers format-based generators from the pattern registry.
func registerFormatPatterns(mgr *generator.Manager, registry *pattern.Registry) error {
	for _, name := range registry.List() {
		p, _ := registry.Get(name)
		if p.IsFormatPattern() {
			cfg := generator.FormatPatternConfig{
				Name:    p.Name,
				Format:  p.Format,
				Type:    p.Type,
				Min:     p.Min,
				Max:     p.Max,
				MinYear: p.MinYear,
				MaxYear: p.MaxYear,
			}
			if err := mgr.RegisterFormatPattern(cfg); err != nil {
				return fmt.Errorf("failed to register pattern %s: %w", p.Name, err)
			}
		}
	}
	return nil
}

// Run executes the complete anonymization process.
func (a *Anonymizer) Run(ctx context.Context) (*stats.Stats, error) {
	defer a.dictionary.Close()

	// Connect to database
	if err := a.connector.Connect(ctx); err != nil {
		return nil, err
	}
	defer a.connector.Close()

	// Validate columns exist
	columns, err := a.config.GetColumnRefs()
	if err != nil {
		return nil, err
	}

	validator := database.NewSchemaValidator(a.connector.DB())
	missing, err := validator.ValidateColumns(ctx, columns)
	if err != nil {
		return nil, err
	}
	if len(missing) > 0 {
		return nil, errors.NewValidationError(
			"columns not found in database", missing)
	}

	// Analyze foreign keys and get processing order
	fkAnalyzer := database.NewFKAnalyzer(a.connector.DB())
	orderedColumns, err := fkAnalyzer.GetProcessingOrder(ctx, columns)
	if err != nil {
		return nil, err
	}

	// Get CASCADE targets to skip
	cascadeTargets, err := fkAnalyzer.GetCascadeTargets(ctx, columns)
	if err != nil {
		return nil, err
	}
	skipSet := make(map[string]bool)
	for _, col := range cascadeTargets {
		skipSet[col.String()] = true
	}

	// Build column-to-config mapping
	columnConfigMap := make(map[string]config.ColumnConfig)
	for _, cc := range a.config.Columns {
		columnConfigMap[cc.Column] = cc
	}

	// Start transaction
	tx, err := a.connector.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	// Ensure rollback on error
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// Process each column
	collector := stats.NewCollector()
	startTime := time.Now()

	for _, col := range orderedColumns {
		// Skip CASCADE targets
		if skipSet[col.String()] {
			if !a.quiet {
				fmt.Printf("Skipping %s (CASCADE target)\n", col.String())
			}
			continue
		}

		// Get column config
		colConfig, ok := columnConfigMap[col.String()]
		if !ok {
			return nil, fmt.Errorf("no config found for column %s", col.String())
		}

		// Get column data type for proper casting
		dataType, err := validator.GetColumnDataType(ctx, col)
		if err != nil {
			return nil, fmt.Errorf("failed to get data type for %s: %w",
				col.String(), err)
		}

		if !a.quiet {
			estimate, _ := validator.GetTableRowEstimate(ctx,
				col.Schema, col.Table)
			fmt.Printf("Processing %s (est. %d rows)...\n",
				col.String(), estimate)
		}

		// Process column - different handling for JSON vs simple columns
		colStart := time.Now()
		var result *ProcessResult

		if colConfig.IsJSONColumn() {
			// JSON column: process with JSON path extraction
			result, err = a.processJSONColumn(ctx, tx, col, dataType, colConfig)
		} else {
			// Simple column: process with single pattern
			result, err = a.processSimpleColumn(ctx, tx, col, dataType,
				colConfig.Pattern, validator)
		}

		if err != nil {
			return nil, errors.NewAnonymizationError(col, 0, "",
				fmt.Sprintf("processing failed: %v", err), err)
		}

		// Record statistics
		collector.RecordColumn(stats.ColumnStats{
			Column:           col,
			RowsProcessed:    result.RowsProcessed,
			ValuesAnonymized: result.ValuesAnonymized,
			UniqueValues:     result.UniqueValues,
			Duration:         time.Since(colStart),
		})

		if !a.quiet {
			fmt.Printf("  Completed: %d rows, %d values anonymized\n",
				result.RowsProcessed, result.ValuesAnonymized)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, errors.NewDatabaseError("commit",
			fmt.Sprintf("failed to commit transaction: %v", err), err)
	}
	committed = true

	// Finalize statistics
	finalStats := collector.Finalize(time.Since(startTime))

	return finalStats, nil
}

// Close releases resources held by the anonymizer.
func (a *Anonymizer) Close() error {
	if a.dictionary != nil {
		a.dictionary.Close()
	}
	if a.connector != nil {
		a.connector.Close()
	}
	return nil
}

// processSimpleColumn processes a column with a single pattern.
func (a *Anonymizer) processSimpleColumn(
	ctx context.Context,
	tx *sql.Tx,
	col errors.ColumnRef,
	dataType string,
	patternName string,
	validator *database.SchemaValidator,
) (*ProcessResult, error) {
	// Get generator for pattern
	gen, ok := a.generators.Get(patternName)
	if !ok {
		return nil, fmt.Errorf("unknown pattern %q for column %s",
			patternName, col.String())
	}

	// Check if column has a unique constraint
	hasUnique, err := validator.HasUniqueConstraint(ctx, col)
	if err != nil {
		return nil, fmt.Errorf("failed to check unique constraint for %s: %w",
			col.String(), err)
	}

	processor := NewColumnProcessor(tx, col, dataType, gen, a.dictionary,
		database.DefaultBatchSize, hasUnique)

	var lastProgress int64
	return processor.Process(ctx, func(processed int64) {
		if !a.quiet && processed-lastProgress >= 10000 {
			fmt.Printf("  %d rows processed\n", processed)
			lastProgress = processed
		}
	})
}

// processJSONColumn processes a JSON/JSONB column with multiple path patterns.
func (a *Anonymizer) processJSONColumn(
	ctx context.Context,
	tx *sql.Tx,
	col errors.ColumnRef,
	dataType string,
	colConfig config.ColumnConfig,
) (*ProcessResult, error) {
	// Build generator map for each JSON path
	generators := make(map[string]generator.Generator)
	for _, jp := range colConfig.JSONPaths {
		gen, ok := a.generators.Get(jp.Pattern)
		if !ok {
			return nil, fmt.Errorf("unknown pattern %q for JSON path %s in column %s",
				jp.Pattern, jp.Path, col.String())
		}
		generators[jp.Path] = gen
	}

	processor := NewJSONColumnProcessor(
		tx, col, dataType, colConfig.JSONPaths, generators,
		a.dictionary, database.DefaultBatchSize, a.quiet)

	var lastProgress int64
	return processor.Process(ctx, func(processed int64) {
		if !a.quiet && processed-lastProgress >= 10000 {
			fmt.Printf("  %d rows processed\n", processed)
			lastProgress = processed
		}
	})
}
