package test

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DbCreator struct{}

func (c *DbCreator) Create() (*sql.DB, error) {
	dbUser := os.Getenv("TEST_DB_USER")
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	dbHost := os.Getenv("TEST_DB_HOST")
	dbName := os.Getenv("TEST_DB_NAME")

	dbSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dbSourceName)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to test mysql server, %s", err)
	}

	return db, nil
}
