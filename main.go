package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var chiLambda *chiadapter.ChiLambda
var uploader *manager.Uploader
var region string
var uploadS3Bucket string

func init() {
	region = os.Getenv("REGION")
	uploadS3Bucket = os.Getenv("UPLOAD_S3_BUCKET_NAME")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		// TODO ここでエラーが発生した場合、致命的な問題が起きているのでちゃんとしたログを出すように改修する
		log.Fatalln(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader = manager.NewUploader(s3Client)
}

type RequestBody struct {
	Image          string `json:"image"`
	ImageExtension string `json:"imageExtension"`
}

type ResponseBody struct {
	ImageURL string `json:"imageURL"`
}

type ResponseErrorBody struct {
	Message string `json:"message"`
}

func uploadToS3(
	ctx context.Context,
	uploader *manager.Uploader,
	bucket string,
	body *bytes.Buffer,
	contentType string,
	key string,
) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Body:        body,
		ContentType: aws.String(contentType),
		Key:         aws.String(key),
	}

	_, err := uploader.Upload(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func decideS3ContentType(ext string) string {
	contentType := ""

	switch ext {
	case ".png":
		contentType = "image/png"
	default:
		contentType = "image/jpeg"
	}

	return contentType
}

func RenderErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	resBody := &ResponseErrorBody{Message: message}
	resBodyJson, _ := json.Marshal(resBody)

	fmt.Fprint(w, string(resBodyJson))
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
}

func CreateLgtmImage(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RenderErrorResponse(w, 500, "Failed Read Request Body")
		return
	}

	// TODO 画像形式のバリデーション
	var reqBody RequestBody
	if err := json.Unmarshal(req, &reqBody); err != nil {
		RenderErrorResponse(w, 400, "Bad Reques")
		return
	}

	decodedImg, err := base64.StdEncoding.DecodeString(reqBody.Image)
	if err != nil {
		RenderErrorResponse(w, 500, "Failed Decode Base64 Image")
		return
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		RenderErrorResponse(w, 500, "Failed Generate UUID")
		return
	}

	buffer := new(bytes.Buffer)
	buffer.Write(decodedImg)

	// TODO YYYY/MM/DD/HH/UUIDV4 の形式に変更する
	uploadKey := "tmp/" + uid.String() + reqBody.ImageExtension
	ctx := context.Background()
	err = uploadToS3(
		ctx,
		uploader,
		uploadS3Bucket,
		buffer,
		decideS3ContentType(reqBody.ImageExtension),
		uploadKey,
	)

	if err != nil {
		RenderErrorResponse(w, 500, "Failed Upload To S3")
		return
	}

	// TODO ImageURLを生成した値に変更
	response := &ResponseBody{ImageURL: "https://lgtm-images.lgtmeow.com/2021/03/16/22/a66dd7da-3105-4806-8c58-6fc66a0a3d04.webp"}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
	w.WriteHeader(202)
	w.Header().Add("Content-Type", "application/json")

	return
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if chiLambda == nil {
		r := chi.NewRouter()
		r.Post("/lgtm-image", CreateLgtmImage)
		chiLambda = chiadapter.New(r)
	}
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
