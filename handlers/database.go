package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Connecter struct {
	db *sql.DB
}

func OpenOrCreate(name string) (*Connecter, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %v", err)
	}

	dbFile := filepath.Join(currentDir, name)

	_, err = os.Stat(dbFile)
	var create bool
	if err != nil {
		create = true
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if create {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database: %v", err)
		}
	}

	return &Connecter{db: db}, nil
}

func (c *Connecter) Close() {
	defer c.db.Close()
}

func createDatabase(db *sql.DB) error {
	query := `CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT NOT NULL,
		repeat TEXT NOT NULL
	);
	CREATE INDEX date_index ON scheduler(date);`
	_, err := db.Exec(query)
	return err
}
