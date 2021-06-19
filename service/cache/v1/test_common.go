package v1

import (
	"fmt"

	"github.com/itiky/eth-block-proxy/service"
	blockReaderV1 "github.com/itiky/eth-block-proxy/service/block/reader/v1"
	"github.com/itiky/eth-block-proxy/service/cache/testutil"
)

// NewBlockCacheSvcTestResource creates an instance of testutil.BlockCacheSvcTestResource using mock dependencies.
func NewBlockCacheSvcTestResource() (*testutil.BlockCacheSvcTestResource, error) {
	blockReaderRes, err := blockReaderV1.NewBlockReaderSvcTestResource()
	if err != nil {
		return nil, fmt.Errorf("NewBlockReaderSvcTestResource: %w", err)
	}

	svc, err := NewCacheSvc(service.NewTestBaseSvc(), blockReaderRes.Svc)
	if err != nil {
		return nil, fmt.Errorf("NewCacheSvc: %w", err)
	}

	return &testutil.BlockCacheSvcTestResource{
		Svc:          svc,
		ReaderSvcRes: blockReaderRes,
	}, nil
}
