package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/usecase/createltgmimage"
)

type createLgtmImageHandler struct {
	useCase *createltgmimage.UseCase
}

func NewCreateLgtmImageHandler(c *createltgmimage.UseCase) *createLgtmImageHandler {
	return &createLgtmImageHandler{
		useCase: c,
	}
}

type CreateLgtmImageResponse struct {
	ImageUrl string `json:"imageUrl"`
}

func (h *createLgtmImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		RenderErrorResponse(w, InternalServerError)
		return
	}

	var reqBody createltgmimage.RequestBody
	if err := json.Unmarshal(req, &reqBody); err != nil {
		log.Println(err)
		RenderErrorResponse(w, BadRequest)
	}

	image, err := h.useCase.CreateLgtmImage(r.Context(), reqBody)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidImageExtension):
			log.Println(err)
			RenderErrorResponse(w, UnprocessableEntity)
		default:
			log.Println(err)
			RenderErrorResponse(w, InternalServerError)
		}

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response := &CreateLgtmImageResponse{ImageUrl: image.Url}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
}
