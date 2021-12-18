package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ExtractRandomImagesResponseBody TODO パッケージ構成を見直した際に名前を変更
type ExtractRandomImagesResponseBody struct {
	LgtmImages []LgtmImage `json:"lgtmImages"`
}

type LgtmImage struct {
	Id       string `json:"id"`
	ImageUrl string `json:"url"`
}

func ExtractRandomImages(w http.ResponseWriter, r *http.Request) {
	length := r.URL.Query().Get("count")
	fmt.Println(length)

	lgtmImage := []LgtmImage{
		{
			Id:       "1",
			ImageUrl: "https://stg.example.com/2021/12/18/12/577c893a-f830-4c14-bf00-33255231ad31.webp",
		},
		{
			Id:       "2",
			ImageUrl: "https://stg.example.com/2021/12/18/12/577c893a-f830-4c14-bf00-33255231ad31.webp",
		},
	}

	response := &ExtractRandomImagesResponseBody{LgtmImages: lgtmImage}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")

	return
}
