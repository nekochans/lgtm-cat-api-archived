package usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"time"

	"github.com/nekochans/lgtm-cat-api/domain"
)

type CreateLgtmImageUseCase struct {
	repository domain.S3Repository
	cdnDomain  string
}

func NewCreateLgtmImageUseCase(r domain.S3Repository, c string) *CreateLgtmImageUseCase {
	return &CreateLgtmImageUseCase{
		repository: r,
		cdnDomain:  c,
	}
}

type RequestBody struct {
	Image          string `json:"image"`
	ImageExtension string `json:"imageExtension"`
}

func (u *CreateLgtmImageUseCase) CreateLgtmImage(ctx context.Context, reqBody RequestBody) (*domain.UploadedLgtmImage, error) {

	if !domain.CanConvertImageExtension(reqBody.ImageExtension) {
		return nil, domain.ErrInvalidImageExtension
	}

	decodedImg, err := base64.StdEncoding.DecodeString(reqBody.Image)
	if err != nil {
		return nil, domain.ErrDecodeImage
	}

	buffer := new(bytes.Buffer)
	buffer.Write(decodedImg)

	prefix, err := domain.BuildS3Prefix(time.Now().UTC())
	if err != nil {
		return nil, domain.ErrTimeLoadLocation
	}

	imageName, err := domain.GenerateImageName()
	if err != nil {
		return nil, domain.ErrGenerateUuid
	}

	uploadS3param := domain.CreateUploadS3param(
		buffer,
		prefix,
		imageName,
		reqBody.ImageExtension,
	)

	err = u.repository.Upload(ctx, uploadS3param)
	if err != nil {
		return nil, domain.ErrUploadToS3
	}

	return domain.CreateUploadedLgtmImage(u.cdnDomain, prefix, imageName), nil
}
