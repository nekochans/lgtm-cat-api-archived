package infrastructure

import (
	"context"
	"time"

	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/derrors"
	"github.com/nekochans/lgtm-cat-api/domain"
)

type lgtmImageRepository struct {
	db *db.Queries
}

const dbTimeoutSecond = 10

func NewLgtmImageRepository(db *db.Queries) *lgtmImageRepository {
	return &lgtmImageRepository{db: db}
}

func (r *lgtmImageRepository) FindAllIds(c context.Context) (allIds []int32, err error) {
	defer derrors.Wrap(&err, "lgtmImageRepository.FindAllIds()")

	ctx, cancel := context.WithTimeout(c, dbTimeoutSecond*time.Second)
	defer cancel()

	ids, err := r.db.ListLgtmImageIds(ctx)
	if err != nil {
		ReportError(ctx, err)
		return nil, err
	}

	return ids, nil
}

func (r *lgtmImageRepository) FindByIds(
	c context.Context,
	ids []int32,
) (lgtmImageObjects []domain.LgtmImageObject, err error) {
	defer derrors.Wrap(&err, "lgtmImageRepository.FindByIds(%v)", ids)

	ctx, cancel := context.WithTimeout(c, dbTimeoutSecond*time.Second)
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
		ReportError(ctx, err)
		return nil, err
	}

	var lgtmImage []domain.LgtmImageObject
	for _, v := range rows {
		lgtmImage = append(lgtmImage, domain.LgtmImageObject{Id: v.ID, Path: v.Path, Filename: v.Filename})
	}

	return lgtmImage, nil
}

func (r *lgtmImageRepository) FindRecentlyCreated(
	c context.Context,
	count int,
) (lgtmImageObjects []domain.LgtmImageObject, err error) {
	defer derrors.Wrap(&err, "lgtmImageRepository.FindRecentlyCreated(%v)", count)

	ctx, cancel := context.WithTimeout(c, dbTimeoutSecond*time.Second)
	defer cancel()

	rows, err := r.db.ListRecentlyCreatedLgtmImages(ctx, int32(count))
	if err != nil {
		ReportError(ctx, err)
		return nil, err
	}

	var lgtmImage []domain.LgtmImageObject
	for _, v := range rows {
		lgtmImage = append(lgtmImage, domain.LgtmImageObject{Id: v.ID, Path: v.Path, Filename: v.Filename})
	}

	return lgtmImage, nil
}
