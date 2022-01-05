package usecase

import (
	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
	"log"
	"math/rand"
	"time"
)

type ExtractRandomImagesUseCase struct {
	Repository *infrastructure.LgtmImageRepository
	CdnDomain  string
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
		n := rand.Intn(recordCount - 1)
		if contains(randomIds, ids[n]) {
			continue
		}
		i++
		randomIds = append(randomIds, ids[n])
	}

	return randomIds
}

func (u *ExtractRandomImagesUseCase) ExtractRandomImages() ([]domain.LgtmImage, error) {
	ids, err := u.Repository.FindAllIds()
	if err != nil {
		return nil, domain.ErrCountRecords
	}
	if len(ids) < domain.FetchLgtmImageCount {
		log.Println("The total record count is less than fetchLgtmImageCount")
		return nil, domain.ErrFetchImages
	}

	var randomIds = pickupRandomIdsNoDuplicates(ids, domain.FetchLgtmImageCount)

	rows, err := u.Repository.FindByIds(randomIds)
	if err != nil {
		return nil, domain.ErrFetchImages
	}

	var lgtmImages []domain.LgtmImage
	for _, row := range rows {
		lgtmImage := domain.CreateLgtmImage(row.Id, u.CdnDomain, row.Path, row.Filename)
		lgtmImages = append(lgtmImages, *lgtmImage)
	}

	return lgtmImages, nil
}
