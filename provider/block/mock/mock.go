package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/itiky/eth-block-proxy/provider/block"
)

var _ block.BlockProvider = (*MockBlockProvider)(nil)

// MockBlockProvider is a mock BlockProvider with predefined fixtures and response altering methods.
type MockBlockProvider struct {
	Fixtures Fixtures
	//
	latestIdxOverride int // latest block index override
	//
	respFails int           // number of consecutive response fails with respErr and respDelay
	respErr   error         // next response error
	respDelay time.Duration // new response delay
}

// Close implements block.BlockProvider interface.
func (p *MockBlockProvider) Close() error {
	return nil
}

// GetBlockLatest implements block.BlockProvider interface.
func (p *MockBlockProvider) GetBlockLatest(ctx context.Context) (*types.Block, error) {
	if err := p.checkFailResponse(); err != nil {
		return nil, err
	}

	if p.latestIdxOverride >= 0 {
		blockCopy := p.Fixtures.Blocks[p.latestIdxOverride]
		return &blockCopy, nil
	}
	blockCopy := p.Fixtures.Blocks[p.Fixtures.LatestBlockIdx]

	return &blockCopy, nil
}

// GetBlockByNumber implements block.BlockProvider interface.
func (p *MockBlockProvider) GetBlockByNumber(ctx context.Context, blockIdx uint64) (*types.Block, error) {
	if err := p.checkFailResponse(); err != nil {
		return nil, err
	}

	for _, block := range p.Fixtures.Blocks {
		if block.NumberU64() == blockIdx {
			blockCopy := block
			return &blockCopy, nil
		}
	}

	return nil, nil
}

// GetLatestBlockNumber implements block.BlockProvider interface.
func (p *MockBlockProvider) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	if err := p.checkFailResponse(); err != nil {
		return 0, err
	}

	if p.latestIdxOverride >= 0 {
		return p.Fixtures.Blocks[p.latestIdxOverride].NumberU64(), nil
	}

	return p.Fixtures.Blocks[p.Fixtures.LatestBlockIdx].NumberU64(), nil
}

// SetLatestBlockIdxOverride sets fixtures block slice index for the latest block (-1 - no override).
func (p *MockBlockProvider) SetLatestBlockIdxOverride(idx int) {
	p.latestIdxOverride = idx
}

// SetNextFails sets next number of consecutive response fails with delay and error.
func (p *MockBlockProvider) SetNextFails(count int, delay time.Duration, err error) {
	p.respFails, p.respDelay, p.respErr = count, delay, err
}

// checkFailResponse alters provider response.
func (p *MockBlockProvider) checkFailResponse() error {
	if p.respFails == 0 {
		return nil
	}

	p.respFails--
	time.Sleep(p.respDelay)

	return p.respErr
}

// NewMockBlockProvider creates a new MockBlockProvider instance with fixtures.
func NewMockBlockProvider() (*MockBlockProvider, error) {
	fixtures, err := NewFixtures()
	if err != nil {
		return nil, fmt.Errorf("NewFixtures: %w", err)
	}

	return &MockBlockProvider{
		Fixtures:          fixtures,
		latestIdxOverride: -1,
	}, nil
}
