package usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/nekochans/lgtm-cat-api/domain"
)

type CreateLgtmImageUseCase struct {
	Repository domain.S3Repository
	CdnDomain  string
}

type requestBody struct {
	Image          string `json:"image"`
	ImageExtension string `json:"imageExtension"`
}

func (u *CreateLgtmImageUseCase) CreateLgtmImage(ctx context.Context, req []byte) (*domain.UploadedLgtmImage, error) {

	var reqBody requestBody
	if err := json.Unmarshal(req, &reqBody); err != nil {
		return nil, domain.ErrBadRequest
	}

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

	err = u.Repository.Upload(ctx, uploadS3param)
	if err != nil {
		return nil, domain.ErrUploadToS3
	}

	return domain.CreateUploadedLgtmImage(u.CdnDomain, prefix, imageName), nil
}
