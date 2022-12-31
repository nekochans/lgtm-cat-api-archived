package handler

import (
	"context"
	"net/http"

	"github.com/nekochans/lgtm-cat-api/infrastructure"
)

type contextKey string

const logKey contextKey = "log"

func withLogger(logger infrastructure.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get("X-Request-Id")
			withLogger := logger.With(infrastructure.Field{
				Key:   "x_request_id",
				Value: requestId,
			})

			ctx := context.WithValue(r.Context(), logKey, withLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func extractLogger(ctx context.Context) infrastructure.Logger {
	return ctx.Value(logKey).(infrastructure.Logger)
}
