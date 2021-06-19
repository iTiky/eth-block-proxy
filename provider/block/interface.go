package block

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/itiky/eth-block-proxy/pkg"
)

// BlockProvider interface defines the Ethereum chain block getter calls.
type BlockProvider interface {
	pkg.Closer
	// GetBlockLatest returns the latest chain block.
	GetBlockLatest(ctx context.Context) (*types.Block, error)
	// GetBlockByNumber returns block by its number (index).
	GetBlockByNumber(ctx context.Context, blockIdx uint64) (*types.Block, error)
	// GetLatestBlockNumber returns the latest block number.
	GetLatestBlockNumber(ctx context.Context) (uint64, error)
}
