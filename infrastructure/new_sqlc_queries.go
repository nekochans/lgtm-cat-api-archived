package infrastructure

import (
	"database/sql"
	"log"
	"os"

	. "github.com/nekochans/lgtm-cat-api/db/sqlc"
)

func NewSqlcQueries() *Queries {
	host := os.Getenv("DB_HOSTNAME")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dataSourceName := user + ":" + password + "@tcp(" + host + ")/" + dbName + "?tls=true"
	m, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	return New(m)
}
