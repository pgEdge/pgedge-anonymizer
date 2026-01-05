/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/pgedge/pgedge-anonymizer/internal/config"
	"github.com/pgedge/pgedge-anonymizer/internal/database"
	"github.com/pgedge/pgedge-anonymizer/internal/generator"
	"github.com/pgedge/pgedge-anonymizer/internal/pattern"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration without running",
	Long: `Validate the configuration file, patterns, and database schema
without performing any anonymization.

This command checks:
  - Configuration file syntax and required fields
  - Pattern file loading and pattern name validity
  - Database connectivity
  - Column existence in the database
  - Foreign key relationship analysis

Example:
  pgedge-anonymizer validate
  pgedge-anonymizer validate --config myconfig.yaml`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return runValidation()
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidation() error {
	// Check that a config file was loaded
	if err := CheckConfigLoaded(); err != nil {
		return err
	}

	fmt.Println("Validating configuration...")

	// Load configuration
	cfg, err := config.LoadFromViper()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}
	fmt.Println("  Configuration file: OK")

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation error: %w", err)
	}
	fmt.Println("  Configuration validation: OK")

	// Load patterns
	defaultPath := config.FindDefaultPatternsFile(cfg.Patterns.DefaultPath)
	if defaultPath == "" && !cfg.Patterns.DisableDefaults {
		fmt.Println("  Warning: default patterns file not found")
	}

	registry, err := pattern.LoadPatterns(
		defaultPath,
		cfg.Patterns.UserPath,
		cfg.Patterns.DisableDefaults,
	)
	if err != nil {
		return fmt.Errorf("pattern loading error: %w", err)
	}
	fmt.Printf("  Patterns loaded: %d\n", registry.Count())

	// Verify all configured patterns exist
	genMgr := generator.NewManager()
	for _, col := range cfg.Columns {
		if _, ok := registry.Get(col.Pattern); !ok {
			// Check if it's a built-in generator
			if _, ok := genMgr.Get(col.Pattern); !ok {
				return fmt.Errorf("unknown pattern %q for column %s",
					col.Pattern, col.Column)
			}
		}
	}
	fmt.Println("  Pattern references: OK")

	// Test database connection
	fmt.Println("\nValidating database connection...")
	connector := database.NewConnector(&cfg.Database)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := connector.Connect(ctx); err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}
	defer connector.Close()
	fmt.Println("  Database connection: OK")

	// Validate columns exist
	columns, err := cfg.GetColumnRefs()
	if err != nil {
		return fmt.Errorf("column parsing error: %w", err)
	}

	validator := database.NewSchemaValidator(connector.DB())
	missing, err := validator.ValidateColumns(ctx, columns)
	if err != nil {
		return fmt.Errorf("column validation error: %w", err)
	}
	if len(missing) > 0 {
		fmt.Println("\n  Missing columns:")
		for _, col := range missing {
			fmt.Printf("    - %s\n", col.String())
		}
		return fmt.Errorf("%d columns not found in database", len(missing))
	}
	fmt.Printf("  Column validation: OK (%d columns)\n", len(columns))

	// Analyze foreign keys
	fkAnalyzer := database.NewFKAnalyzer(connector.DB())
	fks, err := fkAnalyzer.Analyze(ctx, columns)
	if err != nil {
		return fmt.Errorf("foreign key analysis error: %w", err)
	}

	cascadeTargets, _ := fkAnalyzer.GetCascadeTargets(ctx, columns)

	fmt.Printf("\n  Foreign key relationships: %d\n", len(fks))
	if len(cascadeTargets) > 0 {
		fmt.Printf("  CASCADE targets (will be skipped): %d\n", len(cascadeTargets))
		for _, col := range cascadeTargets {
			fmt.Printf("    - %s\n", col.String())
		}
	}

	// Get processing order
	ordered, err := fkAnalyzer.GetProcessingOrder(ctx, columns)
	if err != nil {
		return fmt.Errorf("ordering error: %w", err)
	}

	fmt.Println("\n  Processing order:")
	for i, col := range ordered {
		skip := ""
		for _, ct := range cascadeTargets {
			if ct.String() == col.String() {
				skip = " (CASCADE - will skip)"
				break
			}
		}
		fmt.Printf("    %d. %s%s\n", i+1, col.String(), skip)
	}

	fmt.Println("\nValidation complete. Configuration is valid.")
	return nil
}
