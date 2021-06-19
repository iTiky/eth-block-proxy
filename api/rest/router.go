package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	v1 "github.com/itiky/eth-block-proxy/api/rest/handlers/v1"
	"github.com/itiky/eth-block-proxy/service/cache"
	"github.com/rs/zerolog"
)

// NewRouter creates a new HTTP router.
// nolint: errcheck
func NewRouter(logger zerolog.Logger, blockCacheSvc cache.BlockCache) *chi.Mux {
	logger = logger.With().Str("api", "v1").Logger()
	r := chi.NewRouter()

	// Common middlewares
	r.Use(middleware.RequestID)     // Unique request ID middleware (adds ID to the request context)
	r.Use(LoggerMiddleware(logger)) // Requests log middleware
	r.Use(RecoveryMiddleware)       // Panic recovery middleware (throws 500 HTTP status code)

	// API: common
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		w.Write([]byte("pong"))
	})
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	// API: v1
	handlersV1 := v1.NewApiHandlers(blockCacheSvc)
	r.Group(func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON)) // Content-type set to JSON middleware

		r.Route("/v1", func(r chi.Router) {
			r.Route("/block", func(r chi.Router) {
				r.Route(v1.UrlParamBlockPattern, func(r chi.Router) {
					r.Use(handlersV1.BlockCtx)      // Stores block artifact into request context
					r.Get("/", handlersV1.BlockGet) // GET /v1/block/{blockNumber} (with pattern check)

					r.Route("/txs", func(r chi.Router) {
						r.Route(v1.UrlParamTxHashPattern, func(r chi.Router) {
							r.Get("/", handlersV1.BlockGetTx) // GET /v1/block/{blockNumber}/txs/{txHash} (with pattern check)
						})
					})
				})
			})
		})
	})

	return r
}
