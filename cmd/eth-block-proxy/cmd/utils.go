package cmd

import (
	"context"
	"fmt"

	"github.com/itiky/eth-block-proxy/pkg"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Cmd context keys.
var cmdCtxKeyLogger = pkg.ContextKey("cmd.context.logger")

// setCmdCtxLogger adds an instance of zerolog.Logger to cobra.Command context.
func setCmdCtxLogger(cmd *cobra.Command, logger zerolog.Logger) {
	v := cmd.Context().Value(cmdCtxKeyLogger)
	if v == nil {
		panic(fmt.Errorf("%s context: not set", cmdCtxKeyLogger))
	}

	ctxPtr := v.(*zerolog.Logger)
	*ctxPtr = logger
}

// getCmdCtxLogger gets an instance of zerolog.Logger from cobra.Command context.
func getCmdCtxLogger(cmd *cobra.Command) zerolog.Logger {
	v := cmd.Context().Value(cmdCtxKeyLogger)
	if v == nil {
		panic(fmt.Errorf("%s context: not set", cmdCtxKeyLogger))
	}
	logger := v.(*zerolog.Logger)

	return *logger
}

// getBaseCmdCtx creates a context.Context to use with cobra.Command.
func getBaseCmdCtx() context.Context {
	return context.WithValue(context.Background(), cmdCtxKeyLogger, &zerolog.Logger{})
}
