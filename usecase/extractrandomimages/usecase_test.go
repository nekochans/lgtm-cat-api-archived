package extractrandomimages

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/nekochans/lgtm-cat-api/domain"
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

//nolint:funlen
func TestExtractRandomImages(t *testing.T) {
	t.Run("Success extract random images", func(t *testing.T) {
		var findAllIdsResponse []int32
		for i := 1; i <= domain.FetchLgtmImageCount; i++ {
			findAllIdsResponse = append(
				findAllIdsResponse,
				int32(i),
			)
		}

		var findByIdsMockResponse []domain.LgtmImageObject
		for i := 1; i < 20; i++ {
			findByIdsMockResponse = append(
				findByIdsMockResponse,
				domain.LgtmImageObject{Id: int32(i), Path: "2022/02/22/22", Filename: "image-name" + fmt.Sprint(i)})
		}

		mock := &mockLgtmImageRepository{
			FakeFindAllIds: func(context.Context) ([]int32, error) {
				return findAllIdsResponse, nil
			},
			FakeFindByIds: func(context.Context, []int32) ([]domain.LgtmImageObject, error) {
				return findByIdsMockResponse, nil
			},
		}
		u := NewUseCase(mock, cdnDomain)

		ctx := context.Background()
		res, err := u.ExtractRandomImages(ctx)
		if err != nil {
			t.Fatalf("unexpected err = %s", err)
		}

		var want []domain.LgtmImage
		for _, v := range findByIdsMockResponse {
			want = append(want, domain.LgtmImage{
				Id:  fmt.Sprint(v.Id),
				Url: "https://" + u.cdnDomain + "/" + v.Path + "/" + v.Filename + ".webp",
			})
		}

		if reflect.DeepEqual(res, want) == false {
			t.Errorf("\nwant\n%s\ngot\n%s", want, res)
		}
	})

	t.Run("Failure error record count", func(t *testing.T) {
		var findAllIdsResponse []int32
		for i := 1; i <= domain.FetchLgtmImageCount-1; i++ {
			findAllIdsResponse = append(
				findAllIdsResponse,
				int32(i),
			)
		}

		mock := &mockLgtmImageRepository{
			FakeFindAllIds: func(context.Context) ([]int32, error) {
				return findAllIdsResponse, nil
			},
			FakeFindByIds: func(context.Context, []int32) ([]domain.LgtmImageObject, error) {
				return nil, nil
			},
		}
		u := NewUseCase(mock, cdnDomain)

		ctx := context.Background()
		_, err := u.ExtractRandomImages(ctx)
		if err == nil {
			t.Fatal("expected to return an error, but no error")
		}
		if !errors.Is(err, domain.ErrRecordCount) {
			t.Fatalf("\nwant\n%s\ngot\n%s", domain.ErrRecordCount, err)
		}
	})

	t.Run("Failure find all ids", func(t *testing.T) {
		mock := &mockLgtmImageRepository{
			FakeFindAllIds: func(context.Context) ([]int32, error) {
				return nil, &domain.LgtmImageError{
					Op:  "FindAllIds",
					Err: errors.New("FindAllIds dummy error"),
				}
			},
			FakeFindByIds: func(context.Context, []int32) ([]domain.LgtmImageObject, error) {
				return nil, nil
			},
		}

		u := NewUseCase(mock, cdnDomain)

		ctx := context.Background()
		_, err := u.ExtractRandomImages(ctx)
		if err == nil {
			t.Fatal("expected to return an error, but no error")
		}
		var want *domain.LgtmImageError
		if !errors.As(err, &want) {
			t.Errorf("\nwant\n%T\ngot\n%T", want, errors.Unwrap(err))
		}
	})

	t.Run("Failure find all ids", func(t *testing.T) {
		var findAllIdsResponse []int32
		for i := 1; i <= domain.FetchLgtmImageCount; i++ {
			findAllIdsResponse = append(
				findAllIdsResponse,
				int32(i),
			)
		}

		mock := &mockLgtmImageRepository{
			FakeFindAllIds: func(context.Context) ([]int32, error) {
				return findAllIdsResponse, nil
			},
			FakeFindByIds: func(context.Context, []int32) ([]domain.LgtmImageObject, error) {
				return nil, &domain.LgtmImageError{
					Op:  "FindByIds",
					Err: errors.New("FindByIds dummy error"),
				}
			},
		}

		u := NewUseCase(mock, cdnDomain)

		ctx := context.Background()
		_, err := u.ExtractRandomImages(ctx)
		if err == nil {
			t.Fatal("expected to return an error, but no error")
		}
		var want *domain.LgtmImageError
		if !errors.As(err, &want) {
			t.Errorf("\nwant\n%T\ngot\n%T", want, errors.Unwrap(err))
		}
	})
}

func TestRetrieveRecentlyCreatedImages(t *testing.T) {
	t.Run("Success retrieve recently created images", func(t *testing.T) {
		var findRecentlyCreatedResponse []domain.LgtmImageObject
		for i := 1; i < 20; i++ {
			findRecentlyCreatedResponse = append(
				findRecentlyCreatedResponse,
				domain.LgtmImageObject{Id: int32(i), Path: "2022/02/22/22", Filename: "image-name" + fmt.Sprint(i)})
		}

		mock := &mockLgtmImageRepository{
			FakeFindRecentlyCreated: func(context.Context, int) ([]domain.LgtmImageObject, error) {
				return findRecentlyCreatedResponse, nil
			},
		}
		u := NewUseCase(mock, cdnDomain)

		ctx := context.Background()
		res, err := u.RetrieveRecentlyCreatedImages(ctx)
		if err != nil {
			t.Fatalf("unexpected err = %s", err)
		}

		var want []domain.LgtmImage
		for _, v := range findRecentlyCreatedResponse {
			want = append(want, domain.LgtmImage{
				Id:  fmt.Sprint(v.Id),
				Url: "https://" + u.cdnDomain + "/" + v.Path + "/" + v.Filename + ".webp",
			})
		}

		if reflect.DeepEqual(res, want) == false {
			t.Errorf("\nwant\n%s\ngot\n%s", want, res)
		}
	})

	t.Run("Failure find recently created images", func(t *testing.T) {
		mock := &mockLgtmImageRepository{
			FakeFindRecentlyCreated: func(context.Context, int) ([]domain.LgtmImageObject, error) {
				return nil, &domain.LgtmImageError{
					Op:  "FindRecentlyCreated",
					Err: errors.New("FindRecentlyCreated dummy error"),
				}
			},
		}
		u := NewUseCase(mock, cdnDomain)

		ctx := context.Background()
		_, err := u.RetrieveRecentlyCreatedImages(ctx)
		if err == nil {
			t.Fatal("expected to return an error, but no error")
		}
		var want *domain.LgtmImageError
		if !errors.As(err, &want) {
			t.Errorf("\nwant\n%T\ngot\n%T", want, errors.Unwrap(err))
		}
	})
}
