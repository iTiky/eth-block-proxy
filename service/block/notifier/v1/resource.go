package v1

import (
	blockReaderRes "github.com/itiky/eth-block-proxy/service/block/reader/testutil"
)

// BlockNotifierSvcTestResource defines test resource for service containing all the necessary dependencies.
type BlockNotifierSvcTestResource struct {
	Svc          *BlockNotifierSvc
	ReaderSvcRes *blockReaderRes.BlockReaderSvcTestResource
}
