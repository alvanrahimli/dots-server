package main

import (
	"database/sql"
	"log"
	"os"
)

func getDbInstance() *sql.DB {
	db, dbErr := sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	return db
}
