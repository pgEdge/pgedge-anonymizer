/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package pattern

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRegistry tests the pattern registry
func TestRegistry(t *testing.T) {
	t.Run("add and get", func(t *testing.T) {
		r := NewRegistry()

		p := Pattern{
			Name:        "TEST_PATTERN",
			Replacement: "XXX",
			Note:        "Test pattern",
		}

		if err := r.Add(p); err != nil {
			t.Fatalf("failed to add pattern: %v", err)
		}

		got, ok := r.Get("TEST_PATTERN")
		if !ok {
			t.Fatal("pattern not found")
		}

		if got.Name != "TEST_PATTERN" {
			t.Errorf("unexpected name: %s", got.Name)
		}
		if got.Replacement != "XXX" {
			t.Errorf("unexpected replacement: %s", got.Replacement)
		}
	})

	t.Run("case insensitive get", func(t *testing.T) {
		r := NewRegistry()
		_ = r.Add(Pattern{Name: "TEST_PATTERN", Replacement: "XXX"})

		_, ok := r.Get("test_pattern")
		if !ok {
			t.Error("case insensitive lookup failed")
		}

		_, ok = r.Get("Test_Pattern")
		if !ok {
			t.Error("mixed case lookup failed")
		}
	})

	t.Run("duplicate pattern error", func(t *testing.T) {
		r := NewRegistry()
		_ = r.Add(Pattern{Name: "DUPLICATE", Replacement: "X"})

		err := r.Add(Pattern{Name: "DUPLICATE", Replacement: "Y"})
		if err == nil {
			t.Error("expected error for duplicate pattern")
		}
	})

	t.Run("list patterns", func(t *testing.T) {
		r := NewRegistry()
		_ = r.Add(Pattern{Name: "A", Replacement: "X"})
		_ = r.Add(Pattern{Name: "B", Replacement: "Y"})
		_ = r.Add(Pattern{Name: "C", Replacement: "Z"})

		names := r.List()
		if len(names) != 3 {
			t.Errorf("expected 3 patterns, got %d", len(names))
		}
	})

	t.Run("count patterns", func(t *testing.T) {
		r := NewRegistry()
		if r.Count() != 0 {
			t.Errorf("expected 0, got %d", r.Count())
		}

		_ = r.Add(Pattern{Name: "A", Replacement: "X"})
		_ = r.Add(Pattern{Name: "B", Replacement: "Y"})

		if r.Count() != 2 {
			t.Errorf("expected 2, got %d", r.Count())
		}
	})

	t.Run("get nonexistent", func(t *testing.T) {
		r := NewRegistry()
		_, ok := r.Get("NONEXISTENT")
		if ok {
			t.Error("expected pattern not found")
		}
	})
}

// TestLoader tests pattern file loading
func TestLoader(t *testing.T) {
	loader := NewLoader()

	t.Run("load valid file", func(t *testing.T) {
		content := `
patterns:
  - name: TEST_ONE
    replacement: "XXX"
    note: "First test pattern"
  - name: TEST_TWO
    replacement: "YYY"
`
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "patterns.yaml")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		pf, err := loader.LoadFile(path)
		if err != nil {
			t.Fatalf("failed to load file: %v", err)
		}

		if len(pf.Patterns) != 2 {
			t.Errorf("expected 2 patterns, got %d", len(pf.Patterns))
		}

		if pf.Patterns[0].Name != "TEST_ONE" {
			t.Errorf("unexpected name: %s", pf.Patterns[0].Name)
		}
		if pf.Patterns[1].Name != "TEST_TWO" {
			t.Errorf("unexpected name: %s", pf.Patterns[1].Name)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := loader.LoadFile("/nonexistent/path.yaml")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		content := `
patterns:
  - name: [invalid
`
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "invalid.yaml")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		_, err := loader.LoadFile(path)
		if err == nil {
			t.Error("expected error for invalid yaml")
		}
	})

	t.Run("empty pattern name", func(t *testing.T) {
		content := `
patterns:
  - name: ""
    replacement: "XXX"
`
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "empty_name.yaml")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		_, err := loader.LoadFile(path)
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("empty replacement", func(t *testing.T) {
		content := `
patterns:
  - name: TEST
    replacement: ""
`
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "empty_repl.yaml")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		_, err := loader.LoadFile(path)
		if err == nil {
			t.Error("expected error for empty replacement")
		}
	})
}

