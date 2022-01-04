package domain

type LgtmImageRepository interface {
	FindAllIds() ([]int32, error)
	FindByIds(ids []int32) ([]LgtmImageObject, error)
}
