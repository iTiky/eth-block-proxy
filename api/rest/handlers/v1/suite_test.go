package v1_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-chi/chi/v5"
	"github.com/itiky/eth-block-proxy/api/rest"
	"github.com/itiky/eth-block-proxy/api/rest/handlers/v1/schema"
	"github.com/itiky/eth-block-proxy/provider/block/mock"
	"github.com/itiky/eth-block-proxy/service/cache/testutil"
	blockCacheSvcV1 "github.com/itiky/eth-block-proxy/service/cache/v1"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type ApiTestSuite struct {
	suite.Suite
	ctx    context.Context
	r      *testutil.BlockCacheSvcTestResource
	router *chi.Mux
	//
	provider *mock.MockBlockProvider
}

func (s *ApiTestSuite) SetupSuite() {
	r, err := blockCacheSvcV1.NewBlockCacheSvcTestResource()
	if err != nil {
		panic(fmt.Errorf("NewBlockCacheSvcTestResource: %w", err))
	}

	s.ctx = context.TODO()
	s.r = r
	s.router = rest.NewRouter(zerolog.Nop(), r.Svc)

	s.provider = r.ReaderSvcRes.ProviderRes.Provider.(*mock.MockBlockProvider)
}

func (s *ApiTestSuite) getBlockResp(block *types.Block, err error) (uint64, string) {
	s.Require().NoError(err)
	s.Require().NotNil(block)

	resp, err := schema.NewBlockResponse(block)
	s.Require().NoError(err)

	respBz, err := json.Marshal(resp)
	s.Require().NoError(err)

	return block.NumberU64(), string(respBz)
}

func (s *ApiTestSuite) getBlockTxResp(block *types.Block, err error) (uint64, string, string) {
	s.Require().NoError(err)
	s.Require().NotNil(block)
	s.Require().NotEmpty(block.Transactions())

	tx := block.Transactions()[0]
	s.Require().NotNil(tx)

	resp := schema.NewTxResponse(block, tx.Hash())
	respBz, err := json.Marshal(resp)
	s.Require().NoError(err)

	txHash := hex.EncodeToString(tx.Hash().Bytes())

	return block.NumberU64(), txHash, string(respBz)
}

func TestApiHandlersV1(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