// TestLoadToRegistry tests loading to registry
func TestLoadToRegistry(t *testing.T) {
	loader := NewLoader()

	content := `
patterns:
  - name: TEST_A
    replacement: "XXX"
  - name: TEST_B
    replacement: "YYY"
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "patterns.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	registry := NewRegistry()
	if err := loader.LoadToRegistry(path, registry); err != nil {
		t.Fatalf("failed to load to registry: %v", err)
	}

	if registry.Count() != 2 {
		t.Errorf("expected 2 patterns, got %d", registry.Count())
	}

	if _, ok := registry.Get("TEST_A"); !ok {
		t.Error("TEST_A not found")
	}
	if _, ok := registry.Get("TEST_B"); !ok {
		t.Error("TEST_B not found")
	}
}

// TestMergeToRegistry tests merging with conflict detection
func TestMergeToRegistry(t *testing.T) {
	loader := NewLoader()

	t.Run("merge without conflicts", func(t *testing.T) {
		// Create default patterns file
		defaultContent := `
patterns:
  - name: DEFAULT_A
    replacement: "XXX"
  - name: DEFAULT_B
    replacement: "YYY"
`
		// Create user patterns file
		userContent := `
patterns:
  - name: USER_A
    replacement: "ZZZ"
`
		tmpDir := t.TempDir()
		defaultPath := filepath.Join(tmpDir, "default.yaml")
		userPath := filepath.Join(tmpDir, "user.yaml")

		_ = os.WriteFile(defaultPath, []byte(defaultContent), 0644)
		_ = os.WriteFile(userPath, []byte(userContent), 0644)

		registry := NewRegistry()
		_ = loader.LoadToRegistry(defaultPath, registry)

		if err := loader.MergeToRegistry(userPath, registry); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if registry.Count() != 3 {
			t.Errorf("expected 3 patterns, got %d", registry.Count())
		}
	})

	t.Run("merge with conflicts", func(t *testing.T) {
		defaultContent := `
patterns:
  - name: CONFLICTING
    replacement: "DEFAULT"
`
		userContent := `
patterns:
  - name: CONFLICTING
    replacement: "USER"
`
		tmpDir := t.TempDir()
		defaultPath := filepath.Join(tmpDir, "default.yaml")
		userPath := filepath.Join(tmpDir, "user.yaml")

		_ = os.WriteFile(defaultPath, []byte(defaultContent), 0644)
		_ = os.WriteFile(userPath, []byte(userContent), 0644)

		registry := NewRegistry()
		_ = loader.LoadToRegistry(defaultPath, registry)

		err := loader.MergeToRegistry(userPath, registry)
		if err == nil {
			t.Error("expected conflict error")
		}
		if !contains(err.Error(), "conflict") {
			t.Errorf("expected conflict message, got: %v", err)
		}
	})
}

// TestLoadPatterns tests the main LoadPatterns function
func TestLoadPatterns(t *testing.T) {
	t.Run("defaults only", func(t *testing.T) {
		defaultContent := `
patterns:
  - name: DEFAULT_ONLY
    replacement: "XXX"
`
		tmpDir := t.TempDir()
		defaultPath := filepath.Join(tmpDir, "default.yaml")
		_ = os.WriteFile(defaultPath, []byte(defaultContent), 0644)

		registry, err := LoadPatterns(defaultPath, "", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if registry.Count() != 1 {
			t.Errorf("expected 1 pattern, got %d", registry.Count())
		}
		if _, ok := registry.Get("DEFAULT_ONLY"); !ok {
			t.Error("DEFAULT_ONLY not found")
		}
	})

	t.Run("user only with defaults disabled", func(t *testing.T) {
		userContent := `
patterns:
  - name: USER_ONLY
    replacement: "YYY"
`
		tmpDir := t.TempDir()
		userPath := filepath.Join(tmpDir, "user.yaml")
		_ = os.WriteFile(userPath, []byte(userContent), 0644)

		registry, err := LoadPatterns("", userPath, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if registry.Count() != 1 {
			t.Errorf("expected 1 pattern, got %d", registry.Count())
		}
		if _, ok := registry.Get("USER_ONLY"); !ok {
			t.Error("USER_ONLY not found")
		}
	})

	t.Run("defaults and user merged", func(t *testing.T) {
		defaultContent := `
patterns:
  - name: DEFAULT
    replacement: "XXX"
`
		userContent := `
patterns:
  - name: USER
    replacement: "YYY"
`
		tmpDir := t.TempDir()
		defaultPath := filepath.Join(tmpDir, "default.yaml")
		userPath := filepath.Join(tmpDir, "user.yaml")
		_ = os.WriteFile(defaultPath, []byte(defaultContent), 0644)
		_ = os.WriteFile(userPath, []byte(userContent), 0644)

		registry, err := LoadPatterns(defaultPath, userPath, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if registry.Count() != 2 {
			t.Errorf("expected 2 patterns, got %d", registry.Count())
		}
		if _, ok := registry.Get("DEFAULT"); !ok {
			t.Error("DEFAULT not found")
		}
		if _, ok := registry.Get("USER"); !ok {
			t.Error("USER not found")
		}
	})

	t.Run("empty paths", func(t *testing.T) {
		registry, err := LoadPatterns("", "", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if registry.Count() != 0 {
			t.Errorf("expected 0 patterns, got %d", registry.Count())
		}
	})
}

// helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
