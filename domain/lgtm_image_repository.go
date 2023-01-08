package domain

import (
	"context"
)

type LgtmImageRepository interface {
	FindAllIds(context.Context) (allIds []int32, err error)
	FindByIds(context.Context, []int32) (lgtmImageObjects []LgtmImageObject, err error)
	FindRecentlyCreated(context.Context, int) (lgtmImageObjects []LgtmImageObject, err error)
}
