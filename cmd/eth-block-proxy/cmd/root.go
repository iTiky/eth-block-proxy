package cmd

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	// Cmd flags
	flagLogLevel = "log-level"
)

// Execute prepares cmd Context and executes the root cobra.Command.
func Execute() error {
	return newRootCmd().ExecuteContext(getBaseCmdCtx())
}

// newRootCmd creates a new root cobra.Command.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eth-block-proxy",
		Short: "Ethereum blocks proxy server",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Logger setup
			logWriter := zerolog.NewConsoleWriter()
			logWriter.TimeFormat = time.RFC3339
			logger := zerolog.New(logWriter).Level(zerolog.InfoLevel).With().Timestamp().Logger()

			logLevelBz, err := cmd.Flags().GetString(flagLogLevel)
			if err != nil {
				logger.Fatal().Err(err).Msgf("%s flag reading", flagLogLevel)
			}
			logLevel, err := zerolog.ParseLevel(logLevelBz)
			if err != nil {
				logger.Fatal().Err(err).Msgf("Parsing logLevel (%s)", logLevelBz)
			}
			setCmdCtxLogger(cmd, logger.Level(logLevel))

			return nil
		},
	}

	cmd.PersistentFlags().String(flagLogLevel, "info", "log level [debug,info,warn,error,fatal]")

	cmd.AddCommand(
		newServerCmd(),
		newVersionCmd(),
		newDefaultConfigCmd(),
	)

	return cmd
}
