package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/handler"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
)

var uploader *manager.Uploader
var queries *db.Queries
var logger infrastructure.Logger

func main() {
	queries = infrastructure.NewSqlcQueries()
	uploader = infrastructure.NewUploader()
	logger = infrastructure.NewLogger()

	r := handler.NewRouter(uploader, queries, logger)

	const timeoutSecond = 10
	server := &http.Server{
		Addr:              ":3333",
		Handler:           r,
		ReadHeaderTimeout: timeoutSecond * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}
