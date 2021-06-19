package v1

import (
	"context"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"
)

// HandleNewBlockEvent handles the NewBlock event updating the latest block number.
// This callback is called by the notifier service.
func (svc *CacheSvc) HandleNewBlockEvent(block types.Block) {
	atomic.StoreUint64(&svc.latestBlockIdx, block.NumberU64())
	svc.Logger(context.TODO()).Info().
		Uint64("block", block.NumberU64()).
		Str("event", "NewBlock").
		Msg("Latest block number updated")
}

// HandleChainForkedEvent handles the ChainForked event invalidating blocks around the latest one.
// This callback is called by the notifier service.
func (svc *CacheSvc) HandleChainForkedEvent() {
	latestBlockIdx := atomic.SwapUint64(&svc.latestBlockIdx, 0)
	if latestBlockIdx == 0 {
		svc.Logger(context.TODO()).Warn().
			Str("event", "ChainForked").
			Msg("Latest block number not set (skip)")
		return
	}

	minBlockIdx, maxBlockIdx := uint64(1), latestBlockIdx+uint64(svc.config.ForkLength)
	if latestBlockIdx > uint64(svc.config.ForkLength) {
		minBlockIdx = latestBlockIdx - uint64(svc.config.ForkLength)
	}

	for idx := minBlockIdx; idx <= maxBlockIdx; idx++ {
		svc.cache.Remove(idx)
		svc.Logger(context.TODO()).Warn().
			Uint64("block", idx).
			Str("event", "ChainForked").
			Msg("Block cache invalidated")
	}
}
