package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nekochans/lgtm-cat-api/domain"
	"github.com/nekochans/lgtm-cat-api/usecase/createltgmimage"
	"github.com/pkg/errors"
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
		RenderErrorResponse(w, 500, "Failed Read Request Body")
		return
	}

	var reqBody createltgmimage.RequestBody
	if err := json.Unmarshal(req, &reqBody); err != nil {
		RenderErrorResponse(w, 400, err.Error())
	}

	image, err := h.useCase.CreateLgtmImage(r.Context(), reqBody)
	if err != nil {
		switch errors.Cause(err) {
		case domain.ErrInvalidImageExtension:
			fmt.Println(err)

			RenderErrorResponse(w, 422, err.Error())
		default:
			RenderErrorResponse(w, 500, err.Error())
		}

		return
	}

	response := &CreateLgtmImageResponse{ImageUrl: image.Url}
	responseJson, _ := json.Marshal(response)
	fmt.Fprint(w, string(responseJson))
	w.WriteHeader(202)
	w.Header().Add("Content-Type", "application/json")

	return
}
