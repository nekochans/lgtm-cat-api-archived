package domain

import (
	"strconv"

	"github.com/pkg/errors"
)

var (
	ErrCountRecords = errors.New("failed to count LGTM images records")
	ErrFetchImages  = errors.New("failed to fetch LGTM images")
)

const FetchLgtmImageCount = 9

type LgtmImage struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func CreateLgtmImage(id int32, domain, path, filename string) *LgtmImage {
	return &LgtmImage{
		Id:  strconv.Itoa(int(id)),
		Url: "https://" + domain + "/" + path + "/" + filename + ".webp",
	}
}
