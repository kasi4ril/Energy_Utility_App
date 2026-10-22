// Package database_test validates database connection configuration and pooling helpers.
package database_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"disco-enterprise-platform/services/core-api/pkg/database"
)

func TestConfig_DSN(t *testing.T) {
	cfg := database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secretpassword",
		DBName:   "disco_enterprise",
		SSLMode:  "require",
	}

	dsn := cfg.DSN()
	if !strings.Contains(dsn, "host=localhost") || !strings.Contains(dsn, "port=5432") ||
		!strings.Contains(dsn, "dbname=disco_enterprise") || !strings.Contains(dsn, "sslmode=require") {
		t.Fatalf("unexpected DSN output: %s", dsn)
	}

	// Default SSL mode check
	cfgDefaultSSL := cfg
	cfgDefaultSSL.SSLMode = ""
	if !strings.Contains(cfgDefaultSSL.DSN(), "sslmode=disable") {
		t.Fatalf("expected sslmode=disable as default, got: %s", cfgDefaultSSL.DSN())
	}
}

func TestPool_NilDBHandling(t *testing.T) {
	cfg := database.Config{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 10 * time.Minute,
	}

	pool := database.NewPool(nil, cfg)
	if pool.DB() != nil {
		t.Fatalf("expected nil *sql.DB")
	}

	err := pool.Ping(context.Background())
	if err == nil {
		t.Fatalf("expected error pinging nil database handle")
	}

	if err := pool.Close(); err != nil {
		t.Fatalf("unexpected error closing pool: %v", err)
	}
}

