// Package database manages database connection pooling, configuration,
// and health inspection for PostgreSQL 16+ with PostGIS.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Config encapsulates configuration for the PostgreSQL connection pool.
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DSN formats the connection string for PostgreSQL.
func (c Config) DSN() string {
	ssl := c.SSLMode
	if ssl == "" {
		ssl = "disable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, ssl)
}

// Pool wraps sql.DB to provide health checks and clean shutdown.
type Pool struct {
	db     *sql.DB
	config Config
}

// NewPool initializes the pool wrapper with configuration parameters.
func NewPool(db *sql.DB, cfg Config) *Pool {
	if db != nil {
		if cfg.MaxOpenConns > 0 {
			db.SetMaxOpenConns(cfg.MaxOpenConns)
		}
		if cfg.MaxIdleConns > 0 {
			db.SetMaxIdleConns(cfg.MaxIdleConns)
		}
		if cfg.ConnMaxLifetime > 0 {
			db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
		}
	}
	return &Pool{db: db, config: cfg}
}

// DB returns the underlying *sql.DB connection handle.
func (p *Pool) DB() *sql.DB {
	return p.db
}

// Ping checks database connectivity with a timeout context.
func (p *Pool) Ping(ctx context.Context) error {
	if p.db == nil {
		return errors.New("database: pool is not initialized")
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return p.db.PingContext(pingCtx)
}

// Close terminates all active pool connections gracefully.
func (p *Pool) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}
