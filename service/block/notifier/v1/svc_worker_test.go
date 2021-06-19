package v1

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func (s *ServiceTestSuite) TestGetLatestBlockNumber() {
	ctx, svc, provider := s.ctx, s.r.Svc, s.provider

	newBlockEventCh := make(chan types.Block)
	newBlockEventHandler := func(block types.Block) {
		newBlockEventCh <- block
	}
	svc.newBlockHandler = newBlockEventHandler

	chainForkedEventCh := make(chan struct{})
	chainForkedEventHandler := func() {
		chainForkedEventCh <- struct{}{}
	}
	svc.chainForkHandler = chainForkedEventHandler

	checkNewBlockEventTriggered := func(expBlock *types.Block) (bool, error) {
		select {
		case rcvBlock := <-newBlockEventCh:
			if expBlock == nil {
				return true, nil
			}
			if !bytes.Equal(expBlock.Hash().Bytes(), rcvBlock.Hash().Bytes()) {
				return true, fmt.Errorf("NewBlockEvent received: expected / received block hashes mismatch")
			}
		case <-time.After(50 * time.Millisecond):
			return false, fmt.Errorf("NewBlockEvent not received: timeout")
		}

		return true, nil
	}

	checkChainForkedEventTriggered := func() bool {
		select {
		case <-chainForkedEventCh:
			return true
		case <-time.After(50 * time.Millisecond):
			return false
		}
	}

	s.Run("Ok: NewBlockEvent triggered (1st time)", func() {
		provider.SetLatestBlockIdxOverride(0)
		expNewBlock, err := provider.GetBlockLatest(ctx)
		s.Require().NoError(err)

		svc.workerStep()
		nbEventReceived, nbEventErr := checkNewBlockEventTriggered(expNewBlock)
		cfEventReceived := checkChainForkedEventTriggered()
		s.Require().True(nbEventReceived)
		s.Require().False(cfEventReceived)
		s.Assert().NoError(nbEventErr)

	})

	s.Run("Ok: NewBlockEvent not triggered", func() {
		expNewBlock, err := provider.GetBlockLatest(ctx)
		s.Require().NoError(err)

		svc.workerStep()
		nbEventReceived, _ := checkNewBlockEventTriggered(expNewBlock)
		cfEventReceived := checkChainForkedEventTriggered()
		s.Require().False(nbEventReceived)
		s.Require().False(cfEventReceived)
	})

	s.Run("Ok: NewBlockEvent triggered (2nd time)", func() {
		provider.SetLatestBlockIdxOverride(1)
		expNewBlock, err := provider.GetBlockLatest(ctx)
		s.Require().NoError(err)

		svc.workerStep()
		nbEventReceived, nbEventErr := checkNewBlockEventTriggered(expNewBlock)
		cfEventReceived := checkChainForkedEventTriggered()
		s.Require().True(nbEventReceived)
		s.Require().False(cfEventReceived)
		s.Assert().NoError(nbEventErr)
	})

	s.Run("Ok: ChainForked event triggered (lower blockNumber)", func() {
		provider.SetLatestBlockIdxOverride(0)

		svc.workerStep()
		nbEventReceived, _ := checkNewBlockEventTriggered(nil)
		cfEventReceived := checkChainForkedEventTriggered()
		s.Require().False(nbEventReceived)
		s.Require().True(cfEventReceived)
	})
}
