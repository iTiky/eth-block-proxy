package v1

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"github.com/itiky/eth-block-proxy/service"
	blockReader "github.com/itiky/eth-block-proxy/service/block/reader"
	"github.com/itiky/eth-block-proxy/service/cache"
)

var _ cache.BlockCache = (*CacheSvc)(nil)

// CacheSvc caches block and tx request results and handles the NewBlock and the ChainForked events.
// Events are emitted by the notifier service.
type CacheSvc struct {
	service.BaseSvc
	config Config                  // service config
	reader blockReader.BlockReader // block reader service
	//
	cache          *lru.Cache // thread-safe LRU cache
	latestBlockIdx uint64     // the latest block number (set externally by the ChainForked event)
}

// NewCacheSvc creates a new CacheSvc instance.
func NewCacheSvc(
	baseSvc service.BaseSvc,
	blockReader blockReader.BlockReader,
) (*CacheSvc, error) {

	cfg := BuildConfig()

	if blockReader == nil {
		return nil, fmt.Errorf("blockReader: nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	cache, err := lru.New(cfg.CacheSize)
	if err != nil {
		return nil, fmt.Errorf("creating LRU cache: %w", err)
	}

	svc := &CacheSvc{
		BaseSvc: baseSvc,
		config:  cfg,
		reader:  blockReader,
		cache:   cache,
	}
	svc.SetServiceName("CacheSvc")

	return svc, nil
}
