package handler

import (
	"encoding/json"
	"net/http"
)

type healthCheckResponse struct {
	StatusCode int `json:"statusCode"`
}

type healthCheckHandler struct{}

func NewHealthCheckHandler() *healthCheckHandler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) Check(w http.ResponseWriter, r *http.Request) {
	logger := extractLogger(r.Context())

	responseJson, err := json.Marshal(&healthCheckResponse{
		StatusCode: http.StatusOK,
	})
	if err != nil {
		logger.Error(err)
		RenderErrorResponse(w, InternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(responseJson); err != nil {
		logger.Error(err)
		RenderErrorResponse(w, InternalServerError)
		return
	}
}
