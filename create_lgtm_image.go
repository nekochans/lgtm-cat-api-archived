package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type RequestBody struct {
	Image          string `json:"image"`
	ImageExtension string `json:"imageExtension"`
}

type ResponseBody struct {
	ImageUrl string `json:"imageUrl"`
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

func canConvertImageExtension(ext string) bool {
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return false
	}
	return true
}

func buildS3Prefix(t time.Time) (string, error) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", err
	}
	timeTokyo := t.In(tokyo)
	return timeTokyo.Format("2006/01/02/15/"), nil
}

func CreateLgtmImage(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RenderErrorResponse(w, 500, "Failed Read Request Body")
		return
	}

	var reqBody RequestBody
	if err := json.Unmarshal(req, &reqBody); err != nil {
		RenderErrorResponse(w, 400, "Bad Request")
		return
	}

	if !canConvertImageExtension(reqBody.ImageExtension) {
		RenderErrorResponse(w, 422, "Invalid Image Extension")
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

	prefix, err := buildS3Prefix(time.Now().UTC())
	if err != nil {
		RenderErrorResponse(w, 500, "Failed Time LoadLocation")
		return
	}

	imageName := uid.String()
	uploadKey := prefix + imageName + reqBody.ImageExtension
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

	response := &ResponseBody{ImageUrl: "https://" + lgtmImagesCdnDomain + "/" + prefix + imageName + ".webp"}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
	w.WriteHeader(202)
	w.Header().Add("Content-Type", "application/json")

	return
}
