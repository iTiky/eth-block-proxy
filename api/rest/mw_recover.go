package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// RecoveryMiddleware router middleware adds panic recovery step.
// Middleware throws 500 HTTP status code on panic events.
func RecoveryMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				logEvent := GetLogEvent(r)
				if logEvent != nil {
					logEvent.Interface("panic", rvr)
				}
				middleware.PrintPrettyStack(rvr)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
