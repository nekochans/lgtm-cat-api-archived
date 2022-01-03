package domain

import "errors"

type UploadedLgtmImage struct {
	Url string
}

var (
	ErrBadRequest            = errors.New("bad request")
	ErrInvalidImageExtension = errors.New("invalid image extension")
	ErrDecodeImage           = errors.New("failed to decode Base64 image")
	ErrGenerateUUID          = errors.New("failed to generate UUID")
	ErrTimeLoadLocation      = errors.New("failed to Time LoadLocation")
	ErrUploadToS3            = errors.New("failed to upload to S3")
)
