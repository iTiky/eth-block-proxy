package v1_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	he "github.com/gavv/httpexpect/v2"
)

func (s *ApiTestSuite) TestBlockGet() {
	ctx, provider := s.ctx, s.provider

	server := httptest.NewServer(s.router)
	defer server.Close()
	e := he.New(s.T(), server.URL)

	s.Run("Fail: invalid {blockNumber}", func() {
		e.GET("/v1/block/abc").
			Expect().
			Status(http.StatusNotFound)
	})

	s.Run("Ok: block not found", func() {
		e.GET("/v1/block/100500").
			Expect().
			Status(http.StatusNotFound)
	})

	s.Run("Ok: latest", func() {
		blockIdx, expResp := s.getBlockResp(provider.GetBlockLatest(ctx))

		e.GET(fmt.Sprintf("/v1/block/%d", blockIdx)).
			Expect().
			Status(http.StatusOK).
			Body().Contains(expResp)
	})

	s.Run("Ok: by {blockNumber}", func() {
		blockIdx, expResp := s.getBlockResp(provider.GetBlockByNumber(ctx, provider.Fixtures.Blocks[0].NumberU64()))

		e.GET(fmt.Sprintf("/v1/block/%d", blockIdx)).
			Expect().
			Status(http.StatusOK).
			Body().Contains(expResp)
	})
}

func (s *ApiTestSuite) TestBlockGetTx() {
	ctx, provider := s.ctx, s.provider

	server := httptest.NewServer(s.router)
	defer server.Close()
	e := he.New(s.T(), server.URL)

	s.Run("Fail: invalid {txHash}", func() {
		e.GET("/v1/block/latest/txs/abc").
			Expect().
			Status(http.StatusBadRequest)
	})

	s.Run("Ok: not found", func() {
		e.GET("/v1/block/latest/txs/0x00010203").
			Expect().
			Status(http.StatusNotFound)
	})

	s.Run("Ok: found", func() {
		blockIdx, txHash, expResp := s.getBlockTxResp(provider.GetBlockLatest(ctx))

		e.GET(fmt.Sprintf("/v1/block/%d/txs/%s", blockIdx, txHash)).
			Expect().
			Status(http.StatusOK).
			Body().Contains(expResp)
	})

	s.Run("Ok: found (with 0x prefix)", func() {
		blockIdx, txHash, expResp := s.getBlockTxResp(provider.GetBlockLatest(ctx))

		e.GET(fmt.Sprintf("/v1/block/%d/txs/0x%s", blockIdx, txHash)).
			Expect().
			Status(http.StatusOK).
			Body().Contains(expResp)
	})
}
