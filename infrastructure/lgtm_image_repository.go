package infrastructure

import (
	"context"

	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/domain"
)

type LgtmImageRepository struct {
	Db *db.Queries
}

func (r *LgtmImageRepository) FindAllIds() ([]int32, error) {
	ctx := context.Background()

	ids, err := r.Db.ListLgtmImageIds(ctx)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *LgtmImageRepository) FindByIds(ids []int32) ([]domain.LgtmImage, error) {
	ctx := context.Background()

	var listLgtmImagesParams = db.ListLgtmImagesParams{
		ID:   ids[0],
		ID_2: ids[1],
		ID_3: ids[2],
		ID_4: ids[3],
		ID_5: ids[4],
		ID_6: ids[5],
		ID_7: ids[6],
		ID_8: ids[7],
		ID_9: ids[8],
	}

	rows, err := r.Db.ListLgtmImages(ctx, listLgtmImagesParams)
	if err != nil {
		return nil, err
	}

	var lgtmImage []domain.LgtmImage
	for _, v := range rows {
		lgtmImage = append(lgtmImage, domain.LgtmImage{ID: v.ID, Path: v.Path, Filename: v.Filename})
	}

	return lgtmImage, nil
}
