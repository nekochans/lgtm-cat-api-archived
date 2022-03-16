package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/usecase/extractrandomimages"
)

type extractRandomImagesHandler struct {
	useCase *extractrandomimages.UseCase
}

func NewExtractRandomImagesHandler(c *extractrandomimages.UseCase) *extractRandomImagesHandler {
	return &extractRandomImagesHandler{
		useCase: c,
	}
}

type ExtractRandomImagesResponse struct {
	LgtmImages []domain.LgtmImage `json:"lgtmImages"`
}

func (h *extractRandomImagesHandler) Extract(w http.ResponseWriter, r *http.Request) {
	lgtmImages, err := h.useCase.ExtractRandomImages(r.Context())
	if err != nil {
		log.Println(err)
		RenderErrorResponse(w, InternalServerError)
		return
	}

	response := &ExtractRandomImagesResponse{LgtmImages: lgtmImages}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}
