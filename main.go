package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
)

var chiLambda *chiadapter.ChiLambda
var uploader *manager.Uploader
var region string
var uploadS3Bucket string
var lgtmImagesCdnDomain string

func init() {
	region = os.Getenv("REGION")
	uploadS3Bucket = os.Getenv("UPLOAD_S3_BUCKET_NAME")
	lgtmImagesCdnDomain = os.Getenv("LGTM_IMAGES_CDN_DOMAIN")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
		log.Fatalln(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if chiLambda == nil {
		r := chi.NewRouter()
		r.Post("/lgtm-images", CreateLgtmImage)
		chiLambda = chiadapter.New(r)
	}
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
