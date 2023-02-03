package infrastructure

import (
	"context"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"

	"github.com/nekochans/lgtm-cat-api/derrors"
)

const flushTimeSeconds = 2

func InitSentry() (err error) {
	defer derrors.Wrap(&err, "InitSentry")

	env := os.Getenv("ENV")
	err = sentry.Init(sentry.ClientOptions{
		Environment:      env,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
	})
	defer sentry.Flush(flushTimeSeconds * time.Second)

	if err != nil {
		return err
	}

	return nil
}

func NewSentryHttp() *sentryhttp.Handler {
	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
}

func ReportError(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
	}
}
