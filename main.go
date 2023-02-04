package main

import (
	"errors"
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

	if err := infrastructure.InitSentry(); err != nil {
		logger.Error(err)
	}

	const timeoutSecond = 10
	server := &http.Server{
		Addr:              ":3333",
		Handler:           r,
		ReadHeaderTimeout: timeoutSecond * time.Second,
		ErrorLog:          log.New(&logForwarder{l: logger}, "", 0),
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(err)
		return
	}
}

type logForwarder struct {
	l infrastructure.Logger
}

func (fw *logForwarder) Write(p []byte) (int, error) {
	fw.l.Error(errors.New(string(p)))
	return len(p), nil
}
