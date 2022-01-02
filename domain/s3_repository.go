package domain

import "bytes"

type S3Repository interface {
	Upload(
		bucket string,
		body *bytes.Buffer,
		contentType string,
		key string,
	) error
}
