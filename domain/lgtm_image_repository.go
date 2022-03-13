package domain

import (
	"context"
	"fmt"
)

type LgtmImageRepository interface {
	FindAllIds(context.Context) ([]int32, error)
	FindByIds(context.Context, []int32) ([]LgtmImageObject, error)
}

type LgtmImageError struct {
	Op  string
	Err error
}

func (e *LgtmImageError) Error() string {
	return fmt.Sprintf("lgtmImageRepository: %s, %s", e.Op, e.Err)
}
