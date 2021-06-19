package v1

import (
	"fmt"

	blockProviderMock "github.com/itiky/eth-block-proxy/provider/block/mock"
	"github.com/itiky/eth-block-proxy/service"
	"github.com/itiky/eth-block-proxy/service/block/reader/testutil"
)

// NewBlockReaderSvcTestResource creates an instance of testutil.BlockReaderSvcTestResource using mock dependencies.
func NewBlockReaderSvcTestResource() (*testutil.BlockReaderSvcTestResource, error) {
	blockProviderRes, err := blockProviderMock.NewBlockProviderTestResource()
	if err != nil {
		return nil, fmt.Errorf("NewBlockProviderTestResource: %w", err)
	}

	svc, err := NewFallbackBlockReaderSvc(service.NewTestBaseSvc(), blockProviderRes.Provider)
	if err != nil {
		return nil, fmt.Errorf("NewFallbackBlockReaderSvc: %w", err)
	}

	return &testutil.BlockReaderSvcTestResource{
		Svc:         svc,
		ProviderRes: blockProviderRes,
	}, nil
}
