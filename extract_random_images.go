package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
)

const fetchLgtmImageCount = 9

// ExtractRandomImagesResponseBody TODO パッケージ構成を見直した際に名前を変更
type ExtractRandomImagesResponseBody struct {
	LgtmImages []LgtmImage `json:"lgtmImages"`
}

type LgtmImage struct {
	Id       string `json:"id"`
	ImageUrl string `json:"url"`
}

func contains(elems []int32, v int32) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func pickupRandomIdsNoDuplicates(ids []int32, listCount int) []int32 {

	rand.Seed(time.Now().Unix())
	recordCount := len(ids)

	var randomIds []int32
	for i := 1; i <= listCount; {
		n := rand.Int31n(int32(recordCount - 1))
		if contains(randomIds, ids[n]) {
			continue
		}
		i++
		randomIds = append(randomIds, ids[n])
	}

	return randomIds
}

func ExtractRandomImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	ids, err := q.ListLgtmImageIds(ctx)
	if err != nil {
		RenderErrorResponse(w, 500, "Failed count LGTM images records")
	}

	if len(ids) < fetchLgtmImageCount {
		log.Println("The total record count is less than fetchLgtmImageCount")
		RenderErrorResponse(w, 500, "Failed fetch LGTM images")
		return
	}

	var randomIds = pickupRandomIdsNoDuplicates(ids, fetchLgtmImageCount)
	var listLgtmImagesParams = db.ListLgtmImagesParams{
		ID:   randomIds[0],
		ID_2: randomIds[1],
		ID_3: randomIds[2],
		ID_4: randomIds[3],
		ID_5: randomIds[4],
		ID_6: randomIds[5],
		ID_7: randomIds[6],
		ID_8: randomIds[7],
		ID_9: randomIds[8],
	}

	rows, err := q.ListLgtmImages(ctx, listLgtmImagesParams)
	if err != nil {
		RenderErrorResponse(w, 500, "Failed fetch LGTM images")
		return
	}

	var lgtmImages []LgtmImage
	for _, row := range rows {
		lgtmImage := LgtmImage{
			Id:       strconv.Itoa(int(row.ID)),
			ImageUrl: "https://" + lgtmImagesCdnDomain + "/" + row.Path + "/" + row.Filename + ".webp",
		}
		lgtmImages = append(lgtmImages, lgtmImage)
	}

	response := &ExtractRandomImagesResponseBody{LgtmImages: lgtmImages}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")

	return
}
