package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var schema = ` CREATE table scheduler (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					date VARCHAR(8) NOT NULL DEFAULT "",
					title VARCHAR(8) NOT NULL DEFAULT "",
					comment TEXT NOT NULL DEFAULT "",
					repeat VARCHAR(8) NOT NULL DEFAULT "" );
				 CREATE INDEX scheduler_date_ind ON scheduler (date);
				`

const DateFormat = "20060102"

var db *sql.DB

func Close() {
	db.Close()
}

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)

	var install bool

	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)

	if err != nil {
		return fmt.Errorf("Database open error: %w", err)
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			db.Close()
			return fmt.Errorf("Schema creation error: %w", err)
		}

	}

	return nil

}
