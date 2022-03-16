package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func RenderErrorResponse(w http.ResponseWriter, apiErr *apiError) {
	resBody := &ErrorResponse{Message: apiErr.message}
	resBodyJson, _ := json.Marshal(resBody)

	fmt.Fprint(w, string(resBodyJson))
	w.WriteHeader(apiErr.status)
	w.Header().Add("Content-Type", "application/json")
}
