package v1

import (
	"context"
	"fmt"
	"testing"

	"github.com/itiky/eth-block-proxy/provider/block/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx context.Context
	r   *BlockNotifierSvcTestResource
	//
	provider *mock.MockBlockProvider
}

func (s *ServiceTestSuite) SetupSuite() {
	r, err := NewBlockNotifierSvcTestResource()
	if err != nil {
		panic(fmt.Errorf("NewBlockNotifierSvcTestResource: %w", err))
	}

	s.ctx = context.TODO()
	s.r = r

	s.provider = r.ReaderSvcRes.ProviderRes.Provider.(*mock.MockBlockProvider)
}

func (s *ServiceTestSuite) SetupTest() {
	s.provider.SetLatestBlockIdxOverride(-1)
}

func TestBockNotifierService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
