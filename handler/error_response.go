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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(apiErr.status)
	fmt.Fprint(w, string(resBodyJson))
}
