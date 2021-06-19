package app

import (
	"github.com/spf13/viper"
)

const (
	// Viper configuration keys.
	cfgKeyPrefix = "app"
	//
	cfgKeyHost        = cfgKeyPrefix + ".Host"
	cfgKeyServicePort = cfgKeyPrefix + ".ServicePort"
	// Configuration defaults.
	defHost        = "0.0.0.0"
	defServicePort = "2412"
)

func init() {
	// Set config defaults
	viper.SetDefault(cfgKeyHost, defHost)
	viper.SetDefault(cfgKeyServicePort, defServicePort)
}
