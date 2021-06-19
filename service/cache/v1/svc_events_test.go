package v1

import "sort"

func (s *ServiceTestSuite) TestNewBlockEvent() {
	ctx, svc, provider := s.ctx, s.svc, s.provider

	expLatestBlock, err := provider.GetBlockLatest(ctx)
	s.Require().NoError(err)

	// Check the latest index is not set
	s.Require().Empty(svc.latestBlockIdx)

	// Trigger event
	svc.HandleNewBlockEvent(*expLatestBlock)

	// Check the latest index was updated
	s.Assert().NotEmpty(svc.latestBlockIdx)
	s.Assert().EqualValues(expLatestBlock.NumberU64(), svc.latestBlockIdx)
}

func (s *ServiceTestSuite) TestChainForkedEvent() {
	ctx, svc, provider := s.ctx, s.svc, s.provider
	svc.config.ForkLength = 2

	getCachedBlockNumbers := func() (retNumbers []uint64) {
		for _, key := range svc.cache.Keys() {
			retNumbers = append(retNumbers, key.(uint64))
		}
		sort.Slice(retNumbers, func(i, j int) bool {
			return retNumbers[i] < retNumbers[j]
		})

		return
	}

	// Cache all fixtures and update latest index
	for _, block := range provider.Fixtures.Blocks {
		resp, err := svc.GetBlock(ctx, block.NumberU64())
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().EqualValues(block.Hash().Bytes(), resp.Hash().Bytes())
	}
	resp, err := svc.GetBlock(ctx, 0)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	cachedNumbersBefore := getCachedBlockNumbers()
	s.Require().Len(cachedNumbersBefore, len(provider.Fixtures.Blocks))
	s.Require().GreaterOrEqual(len(cachedNumbersBefore), 4) // 4 block are going to be cleaned up

	// Shift the latest one to the left
	// On the fork event two block on the left side from the latest one and two on the right side should be cleaned up
	// Total of 4 should be cleaned up in our case
	svc.latestBlockIdx--

	// Trigger event
	svc.HandleChainForkedEvent()
	cachedNumbersAfter := getCachedBlockNumbers()

	// Check cleared
	s.Assert().ElementsMatch(cachedNumbersBefore[:len(cachedNumbersBefore)-4], cachedNumbersAfter)
}
