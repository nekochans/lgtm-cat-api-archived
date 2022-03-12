package extractrandomimages

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nekochans/lgtm-cat-api/domain"
)

type UseCase struct {
	repository domain.LgtmImageRepository
	cdnDomain  string
}

func NewUseCase(r domain.LgtmImageRepository, c string) *UseCase {
	return &UseCase{
		repository: r,
		cdnDomain:  c,
	}
}

func contains(elems []int32, v int32) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func pickupRandomIdsNoDuplicates(ids []int32, listCount int) []int32 {
	rand.Seed(time.Now().Unix())
	recordCount := len(ids)

	var randomIds []int32
	for i := 1; i <= listCount; {
		n := rand.Intn(recordCount - 1) //nolint:gosec
		if contains(randomIds, ids[n]) {
			continue
		}
		i++
		randomIds = append(randomIds, ids[n])
	}

	return randomIds
}

func (u *UseCase) ExtractRandomImages(ctx context.Context) ([]domain.LgtmImage, error) {
	ids, err := u.repository.FindAllIds(ctx)
	if err != nil {
		return nil, fmt.Errorf("faild to extract randam images: %w", err)
	}
	if len(ids) < domain.FetchLgtmImageCount {
		return nil, domain.ErrRecordCount
	}

	var randomIds = pickupRandomIdsNoDuplicates(ids, domain.FetchLgtmImageCount)

	rows, err := u.repository.FindByIds(ctx, randomIds)
	if err != nil {
		return nil, fmt.Errorf("faild to extract randam images: %w", err)
	}

	var lgtmImages []domain.LgtmImage
	for _, row := range rows {
		lgtmImage := domain.CreateLgtmImage(row.Id, u.cdnDomain, row.Path, row.Filename)
		lgtmImages = append(lgtmImages, *lgtmImage)
	}

	return lgtmImages, nil
}
