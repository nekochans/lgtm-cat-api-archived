package domain

import (
	"context"
	"fmt"
)

type LgtmImageRepository interface {
	FindAllIds(context.Context) (allIds []int32, err error)
	FindByIds(context.Context, []int32) (lgtmImageObjects []LgtmImageObject, err error)
	FindRecentlyCreated(context.Context, int) (lgtmImageObjects []LgtmImageObject, err error)
}

type LgtmImageError struct {
	Op  string
	Err error
}

func (e *LgtmImageError) Error() string {
	return fmt.Sprintf("lgtmImageRepository: %s, %s", e.Op, e.Err)
}
