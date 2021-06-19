package testutil

import (
	blockProviderRes "github.com/itiky/eth-block-proxy/provider/block/testutil"
	"github.com/itiky/eth-block-proxy/service/block/reader"
)

// BlockReaderSvcTestResource defines test resource for service containing all the necessary dependencies.
type BlockReaderSvcTestResource struct {
	Svc         reader.BlockReader
	ProviderRes *blockProviderRes.BlockProviderTestResource
}
