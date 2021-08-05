package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

var PORT = os.Getenv("SERVER_PORT")

// Set up logging
var (
	InfoLogger *log.Logger
	WarnLogger *log.Logger
	ErrLogger  *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	InfoLogger.Println("Starting server...")
	InfoLogger.Printf("Server started at port: %s", PORT)
	serverAddress := fmt.Sprintf(":%s", PORT)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/ping", pingHandler).Methods("GET")
	router.HandleFunc("/archives/{name}", getArchiveHandler).Methods("GET")

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/packages/add", addPackageHandler).Methods("POST")

	if err := http.ListenAndServe(serverAddress, router); err != nil {
		log.Fatal(err)
	}
}

func getDbInstance() *sql.DB {
	db, dbErr := sql.Open("sqlite3", "./dots.db")
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	return db
}
