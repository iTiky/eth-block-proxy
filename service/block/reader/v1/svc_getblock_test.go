package v1

import (
	"errors"
	"time"
)

func (s *ServiceTestSuite) TestGetLatestBlockNumber() {
	ctx, svc, provider := s.ctx, s.svc, s.provider
	svc.config.RequestTimeoutDur = 10 * time.Millisecond
	svc.config.MinRetryDur = 1 * time.Millisecond
	svc.config.MinRetryDur = 2 * time.Millisecond

	latestBlockIdx, err := provider.GetLatestBlockNumber(ctx)
	s.Require().NoError(err)

	s.Run("Ok", func() {
		resp, err := svc.GetLatestBlockNumber(ctx)
		s.Require().NoError(err)
		s.Assert().EqualValues(latestBlockIdx, resp)
	})

	s.Run("Fail: provider error after fallback retriees", func() {
		provider.SetNextFails(5, 5*time.Millisecond, errors.New("mock"))
		_, err := svc.GetLatestBlockNumber(ctx)
		s.Require().Error(err)
	})

	s.Run("Ok: after few retries provider returns OK response", func() {
		provider.SetNextFails(2, 0, errors.New("mock"))
		resp, err := svc.GetLatestBlockNumber(ctx)
		s.Require().NoError(err)
		s.Assert().EqualValues(latestBlockIdx, resp)
	})
}

func (s *ServiceTestSuite) TestGetBlock() {
	ctx, svc, provider := s.ctx, s.svc, s.provider

	latestBlockIdx, err := provider.GetLatestBlockNumber(ctx)
	s.Require().NoError(err)
	latestBlock, err := provider.GetBlockLatest(ctx)
	s.Require().NoError(err)
	firstBlock, err := provider.GetBlockByNumber(ctx, provider.Fixtures.Blocks[0].NumberU64())
	s.Require().NoError(err)

	s.Run("Ok: latest", func() {
		resp, err := svc.GetBlock(ctx, 0)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Assert().EqualValues(latestBlock.Hash().Bytes(), resp.Hash().Bytes())
	})

	s.Run("Ok: by number", func() {
		resp, err := svc.GetBlock(ctx, firstBlock.NumberU64())
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Assert().EqualValues(firstBlock.Hash().Bytes(), resp.Hash().Bytes())
	})

	s.Run("Ok: not found", func() {
		resp, err := svc.GetBlock(ctx, latestBlockIdx+1)
		s.Require().NoError(err)
		s.Assert().Nil(resp)
	})
}
