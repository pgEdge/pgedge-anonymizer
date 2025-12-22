/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDatabaseConfigConnectionString tests connection string generation
func TestDatabaseConfigConnectionString(t *testing.T) {
	// Clear relevant env vars for testing
	origHost := os.Getenv("PGHOST")
	origPort := os.Getenv("PGPORT")
	origDB := os.Getenv("PGDATABASE")
	origUser := os.Getenv("PGUSER")
	origPass := os.Getenv("PGPASSWORD")
	origSSL := os.Getenv("PGSSLMODE")
	defer func() {
		os.Setenv("PGHOST", origHost)
		os.Setenv("PGPORT", origPort)
		os.Setenv("PGDATABASE", origDB)
		os.Setenv("PGUSER", origUser)
		os.Setenv("PGPASSWORD", origPass)
		os.Setenv("PGSSLMODE", origSSL)
	}()
	os.Unsetenv("PGHOST")
	os.Unsetenv("PGPORT")
	os.Unsetenv("PGDATABASE")
	os.Unsetenv("PGUSER")
	os.Unsetenv("PGPASSWORD")
	os.Unsetenv("PGSSLMODE")

	t.Run("full config", func(t *testing.T) {
		db := DatabaseConfig{
			Host:     "myhost",
			Port:     5433,
			Database: "mydb",
			User:     "myuser",
			Password: "mypass",
			SSLMode:  "require",
		}
		connStr := db.ConnectionString()

		expected := "host=myhost port=5433 dbname=mydb user=myuser " +
			"sslmode=require password=mypass"
		if connStr != expected {
			t.Errorf("expected %q, got %q", expected, connStr)
		}
	})

	t.Run("defaults applied", func(t *testing.T) {
		db := DatabaseConfig{
			Database: "testdb",
			User:     "testuser",
		}
		connStr := db.ConnectionString()

		// Should have default host=localhost, port=5432, sslmode=prefer
		if connStr != "host=localhost port=5432 dbname=testdb "+
			"user=testuser sslmode=prefer" {
			t.Errorf("unexpected connStr: %q", connStr)
		}
	})

	t.Run("env var fallback", func(t *testing.T) {
		os.Setenv("PGHOST", "envhost")
		os.Setenv("PGPORT", "5434")
		os.Setenv("PGDATABASE", "envdb")
		os.Setenv("PGUSER", "envuser")
		os.Setenv("PGPASSWORD", "envpass")
		os.Setenv("PGSSLMODE", "disable")

		db := DatabaseConfig{}
		connStr := db.ConnectionString()

		expected := "host=envhost port=5434 dbname=envdb user=envuser " +
			"sslmode=disable password=envpass"
		if connStr != expected {
			t.Errorf("expected %q, got %q", expected, connStr)
		}
	})

	t.Run("config overrides env", func(t *testing.T) {
		os.Setenv("PGHOST", "envhost")
		os.Setenv("PGUSER", "envuser")

		db := DatabaseConfig{
			Host:     "confighost",
			Database: "configdb",
			User:     "configuser",
		}
		connStr := db.ConnectionString()

		if !contains(connStr, "host=confighost") {
			t.Errorf("config host should override env: %q", connStr)
		}
		if !contains(connStr, "user=configuser") {
			t.Errorf("config user should override env: %q", connStr)
		}
	})

	t.Run("SSL cert paths", func(t *testing.T) {
		os.Unsetenv("PGHOST")
		os.Unsetenv("PGPORT")
		os.Unsetenv("PGSSLMODE")

		db := DatabaseConfig{
			Host:        "localhost",
			Port:        5432,
			Database:    "testdb",
			User:        "testuser",
			SSLMode:     "verify-full",
			SSLCert:     "/path/to/cert.pem",
			SSLKey:      "/path/to/key.pem",
			SSLRootCert: "/path/to/ca.pem",
		}
		connStr := db.ConnectionString()

		if !contains(connStr, "sslcert=/path/to/cert.pem") {
			t.Errorf("missing sslcert: %q", connStr)
		}
		if !contains(connStr, "sslkey=/path/to/key.pem") {
			t.Errorf("missing sslkey: %q", connStr)
		}
		if !contains(connStr, "sslrootcert=/path/to/ca.pem") {
			t.Errorf("missing sslrootcert: %q", connStr)
		}
	})
}

