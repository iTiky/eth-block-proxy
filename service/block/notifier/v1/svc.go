package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/itiky/eth-block-proxy/service"
	blockReader "github.com/itiky/eth-block-proxy/service/block/reader"
)

// BlockNotifierSvc polls the latest block info and emits NewBlock and ChainForked events (if registered).
type BlockNotifierSvc struct {
	service.BaseSvc
	reader blockReader.BlockReader // block reader service
	//
	pollDur     time.Duration // new block polling timeout
	latestBlock *types.Block  // the latest observed block data
	stopCh      chan struct{} // worker stop channel
	//
	newBlockHandler  NewBlockHandler    // new block event handler
	chainForkHandler ChainForkedHandler // chain forked event handler
}

// Close implements pkg.Closer interface.
func (svc *BlockNotifierSvc) Close() error {
	if svc.stopCh != nil {
		close(svc.stopCh)
	}

	return nil
}

// NewBlockNotifierSvc creates a new BlockNotifierSvc instance.
func NewBlockNotifierSvc(
	baseSvc service.BaseSvc,
	blockReader blockReader.BlockReader,
	newBlockHandler NewBlockHandler, chainForkHandler ChainForkedHandler,
) (*BlockNotifierSvc, error) {

	if blockReader == nil {
		return nil, fmt.Errorf("blockReader: nil")
	}

	svc := &BlockNotifierSvc{
		BaseSvc:          baseSvc,
		reader:           blockReader,
		pollDur:          1 * time.Second,
		stopCh:           make(chan struct{}),
		newBlockHandler:  newBlockHandler,
		chainForkHandler: chainForkHandler,
	}
	svc.SetServiceName("BlockNotifierSvc")

	if newBlockHandler != nil || chainForkHandler != nil {
		go svc.worker()
	} else {
		svc.Logger(context.TODO()).Info().Msg("Event handlers are not defined (service skipped)")
	}

	return svc, nil
}
