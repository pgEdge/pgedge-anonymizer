/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pgedge/pgedge-anonymizer/internal/anonymizer"
	"github.com/pgedge/pgedge-anonymizer/internal/config"
	"github.com/pgedge/pgedge-anonymizer/internal/pattern"
	"github.com/pgedge/pgedge-anonymizer/internal/stats"
)

var (
	// Database connection flags
	dbHost     string
	dbPort     int
	dbName     string
	dbUser     string
	dbPassword string

	// Pattern flags
	patternsPath string
	noDefaults   bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the anonymization process",
	Long: `Execute the anonymization process on the configured database.

This command connects to the PostgreSQL database, validates the configured
columns exist, analyzes foreign key relationships, and then anonymizes the
data within a single transaction.

Example:
  pgedge-anonymizer run
  pgedge-anonymizer run --config myconfig.yaml
  pgedge-anonymizer run --host localhost --database mydb --user admin`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return runAnonymization()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Database flags
	runCmd.Flags().StringVar(&dbHost, "host", "",
		"PostgreSQL host (overrides config)")
	runCmd.Flags().IntVar(&dbPort, "port", 0,
		"PostgreSQL port (overrides config)")
	runCmd.Flags().StringVar(&dbName, "database", "",
		"Database name (overrides config)")
	runCmd.Flags().StringVar(&dbUser, "user", "",
		"Database user (overrides config)")
	runCmd.Flags().StringVar(&dbPassword, "password", "",
		"Database password (overrides config)")

	// Pattern flags
	runCmd.Flags().StringVar(&patternsPath, "patterns", "",
		"Path to user patterns file")
	runCmd.Flags().BoolVar(&noDefaults, "no-defaults", false,
		"Disable default patterns")

	// Bind flags to viper
	_ = viper.BindPFlag("database.host", runCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("database.port", runCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("database.database", runCmd.Flags().Lookup("database"))
	_ = viper.BindPFlag("database.user", runCmd.Flags().Lookup("user"))
	_ = viper.BindPFlag("database.password", runCmd.Flags().Lookup("password"))
	_ = viper.BindPFlag("patterns.user_path", runCmd.Flags().Lookup("patterns"))
	_ = viper.BindPFlag("patterns.disable_defaults", runCmd.Flags().Lookup("no-defaults"))
}

func runAnonymization() error {
	// Load configuration
	cfg, err := config.LoadFromViper()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Apply CLI overrides
	overrides := config.CLIOverrides{}
	if dbHost != "" {
		overrides.Host = &dbHost
	}
	if dbPort != 0 {
		overrides.Port = &dbPort
	}
	if dbName != "" {
		overrides.Database = &dbName
	}
	if dbUser != "" {
		overrides.User = &dbUser
	}
	if dbPassword != "" {
		overrides.Password = &dbPassword
	}
	if patternsPath != "" {
		overrides.UserPatterns = &patternsPath
	}
	if noDefaults {
		overrides.DisableDefaults = &noDefaults
	}
	cfg.ApplyOverrides(overrides)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Load patterns
	defaultPath := config.FindDefaultPatternsFile(cfg.Patterns.DefaultPath)
	if defaultPath == "" && !cfg.Patterns.DisableDefaults {
		fmt.Fprintln(os.Stderr,
			"Warning: default patterns file not found")
	}

	registry, err := pattern.LoadPatterns(
		defaultPath,
		cfg.Patterns.UserPath,
		cfg.Patterns.DisableDefaults,
	)
	if err != nil {
		return fmt.Errorf("failed to load patterns: %w", err)
	}

	if !quiet {
		fmt.Printf("Loaded %d patterns\n", registry.Count())
		fmt.Printf("Processing %d columns\n", len(cfg.Columns))
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintln(os.Stderr, "\nReceived interrupt, cancelling...")
		cancel()
	}()

	// Create and run anonymizer
	anon, err := anonymizer.New(anonymizer.Options{
		Config:   cfg,
		Patterns: registry,
		Quiet:    quiet,
	})
	if err != nil {
		return fmt.Errorf("failed to create anonymizer: %w", err)
	}
	defer anon.Close()

	result, err := anon.Run(ctx)
	if err != nil {
		return fmt.Errorf("anonymization failed: %w", err)
	}

	// Report results
	reporter := stats.NewReporter()
	reporter.Report(result, os.Stdout)

	return nil
}