// TestConfigValidate tests configuration validation
func TestConfigValidate(t *testing.T) {
	// Clear env vars that might affect validation
	origDB := os.Getenv("PGDATABASE")
	origPGUser := os.Getenv("PGUSER")
	origUser := os.Getenv("USER")
	defer func() {
		os.Setenv("PGDATABASE", origDB)
		os.Setenv("PGUSER", origPGUser)
		os.Setenv("USER", origUser)
	}()
	os.Unsetenv("PGDATABASE")
	os.Unsetenv("PGUSER")
	os.Unsetenv("USER")

	t.Run("valid config", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				Database: "mydb",
				User:     "myuser",
			},
			Columns: []ColumnConfig{
				{Column: "public.users.email", Pattern: "EMAIL"},
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid config, got error: %v", err)
		}
	})

	t.Run("missing database", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				User: "myuser",
			},
			Columns: []ColumnConfig{
				{Column: "public.users.email", Pattern: "EMAIL"},
			},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for missing database")
		}
		if !contains(err.Error(), "database name is required") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing user", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				Database: "mydb",
			},
			Columns: []ColumnConfig{
				{Column: "public.users.email", Pattern: "EMAIL"},
			},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for missing user")
		}
		if !contains(err.Error(), "database user is required") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("no columns", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				Database: "mydb",
				User:     "myuser",
			},
			Columns: []ColumnConfig{},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for no columns")
		}
		if !contains(err.Error(), "at least one column") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid column format", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				Database: "mydb",
				User:     "myuser",
			},
			Columns: []ColumnConfig{
				{Column: "invalid_column", Pattern: "EMAIL"},
			},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for invalid column format")
		}
		if !contains(err.Error(), "schema.table.column format") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing pattern", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				Database: "mydb",
				User:     "myuser",
			},
			Columns: []ColumnConfig{
				{Column: "public.users.email", Pattern: ""},
			},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for missing pattern")
		}
		if !contains(err.Error(), "pattern name is required") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("env vars provide database and user", func(t *testing.T) {
		os.Setenv("PGDATABASE", "envdb")
		os.Setenv("PGUSER", "envuser")

		cfg := Config{
			Database: DatabaseConfig{},
			Columns: []ColumnConfig{
				{Column: "public.users.email", Pattern: "EMAIL"},
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid config with env vars, got: %v", err)
		}
	})
}

