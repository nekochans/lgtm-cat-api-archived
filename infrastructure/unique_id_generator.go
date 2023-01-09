package infrastructure

import (
	"github.com/google/uuid"
	"github.com/nekochans/lgtm-cat-api/derrors"
)

type UuidGenerator struct{}

func (g *UuidGenerator) Generate() (id string, err error) {
	defer derrors.Wrap(&err, "UuidGenerator.Generate()")

	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return uid.String(), nil
}
