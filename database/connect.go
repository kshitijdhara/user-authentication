package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = ""
	dbname   = "saga"
)

var db *sql.DB

func connectDB() error {
	if host == "" || port == "" || user == "" || dbname == "" {
		return fmt.Errorf("database environment variables are not set properly")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		db = nil
		return fmt.Errorf("db.Ping: %w", err)
	}

	fmt.Println("Successfully connected!")
	return nil
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	err := db.Close()
	db = nil
	return err
}

func GetDatabaseClient() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	return nil, connectDB()
}
