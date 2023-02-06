package handler

import (
	"net/http"
)

var (
	BadRequest          = newBadRequest("Bad Request")
	UnprocessableEntity = newUnprocessableEntity("Unprocessable Entity")
	InternalServerError = newInternalServerError("Internal Server Error")
	Unauthorized        = newUnauthorized("Unauthorized")
)

type apiError struct {
	message string
	status  int
}

func newError(msg string, st int) *apiError {
	e := &apiError{
		message: msg,
		status:  st,
	}
	return e
}
func newBadRequest(msg string) *apiError {
	e := newError(msg, http.StatusBadRequest)
	return e
}

func newUnprocessableEntity(msg string) *apiError {
	e := newError(msg, http.StatusUnprocessableEntity)
	return e
}

func newInternalServerError(msg string) *apiError {
	e := newError(msg, http.StatusInternalServerError)
	return e
}

func newUnauthorized(msg string) *apiError {
	e := newError(msg, http.StatusUnauthorized)
	return e
}
