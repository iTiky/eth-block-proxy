package cache

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/itiky/eth-block-proxy/pkg"
)

// BlockCache defines the Ethereum block cached reader interface.
type BlockCache interface {
	pkg.Closer
	// GetBlock returns types.Block by blockNumber from cache on hit.
	GetBlock(ctx context.Context, blockNumber uint64) (*types.Block, error)
}
