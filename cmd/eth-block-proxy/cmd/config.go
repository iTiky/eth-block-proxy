package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newDefaultConfigCmd creates a new root.default-config cobra.Command.
func newDefaultConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default-config",
		Short: "Saves default config to file",
		Long:  "Useful if you want to change params without ENVs override",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := getCmdCtxLogger(cmd)

			configPath, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				logger.Fatal().Err(err).Msgf("%s flag reading", flagConfig)
			}

			// Configure Viper environment
			viper.SetConfigFile(configPath)
			viper.SetConfigType("toml")

			if err := viper.WriteConfig(); err != nil {
				logger.Info().Msg("Config file write")
			}

			return nil
		},
	}

	cmd.Flags().String(flagConfig, "./config.toml", "config file path")

	return cmd
}
