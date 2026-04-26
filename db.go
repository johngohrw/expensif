package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dbDir := filepath.Join(home, ".expensif")
	os.MkdirAll(dbDir, 0755)

	dbPath := filepath.Join(dbDir, "expenses.db")

	db, err = sql.Open("sqlite", "file:"+dbPath+"?_fk=1")
	if err != nil {
		log.Fatal(err)
	}

	if err := migrate(); err != nil {
		log.Fatal(err)
	}
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		category TEXT NOT NULL,
		description TEXT,
		date TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category);
	CREATE INDEX IF NOT EXISTS idx_expenses_created_at ON expenses(created_at);
	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Migrate existing tables without date column
	var hasDate int
	db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('expenses') WHERE name = 'date'`).Scan(&hasDate)
	if hasDate == 0 {
		if _, err := db.Exec(`ALTER TABLE expenses ADD COLUMN date TEXT`); err != nil {
			return err
		}
		if _, err := db.Exec(`UPDATE expenses SET date = date('now') WHERE date IS NULL`); err != nil {
			return err
		}
	}

	// Preferences table
	prefsSchema := `
	CREATE TABLE IF NOT EXISTS preferences (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		currency TEXT DEFAULT 'USD',
		dark_mode INTEGER DEFAULT 0
	);
	INSERT OR IGNORE INTO preferences (id, currency, dark_mode) VALUES (1, 'USD', 0);
	`
	_, err := db.Exec(prefsSchema)
	return err
}

func closeDB() {
	if db != nil {
		db.Close()
	}
}
