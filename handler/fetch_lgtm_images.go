package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
	"github.com/nekochans/lgtm-cat-api/usecase/fetchlgtmimages"
)

type fetchImagesHandler struct {
	useCase *fetchlgtmimages.UseCase
}

func NewFetchImagesHandler(c *fetchlgtmimages.UseCase) *fetchImagesHandler {
	return &fetchImagesHandler{
		useCase: c,
	}
}

type ExtractRandomImagesResponse struct {
	LgtmImages []domain.LgtmImage `json:"lgtmImages"`
}

type RetrieveRecentlyCreatedImagesResponse struct {
	LgtmImages []domain.LgtmImage `json:"lgtmImages"`
}

func (h *fetchImagesHandler) Extract(w http.ResponseWriter, r *http.Request) {
	logger := extractLogger(r.Context())

	lgtmImages, err := h.useCase.ExtractRandomImages(r.Context())
	if err != nil {
		logger.Error(err)
		infrastructure.ReportError(r.Context(), err)
		RenderErrorResponse(w, InternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &ExtractRandomImagesResponse{LgtmImages: lgtmImages}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
}

func (h *fetchImagesHandler) RetrieveRecentlyCreated(w http.ResponseWriter, r *http.Request) {
	logger := extractLogger(r.Context())
	lgtmImages, err := h.useCase.RetrieveRecentlyCreatedImages(r.Context())
	if err != nil {
		logger.Error(err)
		infrastructure.ReportError(r.Context(), err)
		RenderErrorResponse(w, InternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := &RetrieveRecentlyCreatedImagesResponse{LgtmImages: lgtmImages}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
}
