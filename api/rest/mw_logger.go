package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// LoggerMiddleware router middleware logs request / response meta.
func LoggerMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logEvent := logger.Info()
			wrappedWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			reqStartedAt := time.Now()
			defer func() {
				httpScheme := "http"
				if r.TLS != nil {
					httpScheme = "https"
				}

				logEvent = logEvent.
					Str("request_id", middleware.GetReqID(r.Context())).
					Str("from", r.RemoteAddr).
					Str("to", fmt.Sprintf("%s://%s%s %s", httpScheme, r.Host, r.RequestURI, r.Proto)).
					Int("status", wrappedWriter.Status()).
					Int("bytes", wrappedWriter.BytesWritten()).
					Dur("elapsed", time.Since(reqStartedAt))
				if wrappedWriter.Status() == http.StatusOK {
					logEvent.Msg("Request handled")
				} else {
					logEvent.Msg("Request failed")
				}
			}()

			next.ServeHTTP(wrappedWriter, WithLogEvent(r, logEvent))
		}

		return http.HandlerFunc(fn)
	}
}

// GetLogEvent returns the in-context zerolog.Event for a request.
func GetLogEvent(r *http.Request) *zerolog.Event {
	entry, _ := r.Context().Value(middleware.LogEntryCtxKey).(*zerolog.Event)
	return entry
}

// WithLogEvent sets the in-context zerolog.Event for a request.
func WithLogEvent(r *http.Request, logEvent *zerolog.Event) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), middleware.LogEntryCtxKey, logEvent))
}
