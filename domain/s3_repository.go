package domain

type S3Repository interface {
	Upload(param *UploadS3param) error
}
