/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package version provides build version information.
package version

// Version and BuildTime are set via ldflags at build time.
var (
	Version   = "1.0.0-alpha1"
	BuildTime = "unknown"
)
