package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/itiky/eth-block-proxy/provider/block"
)

var _ block.BlockProvider = (*CloudflareBlockProvider)(nil)

const (
	// EthGatewayUrl defines CloudFlare Ethereum gateway URL.
	EthGatewayUrl = "https://cloudflare-eth.com/"
)

// CloudflareBlockProvider keeps the Ethereum client and used to get blocks.
type CloudflareBlockProvider struct {
	client *ethclient.Client
}

// GetBlockLatest implements the block.BlockProvider interface.
func (p *CloudflareBlockProvider) GetBlockLatest(ctx context.Context) (*types.Block, error) {
	resp, err := p.client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("block request (latest): %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("block request (latest): nil response")
	}

	if err := p.validateBlockResult(resp); err != nil {
		return nil, fmt.Errorf("block request (latest): invalid response received: %w", err)
	}

	return resp, nil
}

// GetBlockByNumber implements the block.BlockProvider interface.
func (p *CloudflareBlockProvider) GetBlockByNumber(ctx context.Context, blockIdx uint64) (*types.Block, error) {
	reqArg := (&big.Int{}).SetUint64(blockIdx)

	resp, err := p.client.BlockByNumber(ctx, reqArg)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("block request (%d): %w", blockIdx, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("block request (%d): nil response", blockIdx)
	}

	if err := p.validateBlockResult(resp); err != nil {
		return nil, fmt.Errorf("block request (latest): invalid response received: %w", err)
	}

	return resp, nil
}

// GetLatestBlockNumber implements the block.BlockProvider interface.
func (p *CloudflareBlockProvider) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	resp, err := p.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("latest block number request: %w", err)
	}

	return resp, nil
}

// Close implements pkg.Closer interface.
func (p *CloudflareBlockProvider) Close() error {
	if p.client == nil {
		return nil
	}
	p.client.Close()

	return nil
}

// validateBlockResult validates types.Block result to prevent returning an invalid data.
// Rules observed from CloudFlare Ethereum Gateway API use-cases (should not happen, but happens).
func (p *CloudflareBlockProvider) validateBlockResult(block *types.Block) error {
	if block.NumberU64() == 0 {
		return fmt.Errorf("0 block number")
	}
	if len(block.Hash().Bytes()) == 0 {
		return fmt.Errorf("empty hash")
	}

	return nil
}

// NewCloudflareBlockProvider creates a new CloudflareBlockProvider instance.
func NewCloudflareBlockProvider() (*CloudflareBlockProvider, error) {
	client, err := ethclient.Dial(EthGatewayUrl)
	if err != nil {
		return nil, fmt.Errorf("creating Ethereum client: %w", err)
	}

	return &CloudflareBlockProvider{
		client: client,
	}, nil
}
