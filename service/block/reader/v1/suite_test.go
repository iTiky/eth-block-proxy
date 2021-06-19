package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/itiky/eth-block-proxy/provider/block/mock"
	"github.com/itiky/eth-block-proxy/service/block/reader/testutil"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx context.Context
	r   *testutil.BlockReaderSvcTestResource
	//
	defaultSvcCfg Config
	svc           *FallbackBlockReaderSvc
	provider      *mock.MockBlockProvider
}

func (s *ServiceTestSuite) SetupSuite() {
	r, err := NewBlockReaderSvcTestResource()
	if err != nil {
		panic(fmt.Errorf("NewBlockReaderSvcTestResource: %w", err))
	}

	s.ctx = context.TODO()
	s.r = r

	s.svc = r.Svc.(*FallbackBlockReaderSvc)
	s.provider = r.ProviderRes.Provider.(*mock.MockBlockProvider)
	s.defaultSvcCfg = s.svc.config
}

func (s *ServiceTestSuite) SetupTest() {
	s.svc.config = s.defaultSvcCfg
	s.provider.SetNextFails(0, 0, nil)
}

func TestBockReaderServiceV1(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
