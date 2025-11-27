/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package database provides PostgreSQL database connectivity and operations.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver

	"github.com/pgedge/pgedge-anonymizer/internal/config"
	"github.com/pgedge/pgedge-anonymizer/internal/errors"
)

// Connector manages database connections.
type Connector struct {
	db     *sql.DB
	config *config.DatabaseConfig
}

// NewConnector creates a new database connector.
func NewConnector(cfg *config.DatabaseConfig) *Connector {
	return &Connector{
		config: cfg,
	}
}

// Connect establishes a connection to the database.
func (c *Connector) Connect(ctx context.Context) error {
	connStr := c.config.ConnectionString()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return errors.NewDatabaseError("connect",
			fmt.Sprintf("failed to open database: %v", err), err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return errors.NewDatabaseError("connect",
			fmt.Sprintf("failed to ping database: %v", err), err)
	}

	c.db = db
	return nil
}

// Close closes the database connection.
func (c *Connector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// DB returns the underlying database connection.
func (c *Connector) DB() *sql.DB {
	return c.db
}

// BeginTx starts a new transaction.
func (c *Connector) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if c.db == nil {
		return nil, errors.NewDatabaseError("begin",
			"database connection not established", nil)
	}

	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, errors.NewDatabaseError("begin",
			fmt.Sprintf("failed to start transaction: %v", err), err)
	}

	return tx, nil
}
