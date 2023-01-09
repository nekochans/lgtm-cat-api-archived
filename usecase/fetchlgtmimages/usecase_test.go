package fetchlgtmimages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	sqlc "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
	"github.com/nekochans/lgtm-cat-api/test"
)

type mockLgtmImageRepository struct {
	domain.LgtmImageRepository
	FakeFindAllIds          func(context.Context) ([]int32, error)
	FakeFindByIds           func(context.Context, []int32) ([]domain.LgtmImageObject, error)
	FakeFindRecentlyCreated func(context.Context, int) ([]domain.LgtmImageObject, error)
}

func (m *mockLgtmImageRepository) FindAllIds(c context.Context) ([]int32, error) {
	return m.FakeFindAllIds(c)
}

func (m *mockLgtmImageRepository) FindByIds(c context.Context, ids []int32) ([]domain.LgtmImageObject, error) {
	return m.FakeFindByIds(c, ids)
}

func (m *mockLgtmImageRepository) FindRecentlyCreated(c context.Context, count int) ([]domain.LgtmImageObject, error) {
	return m.FakeFindRecentlyCreated(c, count)
}

var cdnDomain = "lgtm-images.lgtmeow.com"

var testDb *sql.DB

var mockErr = errors.New("mock error")

func createFindAllIdsResponse(len int) []int32 {
	var findAllIdsResponse []int32
	for i := 1; i <= len; i++ {
		findAllIdsResponse = append(
			findAllIdsResponse,
			int32(i),
		)
	}

	return findAllIdsResponse
}

func createLgtmImageObjects(len int) []domain.LgtmImageObject {
	var findByIdsMockResponse []domain.LgtmImageObject
	for i := 1; i < len; i++ {
		findByIdsMockResponse = append(
			findByIdsMockResponse,
			domain.LgtmImageObject{Id: int32(i), Path: "2022/02/22/22", Filename: "image-name" + fmt.Sprint(i)})
	}

	return findByIdsMockResponse
}

func createLgtmImages(imageObjects []domain.LgtmImageObject) []domain.LgtmImage {
	var images []domain.LgtmImage
	for _, v := range imageObjects {
		images = append(images, domain.LgtmImage{
			Id:  fmt.Sprint(v.Id),
			Url: "https://" + cdnDomain + "/" + v.Path + "/" + v.Filename + ".webp",
		})
	}

	return images
}

func TestMain(m *testing.M) {
	dbCreator := &test.DbCreator{}
	var err error
	testDb, err = dbCreator.Create()
	if err != nil {
		log.Panic(err)
	}

	seeder := &test.Seeder{Db: testDb}
	err = seeder.TruncateAllTable()
	if err != nil {
		log.Panic(err)
	}

	m.Run()

	_ = seeder.TruncateAllTable()
}

