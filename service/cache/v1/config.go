package v1

import (
	"fmt"

	"github.com/itiky/eth-block-proxy/service"
	"github.com/spf13/viper"
)

const (
	// Viper configuration keys.
	cfgKeyPrefix = service.ConfigPrefix + "Cache_v1"
	//
	cfgKeyCacheSize  = cfgKeyPrefix + ".CacheSize"
	cfgKeyForkLength = cfgKeyPrefix + ".ForkLength"
	// Configuration defaults.
	defCacheSize  = 100
	defForkLength = 20
)

// Config defines FallbackBlockReaderSvc params.
type Config struct {
	// Cache size [blocks]
	CacheSize int
	// Number of blocks to invalidate on chain forked event [blocks]
	ForkLength int
}

// Validate validates Config values.
func (c Config) Validate() error {
	if c.CacheSize < 0 {
		return fmt.Errorf("CacheSize: must be GTE 0")
	}
	if c.ForkLength < 0 {
		return fmt.Errorf("ForkLength: must be GTE 0")
	}

	return nil
}

// BuildConfig creates a new Config instance from the Viper environment.
func BuildConfig() Config {
	cfg := Config{}
	if err := viper.UnmarshalKey(cfgKeyPrefix, &cfg); err != nil {
		panic(fmt.Errorf("service config unmarshal (%s): %w", cfgKeyPrefix, err))
	}

	return cfg
}

func init() {
	// Set config defaults
	viper.SetDefault(cfgKeyCacheSize, defCacheSize)
	viper.SetDefault(cfgKeyForkLength, defForkLength)
}
