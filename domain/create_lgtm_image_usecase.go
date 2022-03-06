package domain

import (
	"bytes"
	"errors"
	"time"
)

var (
	ErrBadRequest            = errors.New("bad request")
	ErrInvalidImageExtension = errors.New("invalid image extension")
	ErrDecodeImage           = errors.New("failed to decode Base64 image")
	ErrGenerateUuid          = errors.New("failed to generate UUID")
	ErrTimeLoadLocation      = errors.New("failed to Time LoadLocation")
	ErrUploadToS3            = errors.New("failed to upload to S3")
)

type UploadedLgtmImage struct {
	Url string
}

type UploadS3param struct {
	Body           *bytes.Buffer
	ImageExtension string
	Key            string
}

type LgtmImageObject struct {
	Id       int32
	Path     string
	Filename string
}

func CreateUploadedLgtmImage(domain, prefix, imageName string) *UploadedLgtmImage {
	return &UploadedLgtmImage{Url: "https://" + domain + "/" + prefix + imageName + ".webp"}
}

func CanConvertImageExtension(ext string) bool {
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return false
	}
	return true
}

func GenerateImageName(u UniqueIdGenerator) (string, error) {
	return u.Generate()
}

func BuildS3Prefix(t time.Time) (string, error) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", err
	}
	timeTokyo := t.In(tokyo)
	return timeTokyo.Format("2006/01/02/15/"), nil
}

func CreateUploadS3param(body *bytes.Buffer, prefix, imageName, imageExtension string) *UploadS3param {

	uploadKey := prefix + imageName + imageExtension

	return &UploadS3param{
		Body:           body,
		ImageExtension: imageExtension,
		Key:            uploadKey,
	}
}
