package createltgmimage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

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
	prefix, _ := domain.BuildS3Prefix(time.Now().UTC())
	mockErr := errors.New("mock error")

	cases := []struct {
		name           string
		imageExtension string
		s3Mock         *mockS3Repository
		idGenMock      *mockUniqueIdGenerator
		want           *domain.UploadedLgtmImage
		expectErr      error
	}{
		{
			name:           "Success create LGTM image",
			imageExtension: ".png",
			s3Mock: &mockS3Repository{
				FakeUpload: func(context.Context, *domain.UploadS3param) error {
					return nil
				},
			},
			idGenMock: &mockUniqueIdGenerator{
				FakeGenerate: func() (string, error) {
					return imageName, nil
				},
			},
			want: &domain.UploadedLgtmImage{
				Url: "https://" + cdnDomain + "/" + prefix + imageName + ".webp",
			},
			expectErr: nil,
		},
		{
			name:           "Failure unexpect image extension",
			imageExtension: ".webp",
			s3Mock: &mockS3Repository{
				FakeUpload: func(context.Context, *domain.UploadS3param) error {
					return nil
				},
			},
			idGenMock: &mockUniqueIdGenerator{
				FakeGenerate: func() (string, error) {
					return imageName, nil
				},
			},
			want:      nil,
			expectErr: domain.ErrInvalidImageExtension,
		},
		{
			name:           "Failure generate image name",
			imageExtension: ".png",
			s3Mock: &mockS3Repository{
				FakeUpload: func(context.Context, *domain.UploadS3param) error {
					return nil
				},
			},
			idGenMock: &mockUniqueIdGenerator{
				FakeGenerate: func() (string, error) {
					return "", mockErr
				},
			},
			want:      nil,
			expectErr: mockErr,
		},
		{
			name:           "Failure upload image to s3",
			imageExtension: ".png",
			s3Mock: &mockS3Repository{
				FakeUpload: func(context.Context, *domain.UploadS3param) error {
					return mockErr
				},
			},
			idGenMock: &mockUniqueIdGenerator{
				FakeGenerate: func() (string, error) {
					return imageName, nil
				},
			},
			want:      nil,
			expectErr: mockErr,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s3Mock := tt.s3Mock
			idGenMock := tt.idGenMock

			u := &UseCase{
				repository:  s3Mock,
				idGenerator: idGenMock,
				cdnDomain:   cdnDomain,
			}

			r := &RequestBody{
				Image:          "",
				ImageExtension: tt.imageExtension,
			}
			ctx := context.Background()
			got, err := u.CreateLgtmImage(ctx, *r)
			if tt.expectErr != nil {
				if err == nil {
					t.Fatal("expected to return an error, but no error")
				}
				if !errors.Is(err, tt.expectErr) {
					t.Errorf("\nwant\n%#v\ngot\n%#v", tt.expectErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected err = %s", err)
				}
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("CreateLgtmImage() value is mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
