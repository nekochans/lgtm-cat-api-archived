package createltgmimage

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/nekochans/lgtm-cat-api/domain"
)

type mockUniqueIdGenerator struct {
	domain.UniqueIdGenerator
	FakeGenerate func() (string, error)
}

func (d *mockUniqueIdGenerator) Generate() (string, error) {
	return d.FakeGenerate()
}

type mockS3Repository struct {
	domain.S3Repository
	FakeUpload func(context.Context, *domain.UploadS3param) error
}

func (d *mockS3Repository) Upload(c context.Context, u *domain.UploadS3param) error {
	return d.FakeUpload(c, u)
}

//nolint:funlen
func TestCreateLgtmImage(t *testing.T) {
	imageName := "test-image-name"
	cdnDomain := "lgtm-images.lgtmeow.com"

	t.Run("Success create LGTM image", func(t *testing.T) {
		s3Mock := &mockS3Repository{
			FakeUpload: func(context.Context, *domain.UploadS3param) error {
				return nil
			},
		}
		idGenMock := &mockUniqueIdGenerator{
			FakeGenerate: func() (string, error) {
				return imageName, nil
			},
		}

		u := &UseCase{
			repository:  s3Mock,
			idGenerator: idGenMock,
			cdnDomain:   cdnDomain,
		}

		r := &RequestBody{
			Image:          "",
			ImageExtension: ".png",
		}
		ctx := context.Background()
		res, err := u.CreateLgtmImage(ctx, *r)
		if err != nil {
			t.Fatalf("unexpected err = %s", err)
		}

		prefix, _ := domain.BuildS3Prefix(time.Now().UTC())
		want := &domain.UploadedLgtmImage{
			Url: "https://" + u.cdnDomain + "/" + prefix + imageName + ".webp",
		}

		if reflect.DeepEqual(res, want) == false {
			t.Errorf("\nwant\n%s\ngot\n%s", want, res)
		}
	})

	t.Run("Failure unexpect image extension", func(t *testing.T) {
		s3Mock := &mockS3Repository{
			FakeUpload: func(context.Context, *domain.UploadS3param) error {
				return nil
			},
		}
		idGenMock := &mockUniqueIdGenerator{
			FakeGenerate: func() (string, error) {
				return imageName, nil
			},
		}
		u := &UseCase{
			repository:  s3Mock,
			idGenerator: idGenMock,
			cdnDomain:   cdnDomain,
		}

		r := &RequestBody{
			Image:          "",
			ImageExtension: ".webp",
		}

		ctx := context.Background()
		_, err := u.CreateLgtmImage(ctx, *r)
		if err == nil {
			t.Fatal("expected to return an error, but no error")
		}
		if !errors.Is(err, domain.ErrInvalidImageExtension) {
			t.Fatalf("\nwant\n%s\ngot\n%s", domain.ErrInvalidImageExtension, err)
		}
	})

	t.Run("Failure generate image name", func(t *testing.T) {
		s3Mock := &mockS3Repository{
			FakeUpload: func(context.Context, *domain.UploadS3param) error {
				return nil
			},
		}
		idGenMock := &mockUniqueIdGenerator{
			FakeGenerate: func() (string, error) {
				return "", errors.New("dummy error")
			},
		}
		u := &UseCase{
			repository:  s3Mock,
			idGenerator: idGenMock,
			cdnDomain:   cdnDomain,
		}

		r := &RequestBody{
			Image:          "",
			ImageExtension: ".png",
		}

		ctx := context.Background()
		_, err := u.CreateLgtmImage(ctx, *r)

		if err == nil {
			t.Fatal("expected to return an error, but no error")
		}
		var want *domain.ErrGenerateImageName
		if !errors.As(err, &want) {
			t.Errorf("\nwant\n%T\ngot\n%T", want, errors.Unwrap(err))
		}
	})

	t.Run("Failure upload image to s3", func(t *testing.T) {
		s3Mock := &mockS3Repository{
			FakeUpload: func(context.Context, *domain.UploadS3param) error {
				return &domain.S3Error{
					Op:  "Upload",
					Err: errors.New("s3 upload dummy error"),
				}
			},
		}
		idGenMock := &mockUniqueIdGenerator{
			FakeGenerate: func() (string, error) {
				return imageName, nil
			},
		}
		u := &UseCase{
			repository:  s3Mock,
			idGenerator: idGenMock,
			cdnDomain:   cdnDomain,
		}

		r := &RequestBody{
			Image:          "",
			ImageExtension: ".png",
		}

		ctx := context.Background()
		_, err := u.CreateLgtmImage(ctx, *r)
		if err == nil {
			t.Fatal("expected to return an error, but no error")
		}
		var want *domain.S3Error
		if !errors.As(err, &want) {
			t.Errorf("\nwant\n%T\ngot\n%T", want, errors.Unwrap(err))
		}
	})
}
