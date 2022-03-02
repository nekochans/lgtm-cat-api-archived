package infrastructure

import (
	"context"
	"time"

	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/domain"
)

type LgtmImageRepository struct {
	db *db.Queries
}

func NewLgtmImageRepository(db *db.Queries) *LgtmImageRepository {
	return &LgtmImageRepository{db: db}
}

func (r *LgtmImageRepository) FindAllIds(c context.Context) ([]int32, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	ids, err := r.db.ListLgtmImageIds(ctx)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *LgtmImageRepository) FindByIds(c context.Context, ids []int32) ([]domain.LgtmImageObject, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

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

	rows, err := r.db.ListLgtmImages(ctx, listLgtmImagesParams)
	if err != nil {
		return nil, err
	}

	var lgtmImage []domain.LgtmImageObject
	for _, v := range rows {
		lgtmImage = append(lgtmImage, domain.LgtmImageObject{Id: v.ID, Path: v.Path, Filename: v.Filename})
	}

	return lgtmImage, nil
}
