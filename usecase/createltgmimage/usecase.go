package createltgmimage

import (
	"bytes"
	"context"
	"encoding/base64"
	"time"

	"github.com/nekochans/lgtm-cat-api/derrors"
	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
)

type UseCase struct {
	repository  domain.S3Repository
	cdnDomain   string
	idGenerator domain.UniqueIdGenerator
}

func NewUseCase(r domain.S3Repository, c string) *UseCase {
	return &UseCase{
		repository:  r,
		cdnDomain:   c,
		idGenerator: &infrastructure.UuidGenerator{},
	}
}

type RequestBody struct {
	Image          string `json:"image"`
	ImageExtension string `json:"imageExtension"`
}

func (u *UseCase) CreateLgtmImage(
	ctx context.Context,
	reqBody RequestBody,
) (uploadedLgtmImage *domain.UploadedLgtmImage, err error) {
	defer derrors.Wrap(&err, "UseCase.CreateLgtmImage(reqBody.ImageExtension: %+v)", reqBody.ImageExtension)

	if !domain.CanConvertImageExtension(reqBody.ImageExtension) {
		return nil, domain.ErrInvalidImageExtension
	}

	decodedImg, err := base64.StdEncoding.DecodeString(reqBody.Image)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(decodedImg)

	prefix, err := domain.BuildS3Prefix(time.Now().UTC())
	if err != nil {
		return nil, err
	}

	imageName, err := domain.GenerateImageName(u.idGenerator)
	if err != nil {
		return nil, err
	}

	uploadS3param := domain.CreateUploadS3param(
		buffer,
		prefix,
		imageName,
		reqBody.ImageExtension,
	)

	err = u.repository.Upload(ctx, uploadS3param)
	if err != nil {
		return nil, err
	}

	return domain.CreateUploadedLgtmImage(u.cdnDomain, prefix, imageName), nil
}
