package infrastructure

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	Uploader *manager.Uploader
}

func (r *S3Repository) Upload(
	bucket string,
	body *bytes.Buffer,
	contentType string,
	key string,
) error {
	ctx := context.Background()

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Body:        body,
		ContentType: aws.String(contentType),
		Key:         aws.String(key),
	}

	_, err := r.Uploader.Upload(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
