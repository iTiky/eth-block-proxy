package v1

import (
	"fmt"

	"github.com/itiky/eth-block-proxy/service"
	blockReaderV1 "github.com/itiky/eth-block-proxy/service/block/reader/v1"
)

// NewBlockNotifierSvcTestResource creates an instance of testutil.BlockNotifierSvcTestResource using mock dependencies.
func NewBlockNotifierSvcTestResource() (*BlockNotifierSvcTestResource, error) {
	blockReaderRes, err := blockReaderV1.NewBlockReaderSvcTestResource()
	if err != nil {
		return nil, fmt.Errorf("NewBlockReaderSvcTestResource: %w", err)
	}

	svc, err := NewBlockNotifierSvc(service.NewTestBaseSvc(), blockReaderRes.Svc, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("NewBlockNotifierSvc: %w", err)
	}

	return &BlockNotifierSvcTestResource{
		Svc:          svc,
		ReaderSvcRes: blockReaderRes,
	}, nil
}
