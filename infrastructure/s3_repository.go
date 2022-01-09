package infrastructure

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nekochans/lgtm-cat-api/domain"
)

type S3Repository struct {
	Uploader *manager.Uploader
	S3Bucket string
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

func (r *S3Repository) Upload(c context.Context, param *domain.UploadS3param) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	input := &s3.PutObjectInput{
		Bucket:      aws.String(r.S3Bucket),
		Body:        param.Body,
		ContentType: aws.String(decideS3ContentType(param.ImageExtension)),
		Key:         aws.String(param.Key),
	}

	_, err := r.Uploader.Upload(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
