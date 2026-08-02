package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/khgreav/chronosplit/menus"

	tea "charm.land/bubbletea/v2"

	_ "modernc.org/sqlite"
)

var DBFile = "app.db?_pragma=foreign_keys(1)"

func main() {
	_, err := os.Stat(DBFile)
	exists := err == nil
	db, err := sql.Open("sqlite", DBFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

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
				_ = tx.Rollback()
				os.Exit(1)
			}
		}
		err = tx.Commit()
		if err != nil {
			fmt.Printf("Failed to commit migration changes: %v", err)
			os.Exit(1)
		}
	}

	app := tea.NewProgram(menus.NewMainMenu(db))
	if _, err := app.Run(); err != nil {
		fmt.Printf("Something is borked: %v", err)
		os.Exit(1)
	}
}
