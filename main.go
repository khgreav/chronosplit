// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/khgreav/chronosplit/menus"

	tea "charm.land/bubbletea/v2"

	_ "modernc.org/sqlite"
)

var DBFile = "chronosplit.db?_pragma=foreign_keys(1)"

func getDBPath() (string, error) {
	if envPath := os.Getenv("CHRONOSPLIT_DB"); envPath != "" {
		return envPath, nil
	}
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("Unable to find user home dir: %w", err)
		}
		dataDir = filepath.Join(home, ".local", "share")
	}

	appDir := filepath.Join(dataDir, "chronosplit")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return "", fmt.Errorf("Unable to create data directory: %w", err)
	}
	return filepath.Join(appDir, DBFile), nil
}

func main() {
	dbPath, err := getDBPath()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	_, err = os.Stat(dbPath)
	exists := err == nil
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

	if !exists {
		tx, err := db.Begin()
		if err != nil {
			fmt.Printf("Failed to begin migration transaction: %v\n", err)
			os.Exit(1)
		}
		for _, migration := range Migrations {
			_, err := db.Exec(migration)
			if err != nil {
				fmt.Printf("Failed to apply migration: %v\n", err)
				_ = tx.Rollback()
				os.Exit(1)
			}
		}
		err = tx.Commit()
		if err != nil {
			fmt.Printf("Failed to commit migration changes: %v\n", err)
			os.Exit(1)
		}
	}

	app := tea.NewProgram(menus.NewMainMenu(db))
	if _, err := app.Run(); err != nil {
		fmt.Printf("Something is borked: %v\n", err)
		os.Exit(1)
	}
}
