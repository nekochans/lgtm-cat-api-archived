package usecase

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
)

type CreateLgtmImageUseCase struct {
	Repository     *infrastructure.S3Repository
	UploadS3Bucket string
	CdnDomain      string
}

type requestBody struct {
	Image          string `json:"image"`
	ImageExtension string `json:"imageExtension"`
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

func (u *CreateLgtmImageUseCase) CreateLgtmImage(req []byte) (*domain.UploadedLgtmImage, error) {

	var reqBody requestBody
	if err := json.Unmarshal(req, &reqBody); err != nil {
		return nil, domain.ErrBadRequest
	}

	if !canConvertImageExtension(reqBody.ImageExtension) {
		return nil, domain.ErrInvalidImageExtension
	}

	decodedImg, err := base64.StdEncoding.DecodeString(reqBody.Image)
	if err != nil {
		return nil, domain.ErrDecodeImage
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, domain.ErrGenerateUUID
	}

	buffer := new(bytes.Buffer)
	buffer.Write(decodedImg)

	prefix, err := buildS3Prefix(time.Now().UTC())
	if err != nil {
		return nil, domain.ErrTimeLoadLocation
	}

	imageName := uid.String()
	uploadKey := prefix + imageName + reqBody.ImageExtension

	err = u.Repository.Upload(
		u.UploadS3Bucket,
		buffer,
		decideS3ContentType(reqBody.ImageExtension),
		uploadKey,
	)

	if err != nil {
		return nil, domain.ErrUploadToS3
	}

	imageUrl := "https://" + u.CdnDomain + "/" + prefix + imageName + ".webp"
	return &domain.UploadedLgtmImage{Url: imageUrl}, nil
}
