// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package tests

import (
	"database/sql"
	"testing"

	"github.com/khgreav/chronosplit/internal/migrations"
	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T, schema string) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if schema != "" {
		_, err = db.Exec(schema)
		if err != nil {
			t.Fatalf("failed to execute schema: %v", err)
		}
	} else {
		for _, migration := range migrations.Migrations {
			if _, err := db.Exec(migration); err != nil {
				t.Fatalf("failed to apply migrations: %v", err)
			}
		}
	}

	return db
}
