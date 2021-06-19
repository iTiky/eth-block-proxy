package testutil

import "github.com/itiky/eth-block-proxy/provider/block"

// BlockProviderTestResource defines test resource for provider containing all the necessary dependencies.
type BlockProviderTestResource struct {
	Provider block.BlockProvider
}
