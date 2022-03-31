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
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/handler"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
	"github.com/nekochans/lgtm-cat-api/usecase/createltgmimage"
	"github.com/nekochans/lgtm-cat-api/usecase/extractrandomimages"
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
	s3Repository := infrastructure.NewS3Repository(uploader, uploadS3Bucket)
	createLgtmImageUseCase := createltgmimage.NewUseCase(s3Repository, lgtmImagesCdnDomain)
	createLgtmImageHandler := handler.NewCreateLgtmImageHandler(createLgtmImageUseCase)

	lgtmImageRepository := infrastructure.NewLgtmImageRepository(q)
	extractRandomImagesUseCase := extractrandomimages.NewUseCase(lgtmImageRepository, lgtmImagesCdnDomain)
	extractRandomImagesHandler := handler.NewExtractRandomImagesHandler(extractRandomImagesUseCase)

	if chiLambda == nil {
		r := chi.NewRouter()

		maxAge := 300
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://localhost:2222"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           maxAge,
		}))

		r.Post("/lgtm-images", createLgtmImageHandler.Create)
		r.Get("/lgtm-images", extractRandomImagesHandler.Extract)
		chiLambda = chiadapter.New(r)
	}
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
