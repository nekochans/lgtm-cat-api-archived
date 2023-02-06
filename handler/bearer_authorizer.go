package handler

import (
	"errors"
	"net/http"
	"strings"
)

type BearerAuthorizer struct{}

func NewBearerAuthorizer() *BearerAuthorizer {
	return &BearerAuthorizer{}
}

func (a *BearerAuthorizer) Authorize(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := extractLogger(r.Context())

		key := r.Header.Get("Authorization")
		_, err := a.extractAccessToken(key)
		if err != nil {
			logger.Error(err)
			RenderErrorResponse(w, Unauthorized)

			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (a *BearerAuthorizer) extractAccessToken(auth string) (string, error) {
	bearerLen := 7
	if len(auth) < bearerLen {
		return "", errors.New("invalid bearer authorization header")
	}

	authType := strings.ToLower(auth[:6])
	if authType != "bearer" {
		return "", errors.New("invalid bearer authorization header")
	}

	return auth[7:], nil
}
