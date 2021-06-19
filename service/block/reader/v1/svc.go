package v1

import (
	"context"
	"fmt"

	blockProvider "github.com/itiky/eth-block-proxy/provider/block"
	"github.com/itiky/eth-block-proxy/service"
	"github.com/itiky/eth-block-proxy/service/block/reader"
)

var _ reader.BlockReader = (*FallbackBlockReaderSvc)(nil)

// FallbackBlockReaderSvc fires the BlockProvider requests with exponential backoff.
type FallbackBlockReaderSvc struct {
	service.BaseSvc
	config   Config                      // service config
	provider blockProvider.BlockProvider // raw block getter
}

// NewFallbackBlockReaderSvc creates a new FallbackBlockReaderSvc instance.
func NewFallbackBlockReaderSvc(
	baseSvc service.BaseSvc,
	blockProvider blockProvider.BlockProvider,
) (*FallbackBlockReaderSvc, error) {

	cfg := BuildConfig()

	if blockProvider == nil {
		return nil, fmt.Errorf("blockProvider: nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	svc := &FallbackBlockReaderSvc{
		BaseSvc:  baseSvc,
		config:   cfg,
		provider: blockProvider,
	}
	svc.SetServiceName("FallbackBlockReaderSvc")

	if cfg.IsRetryDisabled() {
		svc.Logger(context.TODO()).Info().Msg("Retry policy: disabled")
	}

	return svc, nil
}
