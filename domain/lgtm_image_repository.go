package domain

import "context"

type LgtmImageRepository interface {
	FindAllIds(context.Context) ([]int32, error)
	FindByIds(context.Context, []int32) ([]LgtmImageObject, error)
}
