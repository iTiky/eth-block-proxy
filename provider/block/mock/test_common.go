package mock

import (
	"fmt"

	"github.com/itiky/eth-block-proxy/provider/block/testutil"
)

// NewBlockProviderTestResource creates an instance of testutil.BlockProviderTestResource using mock provider.
func NewBlockProviderTestResource() (*testutil.BlockProviderTestResource, error) {
	provider, err := NewMockBlockProvider()
	if err != nil {
		return nil, fmt.Errorf("NewMockBlockProvider: %w", err)
	}

	return &testutil.BlockProviderTestResource{
		Provider: provider,
	}, nil
}
