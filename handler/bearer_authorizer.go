package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/nekochans/lgtm-cat-api/derrors"
	"github.com/nekochans/lgtm-cat-api/domain"
)

type BearerAuthorizer struct {
	JwtValidator domain.JwtValidator
}

func NewBearerAuthorizer(validator domain.JwtValidator) *BearerAuthorizer {
	return &BearerAuthorizer{
		JwtValidator: validator,
	}
}

func (a *BearerAuthorizer) Authorize(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := extractLogger(r.Context())

		key := r.Header.Get("Authorization")
		accessToken, err := a.extractAccessToken(key)
		if err != nil {
			logger.Error(err)
			RenderErrorResponse(w, Unauthorized)

			return
		}
		err = a.JwtValidator.Validate(accessToken)
		if err != nil {
			logger.Error(err)
			RenderErrorResponse(w, Unauthorized)

			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (a *BearerAuthorizer) extractAccessToken(auth string) (accessToken string, err error) {
	defer derrors.Wrap(&err, "BearerAuthorizer.extractAccessToken()")

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
