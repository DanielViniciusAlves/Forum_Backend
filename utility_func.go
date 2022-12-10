package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
	// (*w).Header().Add("Access-Control-Allow-Credentials", "true")
	// (*w).Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	(*w).Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, OPTIONS")
	(*w).Header().Add("Content-Type", "application/json")
}

func dbConn() (db *sql.DB) {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "kiwizao123",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "comments",
		AllowNativePasswords: true,
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		println("Error while starting a connection with the database.")
		log.Fatal(err)
	}

	return
}
