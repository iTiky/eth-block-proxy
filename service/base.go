package service

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

const (
	// ConfigPrefix defines Viper key prefix for all service configurations.
	ConfigPrefix = "Service."
)

// BaseSvc defines common service fields and methods.
type BaseSvc struct {
	logger zerolog.Logger
}

// Logger returns an enriched with context data logger.
func (svc *BaseSvc) Logger(ctx context.Context) *zerolog.Logger {
	logger := svc.logger

	// Unique request ID
	if requestId := middleware.GetReqID(ctx); requestId != "" {
		logger = logger.With().Str("request_id", requestId).Logger()
	}

	return &logger
}

// Close implements pkg.Closer interface.
func (svc *BaseSvc) Close() error {
	return nil
}

// SetServiceName creates a new logger with service unique name context.
func (svc *BaseSvc) SetServiceName(name string) {
	svc.logger = svc.logger.With().Str("service", name).Logger()
}

// NewBaseSvc creates a new BaseSvc instance.
func NewBaseSvc(logger zerolog.Logger) BaseSvc {
	return BaseSvc{
		logger: logger,
	}
}
