package testutil

import (
	blockReaderRes "github.com/itiky/eth-block-proxy/service/block/reader/testutil"
	"github.com/itiky/eth-block-proxy/service/cache"
)

// BlockCacheSvcTestResource defines test resource for service containing all the necessary dependencies.
type BlockCacheSvcTestResource struct {
	Svc          cache.BlockCache
	ReaderSvcRes *blockReaderRes.BlockReaderSvcTestResource
}
