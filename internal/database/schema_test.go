/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------*/

package database

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetTableRowEstimate_handlesNegativeEstimate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	v := &SchemaValidator{db: db}

	// Simulate reltuples = -1 (should clamp to 0)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT GREATEST(COALESCE(c.reltuples, 0)::bigint, 0)
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = $1 AND c.relname = $2`)).
		WithArgs("public", "nostats").
		WillReturnRows(sqlmock.NewRows([]string{"?column?"}).AddRow(int64(0)))

	est, err := v.GetTableRowEstimate(context.Background(), "public", "nostats")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if est != 0 {
		t.Errorf("expected 0 for reltuples -1 (no stats), got %d", est)
	}

	// Simulate reltuples = 42 (normal case)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT GREATEST(COALESCE(c.reltuples, 0)::bigint, 0)
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE n.nspname = $1 AND c.relname = $2`)).
		WithArgs("public", "hasstats").
		WillReturnRows(sqlmock.NewRows([]string{"?column?"}).AddRow(int64(42)))

	est, err = v.GetTableRowEstimate(context.Background(), "public", "hasstats")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if est != 42 {
		t.Errorf("expected 42 for reltuples 42, got %d", est)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}
