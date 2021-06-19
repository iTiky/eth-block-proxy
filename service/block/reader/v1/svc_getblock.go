package v1

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/core/types"
)

// GetBlock implements the reader.BlockReader interface.
func (svc *FallbackBlockReaderSvc) GetBlock(ctx context.Context, blockIdx uint64) (retBlock *types.Block, retErr error) {
	ctx, ctxCancel := svc.enrichCtxWithTimeout(ctx)
	defer ctxCancel()

	// Define retry operation
	var retryOp backoff.Operation
	if blockIdx == 0 {
		retryOp = func() error {
			var err error
			retBlock, err = svc.provider.GetBlockLatest(ctx)
			return err
		}
	} else {
		retryOp = func() error {
			var err error
			retBlock, err = svc.provider.GetBlockByNumber(ctx, blockIdx)
			return err
		}
	}

	// Request with retry
	if err := backoff.Retry(retryOp, svc.getRetryCfg(ctx)); err != nil {
		retErr = err
		return
	}

	return
}

// GetLatestBlockNumber implements the reader.BlockReader interface.
func (svc *FallbackBlockReaderSvc) GetLatestBlockNumber(ctx context.Context) (retBlock uint64, retErr error) {
	ctx, ctxCancel := svc.enrichCtxWithTimeout(ctx)
	defer ctxCancel()

	retryOp := func() error {
		var err error
		retBlock, err = svc.provider.GetLatestBlockNumber(ctx)
		return err
	}

	if err := backoff.Retry(retryOp, svc.getRetryCfg(ctx)); err != nil {
		retErr = err
		return
	}

	return
}

// enrichCtxWithTimeout adds request timeout if it wasn't set by the caller.
func (svc *FallbackBlockReaderSvc) enrichCtxWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	// Check if timeout is already set or disabled: leave ctx "as is"
	if _, found := ctx.Deadline(); found || svc.config.IsRetryDisabled() {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, svc.config.RequestTimeoutDur)
}

// getRetryCfg returns retry policy params (no retry / exponential).
func (svc *FallbackBlockReaderSvc) getRetryCfg(ctx context.Context) backoff.BackOff {
	// No-backoff policy
	if svc.config.IsRetryDisabled() {
		return &backoff.StopBackOff{}
	}

	// Estimate time left for the request
	deadlineValue, deadlineFound := ctx.Deadline()
	if !deadlineFound {
		return &backoff.StopBackOff{}
	}
	retryTimeout := time.Until(deadlineValue)

	// No-backoff policy: no time left for retries
	if retryTimeout <= 0 {
		svc.Logger(ctx).Warn().
			Dur("deadline_dur", retryTimeout).
			Msg("Overdue context deadline (retry skip)")
		return &backoff.StopBackOff{}
	}

	// Exponential backoff policy
	expBackOff := &backoff.ExponentialBackOff{
		InitialInterval:     svc.config.MinRetryDur,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         svc.config.MaxRetryDur,
		MaxElapsedTime:      retryTimeout,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	expBackOff.Reset()

	return expBackOff
}
