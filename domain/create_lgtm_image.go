package domain

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidImageExtension = errors.New("invalid image extension")
)

type ErrDecodeImage struct {
	Err error
}

func (e *ErrDecodeImage) Error() string {
	return fmt.Sprintf("failed to decode Base64 image, %s", e.Err)
}

type ErrGenerateImageName struct {
	Err error
}

func (e *ErrGenerateImageName) Error() string {
	return fmt.Sprintf("failed to generate image name, %s", e.Err)
}

type ErrTimeLoadLocation struct {
	Err error
}

func (e *ErrTimeLoadLocation) Error() string {
	return fmt.Sprintf("failed to Time LoadLocation, %s", e.Err)
}

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
