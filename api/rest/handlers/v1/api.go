package v1

import (
	"github.com/itiky/eth-block-proxy/pkg"
	"github.com/itiky/eth-block-proxy/service/cache"
)

const (
	// URL param keys.
	UrlParamBlock  = "blockNumber"
	UrlParamTxHash = "txHash"

	UrlParamBlockPattern  = "/{" + UrlParamBlock + ":(^latest$)|(^[0-9]+$)}" // latest OR blockNumber (dec)
	UrlParamTxHashPattern = "/{" + UrlParamTxHash + ":^(0x)?[a-fA-F0-9]+$}"
)

// Context middleware keys,
var CtxKeyBlock = pkg.ContextKey("block")

// ApiHandlers keeps API handlers dependencies and helper functions.
type ApiHandlers struct {
	blockCacheSvc cache.BlockCache
}

// NewApiHandlers creates a new ApiHandlers object.
func NewApiHandlers(blockCacheSvc cache.BlockCache) ApiHandlers {
	return ApiHandlers{
		blockCacheSvc: blockCacheSvc,
	}
}