//nolint:funlen
func TestExtractRandomImages(t *testing.T) {
	var findByIdsMockResponse = createLgtmImageObjects(20)
	var want = createLgtmImages(findByIdsMockResponse)

	cases := []struct {
		name      string
		mock      *mockLgtmImageRepository
		want      []domain.LgtmImage
		expectErr error
	}{
		{
			name: "Success create LGTM image",
			mock: &mockLgtmImageRepository{
				FakeFindAllIds: func(context.Context) ([]int32, error) {
					return createFindAllIdsResponse(domain.FetchLgtmImageCount), nil
				},
				FakeFindByIds: func(context.Context, []int32) ([]domain.LgtmImageObject, error) {
					return findByIdsMockResponse, nil
				},
			},
			want:      want,
			expectErr: nil,
		},
		{
			name: "Failure error record count",
			mock: &mockLgtmImageRepository{
				FakeFindAllIds: func(context.Context) ([]int32, error) {
					return createFindAllIdsResponse(domain.FetchLgtmImageCount - 1), nil
				},
			},
			want:      nil,
			expectErr: domain.ErrRecordCount,
		},
		{
			name: "Failure find all ids",
			mock: &mockLgtmImageRepository{
				FakeFindAllIds: func(context.Context) ([]int32, error) {
					return nil, mockErr
				},
			},
			want:      nil,
			expectErr: mockErr,
		},
		{
			name: "Failure find by ids",
			mock: &mockLgtmImageRepository{
				FakeFindAllIds: func(context.Context) ([]int32, error) {
					return createFindAllIdsResponse(domain.FetchLgtmImageCount), nil
				},
				FakeFindByIds: func(context.Context, []int32) ([]domain.LgtmImageObject, error) {
					return nil, mockErr
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

			mock := tt.mock
			u := NewUseCase(mock, cdnDomain)

			ctx := context.Background()
			got, err := u.ExtractRandomImages(ctx)
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
					t.Errorf("ExtractRandomImages() value is mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestExtractRandomImagesConnectToDb(t *testing.T) {
	t.Run("Success extract random images", func(t *testing.T) {
		testDataDir, err := filepath.Abs("./testdata")
		if err != nil {
			t.Fatal("failed to read test data", err)
		}
		seeder := &test.Seeder{Db: testDb, DirPath: testDataDir}
		err = seeder.Execute()
		if err != nil {
			t.Fatal("failed seeder.Execute()", err)
		}

		q := sqlc.New(testDb)
		lgtmImageRepository := infrastructure.NewLgtmImageRepository(q)
		u := NewUseCase(lgtmImageRepository, cdnDomain)

		ctx := context.Background()
		res, err := u.ExtractRandomImages(ctx)
		if err != nil {
			t.Fatalf("unexpected err = %s", err)
		}

		if len(res) != domain.FetchLgtmImageCount {
			t.Fatalf("\nwant count\n%d\ngot  count\n%d", domain.FetchLgtmImageCount, len(res))
		}

		// ランダムに抽出するので型のみテストする
		for _, v := range res {
			_, ok := interface{}(v).(domain.LgtmImage)
			if !ok {
				t.Fatalf("\nwant\n%T\ngot\n%T", v, domain.LgtmImage{})
				return
			}
		}
	})
}

func TestRetrieveRecentlyCreatedImages(t *testing.T) {
	var findRecentlyCreatedResponse = createLgtmImageObjects(20)
	var want = createLgtmImages(findRecentlyCreatedResponse)

	cases := []struct {
		name      string
		mock      *mockLgtmImageRepository
		want      []domain.LgtmImage
		expectErr error
	}{
		{
			name: "Success create LGTM image",
			mock: &mockLgtmImageRepository{
				FakeFindRecentlyCreated: func(context.Context, int) ([]domain.LgtmImageObject, error) {
					return findRecentlyCreatedResponse, nil
				},
			},
			want:      want,
			expectErr: nil,
		},
		{
			name: "Failure find recently created images",
			mock: &mockLgtmImageRepository{
				FakeFindRecentlyCreated: func(context.Context, int) ([]domain.LgtmImageObject, error) {
					return nil, mockErr
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

			mock := tt.mock
			u := NewUseCase(mock, cdnDomain)

			ctx := context.Background()
			got, err := u.RetrieveRecentlyCreatedImages(ctx)
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
					t.Errorf("RetrieveRecentlyCreatedImages() value is mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestRetrieveRecentlyCreatedImagesConnectToDb(t *testing.T) {
	t.Run("Success retrieve recently created images", func(t *testing.T) {
		testDataDir, err := filepath.Abs("./testdata")
		if err != nil {
			t.Fatal("failed to read test data", err)
		}
		seeder := &test.Seeder{Db: testDb, DirPath: testDataDir}
		err = seeder.Execute()
		if err != nil {
			t.Fatal("failed seeder.Execute()", err)
		}

		q := sqlc.New(testDb)
		lgtmImageRepository := infrastructure.NewLgtmImageRepository(q)
		u := NewUseCase(lgtmImageRepository, cdnDomain)

		ctx := context.Background()
		got, err := u.RetrieveRecentlyCreatedImages(ctx)
		if err != nil {
			t.Fatalf("unexpected err = %s", err)
		}

		var want []domain.LgtmImage
		for i := 15; i > 6; i-- {
			var dd string
			if i < 10 {
				dd = "0" + fmt.Sprint(i)
			} else {
				dd = fmt.Sprint(i)
			}

			want = append(want, domain.LgtmImage{
				Id:  fmt.Sprint(i),
				Url: "https://" + u.cdnDomain + "/" + "2022/02/02/" + dd + "/" + "filename" + fmt.Sprint(i) + ".webp",
			})
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("RetrieveRecentlyCreatedImages() value is mismatch (-want +got):\n%s", diff)
		}
	})
}
