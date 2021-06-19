package v1

import "github.com/ethereum/go-ethereum/core/types"

// NewBlockHandler defines handler for NewBlock event.
type NewBlockHandler func(block types.Block)

// ChainForkedHandler defines handler for ChainForked event.
type ChainForkedHandler func()
