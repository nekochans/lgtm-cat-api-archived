package domain

import (
	"context"
	"fmt"
)

type S3Repository interface {
	Upload(context.Context, *UploadS3param) error
}

type S3Error struct {
	Op  string
	Err error
}

func (e *S3Error) Error() string {
	return fmt.Sprintf("s3Repository: %s, %s", e.Op, e.Err)
}
