package createltgmimage

import (
	"context"
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

func TestNewUseCase(t *testing.T) {
	t.Run("Success create LGTM image", func(t *testing.T) {
		s3Mock := &mockS3Repository{
			FakeUpload: func(context.Context, *domain.UploadS3param) error {
				return nil
			},
		}
		idGenMock := &mockUniqueIdGenerator{
			FakeGenerate: func() (string, error) {
				return "testimagename", nil
			},
		}
		imageName := "testimagename"

		u := &UseCase{
			repository:  s3Mock,
			idGenerator: idGenMock,
			cdnDomain:   imageName,
		}

		r := &RequestBody{
			Image:          "",
			ImageExtension: ".png",
		}
		ctx := context.Background()
		res, err := u.CreateLgtmImage(ctx, *r)
		if err != nil {
			t.Fatalf("エラーにならないはずなのにエラーになった %v", err)
		}

		prefix, _ := domain.BuildS3Prefix(time.Now().UTC())
		expected := &domain.UploadedLgtmImage{
			Url: "https://" + u.cdnDomain + "/" + prefix + imageName + ".webp",
		}

		if reflect.DeepEqual(res, expected) == false {
			t.Error("\nwant: ", res, "\ngot: ", expected)
		}
	})
}
