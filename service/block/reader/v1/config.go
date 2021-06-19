package v1

import (
	"fmt"
	"time"

	"github.com/itiky/eth-block-proxy/service/block"
	"github.com/spf13/viper"
)

const (
	// Viper configuration keys.
	cfgKeyPrefix = block.GroupConfigPrefix + "ReaderV1"
	//
	cfgKeyRequestTimeoutDur = cfgKeyPrefix + ".RequestTimeoutDur"
	cfgKeyMinRetryDur       = cfgKeyPrefix + ".MinRetryDur"
	cfgKeyMaxRetryDur       = cfgKeyPrefix + ".MaxRetryDur"
	// Configuration defaults.
	defRequestTimeoutDur = 5 * time.Second
	defMinRetryDur       = 50 * time.Millisecond
	defMaxRetryDur       = 500 * time.Millisecond
)

// Config defines FallbackBlockReaderSvc params.
type Config struct {
	// Max duration for block info provider request (0 - retry disabled)
	RequestTimeoutDur time.Duration
	// Min fallback retry duration
	MinRetryDur time.Duration
	// Max fallback retry duration
	MaxRetryDur time.Duration
}

// Validate validates Config values.
func (c Config) Validate() error {
	if c.RequestTimeoutDur < 0 {
		return fmt.Errorf("RequestTimeoutDur: must be GTE 0")
	}
	if c.MinRetryDur <= 0 {
		return fmt.Errorf("MinRetryDur: must be GT 0")
	}
	if c.MaxRetryDur <= 0 {
		return fmt.Errorf("MaxRetryDur: must be GT 0")
	}

	if c.MaxRetryDur <= c.MinRetryDur {
		return fmt.Errorf("MinRetryDur / MaxRetryDur: max must be GT min")
	}
	if c.RequestTimeoutDur != 0 && c.MaxRetryDur > c.RequestTimeoutDur {
		return fmt.Errorf("MaxRetryDur / RequestTimeoutDur: timeout must be GTE max")
	}

	return nil
}

// IsRetryDisabled checks if retry policy is disabled.
func (c Config) IsRetryDisabled() bool {
	return c.RequestTimeoutDur == 0
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
	viper.SetDefault(cfgKeyRequestTimeoutDur, defRequestTimeoutDur)
	viper.SetDefault(cfgKeyMinRetryDur, defMinRetryDur)
	viper.SetDefault(cfgKeyMaxRetryDur, defMaxRetryDur)
}
