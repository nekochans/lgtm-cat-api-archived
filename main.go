package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/handler"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
	"github.com/nekochans/lgtm-cat-api/usecase"
)

var chiLambda *chiadapter.ChiLambda
var uploader *manager.Uploader
var region string
var uploadS3Bucket string
var lgtmImagesCdnDomain string
var q *db.Queries

func init() {
	region = os.Getenv("REGION")
	uploadS3Bucket = os.Getenv("UPLOAD_S3_BUCKET_NAME")
	lgtmImagesCdnDomain = os.Getenv("LGTM_IMAGES_CDN_DOMAIN")

	host := os.Getenv("DB_HOSTNAME")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
		log.Fatalln(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)

	dataSourceName := user + ":" + password + "@tcp(" + host + ")/" + dbName
	m, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	q = db.New(m)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	s3Repository := &infrastructure.S3Repository{Uploader: uploader, S3Bucket: uploadS3Bucket}
	createLgtmImageUseCase := &usecase.CreateLgtmImageUseCase{Repository: s3Repository, CdnDomain: lgtmImagesCdnDomain}
	createLgtmImageHandler := &handler.CreateLgtmImageHandler{
		CreateLgtmImageUseCase: createLgtmImageUseCase,
	}

	lgtmImageRepository := &infrastructure.LgtmImageRepository{Db: q}
	extractRandomImagesUseCase := &usecase.ExtractRandomImagesUseCase{Repository: lgtmImageRepository, CdnDomain: lgtmImagesCdnDomain}
	extractRandomImagesHandler := &handler.ExtractRandomImagesHandler{
		ExtractRandomImagesUseCase: extractRandomImagesUseCase,
	}

	if chiLambda == nil {
		r := chi.NewRouter()
		r.Post("/lgtm-images", createLgtmImageHandler.Create)
		r.Get("/lgtm-images", extractRandomImagesHandler.Extract)
		chiLambda = chiadapter.New(r)
	}
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}