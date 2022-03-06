package domain

type UniqueIdGenerator interface {
	Generate() (string, error)
}
