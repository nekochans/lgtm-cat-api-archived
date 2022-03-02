package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/usecase"
	"github.com/pkg/errors"
)

type ExtractRandomImagesHandler struct {
	extractRandomImagesUseCase *usecase.ExtractRandomImagesUseCase
}

func NewExtractRandomImagesHandler(c *usecase.ExtractRandomImagesUseCase) *ExtractRandomImagesHandler {
	return &ExtractRandomImagesHandler{
		extractRandomImagesUseCase: c,
	}
}

type ExtractRandomImagesResponse struct {
	LgtmImages []domain.LgtmImage `json:"lgtmImages"`
}

func (h *ExtractRandomImagesHandler) Extract(w http.ResponseWriter, r *http.Request) {

	lgtmImages, err := h.extractRandomImagesUseCase.ExtractRandomImages(r.Context())
	if err != nil {
		switch errors.Cause(err) {
		default:
			RenderErrorResponse(w, 500, err.Error())
		}

		return
	}

	response := &ExtractRandomImagesResponse{LgtmImages: lgtmImages}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")

	return
}
