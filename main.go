package main

import (
	"chronosplit/menus"
	"database/sql"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	_ "modernc.org/sqlite"
)

var DB_FILE = "app.db?_pragma=foreign_keys(1)"

func main() {
	_, err := os.Stat(DB_FILE)
	exists := err == nil
	db, err := sql.Open("sqlite", DB_FILE)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if !exists {
		tx, err := db.Begin()
		if err != nil {
			fmt.Printf("Failed to begin migration transaction: %v", err)
			os.Exit(1)
		}
		for _, migration := range Migrations {
			_, err := db.Exec(migration)
			if err != nil {
				fmt.Printf("Failed to apply migration: %v", err)
				tx.Rollback()
				os.Exit(1)
			}
		}
		tx.Commit()
	}

	app := tea.NewProgram(menus.NewMainMenu(db))
	if _, err := app.Run(); err != nil {
		fmt.Printf("Something is borked: %v", err)
		os.Exit(1)
	}
}
