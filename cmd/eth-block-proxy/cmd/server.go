package cmd

import (
	"os"
	"os/signal"
	"strings"

	"github.com/itiky/eth-block-proxy/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// Cmd flags
	flagConfig = "config"
)

// newServerCmd creates a new root.server cobra.Command.
func newServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Starts the proxy server",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := getCmdCtxLogger(cmd)
			logger.Info().Msgf("App version: %s", GetVersion())

			// Read config
			configPath, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				logger.Fatal().Err(err).Msgf("%s flag reading", flagConfig)
			}

			// Configure Viper environment
			viper.AutomaticEnv()                                   // Enable config params override with ENVs
			viper.SetEnvPrefix("ETH_BLOCK_PROXY")                  // ENVs const prefix
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Config levels ENVs match
			viper.SetConfigFile(configPath)                        // Set default config path
			viper.SetConfigType("toml")                            // Set config format

			if configPath == "" {
				logger.Info().Msg("Config file argument not provided (defaults are used)")
			} else if err := viper.ReadInConfig(); err != nil {
				logger.Fatal().Err(err).Msg("Config file read")
			}

			// Override config parameters with ENVs
			for _, key := range viper.AllKeys() {
				val := viper.Get(key)
				viper.Set(key, val)
			}

			// Start the main app
			app, err := app.NewApp(logger)
			if err != nil {
				logger.Fatal().Err(err).Msg("App: initialization")
			}

			// Wait for signal
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt)
			<-sigCh

			// Stop the main app
			logger.Info().Msg("App: shutting down...")
			if err := app.Stop(); err != nil {
				logger.Fatal().Err(err).Msg("App: stop")
			}

			return nil
		},
	}

	cmd.Flags().String(flagConfig, "", "config file path (defaults are used, if not provided)")

	return cmd
}
