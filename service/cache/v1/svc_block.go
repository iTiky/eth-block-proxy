package v1

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"
)

// GetBlock implements cache.BlockCache interface.
func (svc *CacheSvc) GetBlock(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	if blockNumber == 0 {
		idx, err := svc.getLatestBlockIdx(ctx)
		if err != nil {
			return nil, err
		}
		blockNumber = idx
	}

	return svc.getCachedBlockByNumber(ctx, blockNumber)
}

// getLatestBlockIdx returns the latest block number and updates service state if the value is not set.
func (svc *CacheSvc) getLatestBlockIdx(ctx context.Context) (uint64, error) {
	blockIdx := atomic.LoadUint64(&svc.latestBlockIdx)
	if blockIdx > 0 {
		return blockIdx, nil
	}

	latestBlockIdx, err := svc.reader.GetLatestBlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("reading latestBlockNumber: %w", err)
	}
	atomic.StoreUint64(&svc.latestBlockIdx, latestBlockIdx)

	return latestBlockIdx, nil
}

// getCachedBlockByNumber returns cached block data or updates the cache on miss.
func (svc *CacheSvc) getCachedBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	// Get from cache
	blockRaw, found := svc.cache.Get(blockNumber)
	if found {
		svc.Logger(ctx).Debug().
			Uint64("block", blockNumber).
			Msg("Cache used")
		block := blockRaw.(types.Block)
		return &block, nil
	}

	// Request
	block, err := svc.reader.GetBlock(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, nil
	}

	// Update cache
	blockEvicted := svc.cache.Add(blockNumber, *block)
	svc.Logger(ctx).Debug().
		Uint64("block", blockNumber).
		Msgf("Cache updated (evicted: %v)", blockEvicted)

	return block, nil
}
