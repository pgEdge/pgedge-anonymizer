/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package config handles configuration loading and validation for
// pgedge-anonymizer.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/pgedge/pgedge-anonymizer/internal/errors"
)

// Config represents the complete application configuration.
type Config struct {
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Patterns PatternsConfig `yaml:"patterns" mapstructure:"patterns"`
	Columns  []ColumnConfig `yaml:"columns" mapstructure:"columns"`
}

// DatabaseConfig holds PostgreSQL connection parameters.
type DatabaseConfig struct {
	Host        string `yaml:"host" mapstructure:"host"`
	Port        int    `yaml:"port" mapstructure:"port"`
	Database    string `yaml:"database" mapstructure:"database"`
	User        string `yaml:"user" mapstructure:"user"`
	Password    string `yaml:"password,omitempty" mapstructure:"password"`
	SSLMode     string `yaml:"sslmode" mapstructure:"sslmode"`
	SSLCert     string `yaml:"sslcert,omitempty" mapstructure:"sslcert"`
	SSLKey      string `yaml:"sslkey,omitempty" mapstructure:"sslkey"`
	SSLRootCert string `yaml:"sslrootcert,omitempty" mapstructure:"sslrootcert"`
}

// PatternsConfig defines pattern file locations.
type PatternsConfig struct {
	DefaultPath     string `yaml:"default_path,omitempty" mapstructure:"default_path"`
	UserPath        string `yaml:"user_path,omitempty" mapstructure:"user_path"`
	DisableDefaults bool   `yaml:"disable_defaults" mapstructure:"disable_defaults"`
}

// ColumnConfig maps a database column to an anonymization pattern.
type ColumnConfig struct {
	Column  string `yaml:"column" mapstructure:"column"`
	Pattern string `yaml:"pattern" mapstructure:"pattern"`
}

// CLIOverrides represents command-line overrides for config.
type CLIOverrides struct {
	Host            *string
	Port            *int
	Database        *string
	User            *string
	Password        *string
	DefaultPatterns *string
	UserPatterns    *string
	DisableDefaults *bool
}

// ConnectionString returns a PostgreSQL connection string, falling back to
// libpq environment variables for missing values.
func (d *DatabaseConfig) ConnectionString() string {
	host := d.Host
	if host == "" {
		host = os.Getenv("PGHOST")
	}
	if host == "" {
		host = "localhost"
	}

	port := d.Port
	if port == 0 {
		if envPort := os.Getenv("PGPORT"); envPort != "" {
			_, _ = fmt.Sscanf(envPort, "%d", &port)
		}
	}
	if port == 0 {
		port = 5432
	}

	database := d.Database
	if database == "" {
		database = os.Getenv("PGDATABASE")
	}

	user := d.User
	if user == "" {
		user = os.Getenv("PGUSER")
	}
	if user == "" {
		user = os.Getenv("USER") // Fall back to OS user like libpq
	}

	password := d.Password
	if password == "" {
		password = os.Getenv("PGPASSWORD")
	}

	sslmode := d.SSLMode
	if sslmode == "" {
		sslmode = os.Getenv("PGSSLMODE")
	}
	if sslmode == "" {
		sslmode = "prefer"
	}

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s sslmode=%s",
		host, port, database, user, sslmode)

	if password != "" {
		connStr += fmt.Sprintf(" password=%s", password)
	}

	if d.SSLCert != "" {
		connStr += fmt.Sprintf(" sslcert=%s", d.SSLCert)
	}
	if d.SSLKey != "" {
		connStr += fmt.Sprintf(" sslkey=%s", d.SSLKey)
	}
	if d.SSLRootCert != "" {
		connStr += fmt.Sprintf(" sslrootcert=%s", d.SSLRootCert)
	}

	return connStr
}

// Load loads configuration from the specified file path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.NewConfigError(path, "failed to read config file", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, errors.NewConfigError(path, "failed to parse config file", err)
	}

	return &cfg, nil
}

// LoadFromViper loads configuration from viper settings.
func LoadFromViper() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.NewConfigError("", "failed to unmarshal config", err)
	}
	return &cfg, nil
}

// ApplyOverrides applies CLI overrides to the configuration.
func (c *Config) ApplyOverrides(overrides CLIOverrides) {
	if overrides.Host != nil {
		c.Database.Host = *overrides.Host
	}
	if overrides.Port != nil {
		c.Database.Port = *overrides.Port
	}
	if overrides.Database != nil {
		c.Database.Database = *overrides.Database
	}
	if overrides.User != nil {
		c.Database.User = *overrides.User
	}
	if overrides.Password != nil {
		c.Database.Password = *overrides.Password
	}
	if overrides.DefaultPatterns != nil {
		c.Patterns.DefaultPath = *overrides.DefaultPatterns
	}
	if overrides.UserPatterns != nil {
		c.Patterns.UserPath = *overrides.UserPatterns
	}
	if overrides.DisableDefaults != nil {
		c.Patterns.DisableDefaults = *overrides.DisableDefaults
	}
}

// Validate checks the configuration for completeness and correctness.
func (c *Config) Validate() error {
	var errs []string

	// Database validation - either config or env vars must provide these
	if c.Database.Database == "" && os.Getenv("PGDATABASE") == "" {
		errs = append(errs, "database name is required")
	}
	// User can come from config, PGUSER, or fall back to $USER (like libpq)
	if c.Database.User == "" && os.Getenv("PGUSER") == "" && os.Getenv("USER") == "" {
		errs = append(errs, "database user is required")
	}

	// Columns validation
	if len(c.Columns) == 0 {
		errs = append(errs, "at least one column must be specified")
	}

	for i, col := range c.Columns {
		if col.Column == "" {
			errs = append(errs, fmt.Sprintf("column[%d]: column name is required", i))
		} else {
			// Validate schema.table.column format
			parts := strings.Split(col.Column, ".")
			if len(parts) != 3 {
				errs = append(errs, fmt.Sprintf(
					"column[%d]: %q must be in schema.table.column format",
					i, col.Column))
			}
		}
		if col.Pattern == "" {
			errs = append(errs, fmt.Sprintf(
				"column[%d]: pattern name is required", i))
		}
	}

	if len(errs) > 0 {
		return errors.NewConfigError("", strings.Join(errs, "; "), nil)
	}

	return nil
}

// FindDefaultPatternsFile searches for the default patterns file in standard
// locations.
func FindDefaultPatternsFile(configPath string) string {
	// Search order per DESIGN.md:
	// 1. Path specified in config (configPath)
	// 2. /etc/pgedge
	// 3. Directory containing the binary

	searchPaths := []string{}

	if configPath != "" {
		searchPaths = append(searchPaths, configPath)
	}

	searchPaths = append(searchPaths, "/etc/pgedge/pgedge-anonymizer-patterns.yaml")

	if exe, err := os.Executable(); err == nil {
		searchPaths = append(searchPaths,
			filepath.Join(filepath.Dir(exe), "pgedge-anonymizer-patterns.yaml"))
	}

	// Also check current directory
	searchPaths = append(searchPaths, "pgedge-anonymizer-patterns.yaml")

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// GetColumnRefs converts ColumnConfig slice to ColumnRef slice.
func (c *Config) GetColumnRefs() ([]errors.ColumnRef, error) {
	refs := make([]errors.ColumnRef, len(c.Columns))
	for i, col := range c.Columns {
		ref, err := errors.ParseColumnRef(col.Column)
		if err != nil {
			return nil, err
		}
		refs[i] = ref
	}
	return refs, nil
}
