package reader

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/itiky/eth-block-proxy/pkg"
)

// BlockReader defines the Ethereum chain block reader interface.
type BlockReader interface {
	pkg.Closer
	// GetBlock returns a block by its number (0 - the latest).
	GetBlock(ctx context.Context, blockIdx uint64) (*types.Block, error)
	// GetLatestBlockNumber returns the latest block number.
	GetLatestBlockNumber(ctx context.Context) (uint64, error)
}
