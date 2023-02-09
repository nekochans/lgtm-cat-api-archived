package fetchlgtmimages

import (
	"context"
	"math/rand"
	"time"

	"github.com/nekochans/lgtm-cat-api/derrors"
	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
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
		n := rand.Intn(recordCount) //nolint:gosec
		if contains(randomIds, ids[n]) {
			continue
		}
		i++
		randomIds = append(randomIds, ids[n])
	}

	return randomIds
}

func (u *UseCase) ExtractRandomImages(ctx context.Context) (randomImages []domain.LgtmImage, err error) {
	defer derrors.Wrap(&err, "UseCase.ExtractRandomImages()")

	ids, err := u.repository.FindAllIds(ctx)
	if err != nil {
		return nil, err
	}
	if len(ids) < domain.FetchLgtmImageCount {
		infrastructure.ReportError(ctx, domain.ErrRecordCount)
		return nil, domain.ErrRecordCount
	}

	var randomIds = pickupRandomIdsNoDuplicates(ids, domain.FetchLgtmImageCount)

	rows, err := u.repository.FindByIds(ctx, randomIds)
	if err != nil {
		return nil, err
	}

	var lgtmImages []domain.LgtmImage
	for _, row := range rows {
		lgtmImage := domain.CreateLgtmImage(row.Id, u.cdnDomain, row.Path, row.Filename)
		lgtmImages = append(lgtmImages, *lgtmImage)
	}

	return lgtmImages, nil
}

func (u *UseCase) RetrieveRecentlyCreatedImages(ctx context.Context) (recentlyImages []domain.LgtmImage, err error) {
	defer derrors.Wrap(&err, "UseCase.RetrieveRecentlyCreatedImages()")

	rows, err := u.repository.FindRecentlyCreated(ctx, domain.FetchLgtmImageCount)
	if err != nil {
		return nil, err
	}

	var lgtmImages []domain.LgtmImage
	for _, row := range rows {
		lgtmImage := domain.CreateLgtmImage(row.Id, u.cdnDomain, row.Path, row.Filename)
		lgtmImages = append(lgtmImages, *lgtmImage)
	}

	return lgtmImages, nil
}
