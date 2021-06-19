package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/itiky/eth-block-proxy/provider/block/mock"
	"github.com/itiky/eth-block-proxy/service/cache/testutil"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx context.Context
	r   *testutil.BlockCacheSvcTestResource
	//
	defaultSvcCfg Config
	svc           *CacheSvc
	provider      *mock.MockBlockProvider
}

func (s *ServiceTestSuite) SetupSuite() {
	r, err := NewBlockCacheSvcTestResource()
	if err != nil {
		panic(fmt.Errorf("NewBlockCacheSvcTestResource: %w", err))
	}

	s.ctx = context.TODO()
	s.r = r

	s.svc = r.Svc.(*CacheSvc)
	s.provider = r.ReaderSvcRes.ProviderRes.Provider.(*mock.MockBlockProvider)
	s.defaultSvcCfg = s.svc.config
}

func (s *ServiceTestSuite) SetupTest() {
	s.provider.SetLatestBlockIdxOverride(-1)
	s.svc.cache.Purge()
	s.svc.latestBlockIdx = 0
	s.svc.config = s.defaultSvcCfg
}

func TestBockCacheV1Service(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
