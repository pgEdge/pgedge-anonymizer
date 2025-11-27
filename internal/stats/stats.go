/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package stats provides statistics collection and reporting for
// anonymization operations.
package stats

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/pgedge/pgedge-anonymizer/internal/errors"
)

// ColumnStats holds statistics for a single column.
type ColumnStats struct {
	Column           errors.ColumnRef
	RowsProcessed    int64
	ValuesAnonymized int64
	UniqueValues     int64
	Duration         time.Duration
}

// Stats holds overall anonymization statistics.
type Stats struct {
	Columns         []ColumnStats
	TotalRows       int64
	TotalAnonymized int64
	TotalUnique     int64
	TotalDuration   time.Duration
}

// Collector collects statistics during processing.
type Collector struct {
	mu      sync.Mutex
	columns []ColumnStats
}

// NewCollector creates a new statistics collector.
func NewCollector() *Collector {
	return &Collector{
		columns: make([]ColumnStats, 0),
	}
}

// RecordColumn records statistics for a processed column.
func (c *Collector) RecordColumn(stats ColumnStats) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.columns = append(c.columns, stats)
}

// Finalize calculates totals and returns final statistics.
func (c *Collector) Finalize(totalDuration time.Duration) *Stats {
	c.mu.Lock()
	defer c.mu.Unlock()

	stats := &Stats{
		Columns:       c.columns,
		TotalDuration: totalDuration,
	}

	for _, col := range c.columns {
		stats.TotalRows += col.RowsProcessed
		stats.TotalAnonymized += col.ValuesAnonymized
		stats.TotalUnique += col.UniqueValues
	}

	return stats
}

// Reporter formats and displays statistics.
type Reporter struct{}

// NewReporter creates a new statistics reporter.
func NewReporter() *Reporter {
	return &Reporter{}
}

// Report generates a formatted report of the statistics.
func (r *Reporter) Report(stats *Stats, w io.Writer) {
	// Calculate the maximum column name width (minimum 20, for "Column" header)
	colWidth := 20
	for _, col := range stats.Columns {
		name := col.Column.String()
		if len(name) > colWidth {
			colWidth = len(name)
		}
	}
	// Also check "TOTAL" fits
	if len("TOTAL") > colWidth {
		colWidth = len("TOTAL")
	}

	// Fixed widths for numeric columns
	const numWidth = 10

	// Calculate total inner width: colWidth + 3 numeric columns + separators
	// Format: "║ {col} │ {rows} │ {values} │ {duration} ║"
	// Inner: 1 + colWidth + 3 + numWidth + 3 + numWidth + 3 + numWidth + 1
	innerWidth := 1 + colWidth + 3 + numWidth + 3 + numWidth + 3 + numWidth + 1

	// Build border strings
	topBorder := "╔" + strings.Repeat("═", innerWidth) + "╗"
	midBorder := "╠" + strings.Repeat("═", innerWidth) + "╣"
	botBorder := "╚" + strings.Repeat("═", innerWidth) + "╝"
	rowSep := "╟" + strings.Repeat("─", colWidth+2) + "┼" +
		strings.Repeat("─", numWidth+2) + "┼" +
		strings.Repeat("─", numWidth+2) + "┼" +
		strings.Repeat("─", numWidth+2) + "╢"

	// Center the title
	title := "Anonymization Summary"
	padding := innerWidth - len(title)
	leftPad := padding / 2
	rightPad := padding - leftPad
	titleLine := "║" + strings.Repeat(" ", leftPad) + title +
		strings.Repeat(" ", rightPad) + "║"

	// Header
	fmt.Fprintln(w)
	fmt.Fprintln(w, topBorder)
	fmt.Fprintln(w, titleLine)
	fmt.Fprintln(w, midBorder)

	// Column headers
	fmt.Fprintf(w, "║ %-*s │ %*s │ %*s │ %*s ║\n",
		colWidth, "Column", numWidth, "Rows", numWidth, "Values", numWidth, "Duration")
	fmt.Fprintln(w, rowSep)

	// Column rows
	for _, col := range stats.Columns {
		fmt.Fprintf(w, "║ %-*s │ %*d │ %*d │ %*s ║\n",
			colWidth, col.Column.String(),
			numWidth, col.RowsProcessed,
			numWidth, col.ValuesAnonymized,
			numWidth, formatDuration(col.Duration))
	}

	// Totals
	fmt.Fprintln(w, rowSep)
	fmt.Fprintf(w, "║ %-*s │ %*d │ %*d │ %*s ║\n",
		colWidth, "TOTAL",
		numWidth, stats.TotalRows,
		numWidth, stats.TotalAnonymized,
		numWidth, formatDuration(stats.TotalDuration))

	fmt.Fprintln(w, botBorder)

	// Additional info
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Columns processed: %d\n", len(stats.Columns))
	fmt.Fprintf(w, "Unique values anonymized: %d\n", stats.TotalUnique)
	fmt.Fprintf(w, "Total duration: %s\n", formatDuration(stats.TotalDuration))
}

// String returns a string representation of the statistics.
func (r *Reporter) String(stats *Stats) string {
	var sb strings.Builder
	r.Report(stats, &sb)
	return sb.String()
}

// formatDuration formats a duration for display.
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
