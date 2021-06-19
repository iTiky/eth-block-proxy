package v1

func (s *ServiceTestSuite) TestGetBlock() {
	ctx, svc, provider := s.ctx, s.svc, s.provider

	s.Run("Ok: not found", func() {
		resp, err := svc.GetBlock(ctx, 100500)
		s.Require().NoError(err)
		s.Assert().Nil(resp)
	})

	s.Run("Ok: get block by number", func() {
		expBlock := provider.Fixtures.Blocks[0]

		_, cached := svc.cache.Peek(expBlock.NumberU64())
		s.Require().False(cached)

		resp, err := svc.GetBlock(ctx, expBlock.NumberU64())
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Assert().EqualValues(expBlock.Hash().Bytes(), resp.Hash().Bytes())

		_, cached = svc.cache.Peek(expBlock.NumberU64())
		s.Assert().True(cached)
	})

	s.Run("Ok: get latest block", func() {
		expBlock, err := provider.GetBlockLatest(ctx)
		s.Require().NoError(err)

		_, cached := svc.cache.Peek(expBlock.NumberU64())
		s.Require().False(cached)

		s.Require().Empty(svc.latestBlockIdx)

		resp, err := svc.GetBlock(ctx, 0)
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Assert().EqualValues(expBlock.Hash().Bytes(), resp.Hash().Bytes())

		_, cached = svc.cache.Peek(expBlock.NumberU64())
		s.Assert().True(cached)

		s.Assert().NotEmpty(svc.latestBlockIdx)
		s.Assert().EqualValues(expBlock.NumberU64(), svc.latestBlockIdx)
	})
}
