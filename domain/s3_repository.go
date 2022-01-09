package domain

import "context"

type S3Repository interface {
	Upload(context.Context, *UploadS3param) error
}
