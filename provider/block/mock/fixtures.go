package mock

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Fixtures keeps mock block data.
type Fixtures struct {
	Blocks         []types.Block // Sorted by blockNumber Blocks
	LatestBlockIdx int           // Blocks index for the latest block
}

// NewFixtures creates a new Fixtures instance.
func NewFixtures() (Fixtures, error) {
	f := Fixtures{}

	latestBlockNumber := uint64(0)
	for blockIdx, blockHexStr := range AllBlocks {
		blockEncoded, err := hex.DecodeString(blockHexStr)
		if err != nil {
			return Fixtures{}, fmt.Errorf("decoding fixture block [%d]: hex.DecodeString", blockIdx)
		}

		var block types.Block
		if err := rlp.DecodeBytes(blockEncoded, &block); err != nil {
			return Fixtures{}, fmt.Errorf("decoding fixture block [%d]: rlp.DecodeBytes", blockIdx)
		}

		f.Blocks = append(f.Blocks, block)
		if block.NumberU64() > latestBlockNumber {
			latestBlockNumber = block.NumberU64()
			f.LatestBlockIdx = blockIdx
		}
	}

	return f, nil
}
