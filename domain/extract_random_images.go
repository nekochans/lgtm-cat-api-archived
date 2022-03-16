package domain

import (
	"errors"
	"strconv"
)

var (
	ErrRecordCount = errors.New("the total record count is less than fetchLgtmImageCount")
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
