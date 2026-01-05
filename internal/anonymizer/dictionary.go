/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package anonymizer provides the core anonymization logic.
package anonymizer

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	_ "modernc.org/sqlite" // SQLite driver
)

// DefaultCacheSize is the default number of entries in the LRU cache.
const DefaultCacheSize = 1000000 // 1 million entries

// Dictionary maintains consistent value mappings for anonymization.
// It uses a two-tier strategy:
//   - Tier 1: LRU in-memory cache for fast lookups
//   - Tier 2: SQLite disk cache for spillover when LRU evicts entries
//
// It also tracks reverse mappings (anonymized → original) to ensure
// uniqueness when columns have unique constraints.
type Dictionary struct {
	mu       sync.RWMutex
	cache    *lru.Cache[string, string]
	reverse  map[string]bool // tracks used anonymized values
	diskDB   *sql.DB
	diskPath string
}

// NewDictionary creates a new value dictionary.
func NewDictionary(cacheSize int) (*Dictionary, error) {
	if cacheSize <= 0 {
		cacheSize = DefaultCacheSize
	}

	cache, err := lru.New[string, string](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create LRU cache: %w", err)
	}

	d := &Dictionary{
		cache:   cache,
		reverse: make(map[string]bool),
	}

	// Initialize SQLite spillover database
	if err := d.initDiskCache(); err != nil {
		return nil, err
	}

	return d, nil
}

// initDiskCache creates a temporary SQLite database for spillover.
func (d *Dictionary) initDiskCache() error {
	// Create temp file for SQLite
	tmpDir := os.TempDir()
	d.diskPath = filepath.Join(tmpDir,
		fmt.Sprintf("pgedge-anon-%d.db", os.Getpid()))

	db, err := sql.Open("sqlite", d.diskPath)
	if err != nil {
		return fmt.Errorf("failed to open disk cache: %w", err)
	}

	// Create table for value mappings
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS mappings (
            original TEXT PRIMARY KEY,
            anonymized TEXT NOT NULL
        )
    `)
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create mappings table: %w", err)
	}

	// Create index for faster lookups
	_, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_original ON mappings(original)
    `)
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Create index on anonymized for reverse lookups (uniqueness checking)
	_, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_anonymized ON mappings(anonymized)
    `)
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to create anonymized index: %w", err)
	}

	d.diskDB = db
	return nil
}

// Get retrieves an anonymized value for the given original.
// Returns the anonymized value and true if found, empty string and false if not.
func (d *Dictionary) Get(original string) (string, bool) {
	d.mu.RLock()
	// Check LRU cache first (fast path)
	if val, ok := d.cache.Get(original); ok {
		d.mu.RUnlock()
		return val, true
	}
	d.mu.RUnlock()

	// Check disk cache
	d.mu.Lock()
	defer d.mu.Unlock()

	// Double-check LRU after acquiring write lock
	if val, ok := d.cache.Get(original); ok {
		return val, true
	}

	// Query disk cache
	var anonymized string
	err := d.diskDB.QueryRow(
		"SELECT anonymized FROM mappings WHERE original = ?",
		original,
	).Scan(&anonymized)

	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		// Log error but don't fail - treat as not found
		return "", false
	}

	// Promote to LRU cache
	d.cache.Add(original, anonymized)
	return anonymized, true
}

// Set stores a mapping from original to anonymized value.
func (d *Dictionary) Set(original, anonymized string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.setInternal(original, anonymized)
}

// setInternal stores a mapping (caller must hold lock).
func (d *Dictionary) setInternal(original, anonymized string) {
	// Add to LRU cache
	d.cache.Add(original, anonymized)

	// Track in reverse map
	d.reverse[anonymized] = true

	// Always store in disk cache for durability
	_, _ = d.diskDB.Exec(
		"INSERT OR REPLACE INTO mappings (original, anonymized) VALUES (?, ?)",
		original, anonymized,
	)
}

// IsUsed checks if an anonymized value is already in use.
func (d *Dictionary) IsUsed(anonymized string) bool {
	d.mu.RLock()
	// Check in-memory reverse map first
	if d.reverse[anonymized] {
		d.mu.RUnlock()
		return true
	}
	d.mu.RUnlock()

	// Check disk cache
	d.mu.Lock()
	defer d.mu.Unlock()

	// Double-check after acquiring write lock
	if d.reverse[anonymized] {
		return true
	}

	// Query disk cache
	var count int
	err := d.diskDB.QueryRow(
		"SELECT COUNT(*) FROM mappings WHERE anonymized = ?",
		anonymized,
	).Scan(&count)

	if err != nil {
		return false
	}

	if count > 0 {
		// Cache the result
		d.reverse[anonymized] = true
		return true
	}

	return false
}

// SetUnique stores a mapping only if the anonymized value is not already in use.
// Returns true if the mapping was stored, false if the anonymized value was
// already used by another original.
func (d *Dictionary) SetUnique(original, anonymized string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Check if this anonymized value is already used
	if d.reverse[anonymized] {
		// Check if it's used by the same original (that's ok)
		if existing, ok := d.cache.Get(original); ok && existing == anonymized {
			return true
		}
		return false
	}

	// Check disk cache for existing usage
	var existingOriginal string
	err := d.diskDB.QueryRow(
		"SELECT original FROM mappings WHERE anonymized = ?",
		anonymized,
	).Scan(&existingOriginal)

	if err == nil {
		// Found in disk - mark in reverse map
		d.reverse[anonymized] = true
		// It's ok if same original
		if existingOriginal == original {
			return true
		}
		return false
	}

	// Not used - safe to set
	d.setInternal(original, anonymized)
	return true
}

// Size returns the number of entries in the LRU cache.
func (d *Dictionary) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.cache.Len()
}

// PreloadUsedValues marks a list of values as already used in the reverse map.
// This is used to prevent generating values that already exist in the database
// (e.g., from a previous run or existing data).
func (d *Dictionary) PreloadUsedValues(values []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, v := range values {
		d.reverse[v] = true
	}
}

// DiskSize returns the number of entries in the disk cache.
func (d *Dictionary) DiskSize() (int64, error) {
	var count int64
	err := d.diskDB.QueryRow("SELECT COUNT(*) FROM mappings").Scan(&count)
	return count, err
}

// Close cleans up the dictionary resources.
func (d *Dictionary) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.diskDB != nil {
		d.diskDB.Close()
	}

	// Remove the temporary SQLite file
	if d.diskPath != "" {
		os.Remove(d.diskPath)
	}

	return nil
}
