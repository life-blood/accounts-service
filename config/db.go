package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//Configured from .env configuration filed
const (
	dbUser = "DB_USER"
	dbPass = "DB_PASS"
	dbPort = "DB_PORT"
	dbName = "DB_NAME"
)

//CreateDatabaseConn retreieve database connection
func CreateDatabaseConn() (*sql.DB, error) {
	serverName := fmt.Sprintf("localhost:%s", os.Getenv(dbPort))
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv(dbUser),
		os.Getenv(dbPass), serverName, os.Getenv(dbName))

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