// TestConfigApplyOverrides tests CLI override application
func TestConfigApplyOverrides(t *testing.T) {
	cfg := Config{
		Database: DatabaseConfig{
			Host:     "original",
			Port:     5432,
			Database: "origdb",
			User:     "origuser",
		},
		Patterns: PatternsConfig{
			DefaultPath:     "/orig/default",
			UserPath:        "/orig/user",
			DisableDefaults: false,
		},
	}

	host := "newhost"
	port := 5433
	database := "newdb"
	user := "newuser"
	password := "newpass"
	defaultPatterns := "/new/default"
	userPatterns := "/new/user"
	disableDefaults := true

	overrides := CLIOverrides{
		Host:            &host,
		Port:            &port,
		Database:        &database,
		User:            &user,
		Password:        &password,
		DefaultPatterns: &defaultPatterns,
		UserPatterns:    &userPatterns,
		DisableDefaults: &disableDefaults,
	}

	cfg.ApplyOverrides(overrides)

	if cfg.Database.Host != "newhost" {
		t.Errorf("host not overridden: %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("port not overridden: %d", cfg.Database.Port)
	}
	if cfg.Database.Database != "newdb" {
		t.Errorf("database not overridden: %s", cfg.Database.Database)
	}
	if cfg.Database.User != "newuser" {
		t.Errorf("user not overridden: %s", cfg.Database.User)
	}
	if cfg.Database.Password != "newpass" {
		t.Errorf("password not overridden: %s", cfg.Database.Password)
	}
	if cfg.Patterns.DefaultPath != "/new/default" {
		t.Errorf("default patterns not overridden: %s",
			cfg.Patterns.DefaultPath)
	}
	if cfg.Patterns.UserPath != "/new/user" {
		t.Errorf("user patterns not overridden: %s", cfg.Patterns.UserPath)
	}
	if !cfg.Patterns.DisableDefaults {
		t.Error("disable defaults not overridden")
	}
}

// TestConfigLoad tests loading configuration from a file
func TestConfigLoad(t *testing.T) {
	t.Run("valid config file", func(t *testing.T) {
		content := `
database:
  host: testhost
  port: 5432
  database: testdb
  user: testuser

columns:
  - column: public.users.email
    pattern: EMAIL
`
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if cfg.Database.Host != "testhost" {
			t.Errorf("unexpected host: %s", cfg.Database.Host)
		}
		if cfg.Database.Database != "testdb" {
			t.Errorf("unexpected database: %s", cfg.Database.Database)
		}
		if len(cfg.Columns) != 1 {
			t.Errorf("unexpected column count: %d", len(cfg.Columns))
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := Load("/nonexistent/path/config.yaml")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		content := `
database:
  host: testhost
  port: [invalid yaml
`
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "invalid.yaml")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		_, err := Load(path)
		if err == nil {
			t.Error("expected error for invalid yaml")
		}
	})
}

// TestGetColumnRefs tests conversion of column configs to refs
func TestGetColumnRefs(t *testing.T) {
	cfg := Config{
		Columns: []ColumnConfig{
			{Column: "public.users.email", Pattern: "EMAIL"},
			{Column: "hr.employees.ssn", Pattern: "US_SSN"},
		},
	}

	refs, err := cfg.GetColumnRefs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(refs) != 2 {
		t.Errorf("expected 2 refs, got %d", len(refs))
	}

	if refs[0].Schema != "public" || refs[0].Table != "users" ||
		refs[0].Column != "email" {
		t.Errorf("unexpected first ref: %+v", refs[0])
	}

	if refs[1].Schema != "hr" || refs[1].Table != "employees" ||
		refs[1].Column != "ssn" {
		t.Errorf("unexpected second ref: %+v", refs[1])
	}
}

// TestFindDefaultPatternsFile tests pattern file search
func TestFindDefaultPatternsFile(t *testing.T) {
	t.Run("finds file in specified path", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "patterns.yaml")
		if err := os.WriteFile(path, []byte("patterns: []"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		found := FindDefaultPatternsFile(path)
		if found != path {
			t.Errorf("expected %s, got %s", path, found)
		}
	})

	t.Run("returns empty for nonexistent paths", func(t *testing.T) {
		found := FindDefaultPatternsFile("/nonexistent/path.yaml")
		// This may or may not find a file depending on whether
		// pgedge-anonymizer-patterns.yaml exists in cwd or exe dir
		// Just verify it doesn't panic
		_ = found
	})
}

// TestJSONColumnValidation tests JSON column configuration validation
func TestJSONColumnValidation(t *testing.T) {
	// Clear env vars that might affect validation
	origDB := os.Getenv("PGDATABASE")
	origPGUser := os.Getenv("PGUSER")
	origUser := os.Getenv("USER")
	defer func() {
		os.Setenv("PGDATABASE", origDB)
		os.Setenv("PGUSER", origPGUser)
		os.Setenv("USER", origUser)
	}()
	os.Setenv("PGDATABASE", "testdb")
	os.Setenv("PGUSER", "testuser")

	t.Run("valid JSON column config", func(t *testing.T) {
		cfg := Config{
			Columns: []ColumnConfig{
				{
					Column: "public.users.profile",
					JSONPaths: []JSONPathConfig{
						{Path: "$.email", Pattern: "EMAIL"},
						{Path: "$.phone", Pattern: "US_PHONE"},
					},
				},
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid config, got error: %v", err)
		}
	})

	t.Run("JSON column with pattern and json_paths", func(t *testing.T) {
		cfg := Config{
			Columns: []ColumnConfig{
				{
					Column:  "public.users.profile",
					Pattern: "EMAIL", // Should not be allowed with json_paths
					JSONPaths: []JSONPathConfig{
						{Path: "$.email", Pattern: "EMAIL"},
					},
				},
			},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for pattern with json_paths")
		}
		if !contains(err.Error(), "cannot specify both") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("JSON path missing path", func(t *testing.T) {
		cfg := Config{
			Columns: []ColumnConfig{
				{
					Column: "public.users.profile",
					JSONPaths: []JSONPathConfig{
						{Path: "", Pattern: "EMAIL"},
					},
				},
			},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for missing path")
		}
		if !contains(err.Error(), "path is required") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("JSON path missing pattern", func(t *testing.T) {
		cfg := Config{
			Columns: []ColumnConfig{
				{
					Column: "public.users.profile",
					JSONPaths: []JSONPathConfig{
						{Path: "$.email", Pattern: ""},
					},
				},
			},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for missing pattern")
		}
		if !contains(err.Error(), "pattern is required") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("JSON path not starting with $", func(t *testing.T) {
		cfg := Config{
			Columns: []ColumnConfig{
				{
					Column: "public.users.profile",
					JSONPaths: []JSONPathConfig{
						{Path: "email", Pattern: "EMAIL"},
					},
				},
			},
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("expected error for path not starting with $")
		}
		if !contains(err.Error(), "must start with '$'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("JSON path with array wildcard", func(t *testing.T) {
		cfg := Config{
			Columns: []ColumnConfig{
				{
					Column: "public.users.profile",
					JSONPaths: []JSONPathConfig{
						{Path: "$.contacts[*].email", Pattern: "EMAIL"},
					},
				},
			},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid config with array wildcard, got: %v", err)
		}
	})
}

// TestIsJSONColumn tests the IsJSONColumn helper method
func TestIsJSONColumn(t *testing.T) {
	t.Run("simple column", func(t *testing.T) {
		col := ColumnConfig{
			Column:  "public.users.email",
			Pattern: "EMAIL",
		}
		if col.IsJSONColumn() {
			t.Error("simple column should not be JSON column")
		}
	})

	t.Run("JSON column", func(t *testing.T) {
		col := ColumnConfig{
			Column: "public.users.profile",
			JSONPaths: []JSONPathConfig{
				{Path: "$.email", Pattern: "EMAIL"},
			},
		}
		if !col.IsJSONColumn() {
			t.Error("column with json_paths should be JSON column")
		}
	})

	t.Run("empty json_paths", func(t *testing.T) {
		col := ColumnConfig{
			Column:    "public.users.profile",
			Pattern:   "EMAIL",
			JSONPaths: []JSONPathConfig{},
		}
		if col.IsJSONColumn() {
			t.Error("column with empty json_paths should not be JSON column")
		}
	})
}

// helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsLoop(s, substr))
}

func containsLoop(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
