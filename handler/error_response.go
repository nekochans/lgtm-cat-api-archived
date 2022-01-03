package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func RenderErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	resBody := &ErrorResponse{Message: message}
	resBodyJson, _ := json.Marshal(resBody)

	fmt.Fprint(w, string(resBodyJson))
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
}
