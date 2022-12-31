package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/handler"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
)

var chiLambda *chiadapter.ChiLambda
var uploader *manager.Uploader
var queries *db.Queries
var logger infrastructure.Logger

func init() {
	queries = infrastructure.NewSqlcQueries()
	uploader = infrastructure.NewUploader()
	logger = infrastructure.NewLogger()
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if chiLambda == nil {
		r := handler.NewRouter(uploader, queries, logger)
		chiLambda = chiadapter.New(r)
	}
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
