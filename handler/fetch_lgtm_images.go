package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nekochans/lgtm-cat-api/domain"
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

func (h *fetchImagesHandler) RetrieveRecentlyCreated(w http.ResponseWriter, r *http.Request) {
	lgtmImages, err := h.useCase.RetrieveRecentlyCreatedImages(r.Context())
	if err != nil {
		log.Println(err)
		RenderErrorResponse(w, InternalServerError)
		return
	}

	response := &RetrieveRecentlyCreatedImagesResponse{LgtmImages: lgtmImages}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}
