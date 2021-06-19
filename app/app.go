package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/itiky/eth-block-proxy/api/rest"
	blockProvider "github.com/itiky/eth-block-proxy/provider/block"
	"github.com/itiky/eth-block-proxy/provider/block/cloudflare"
	"github.com/itiky/eth-block-proxy/service"
	notifSvcV1 "github.com/itiky/eth-block-proxy/service/block/notifier/v1"
	blockReader "github.com/itiky/eth-block-proxy/service/block/reader"
	readerSvcV1 "github.com/itiky/eth-block-proxy/service/block/reader/v1"
	blockCache "github.com/itiky/eth-block-proxy/service/cache"
	cacheSvcV1 "github.com/itiky/eth-block-proxy/service/cache/v1"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type App struct {
	blockProvider  blockProvider.BlockProvider
	blockReaderSvc blockReader.BlockReader
	blockCacheSvc  blockCache.BlockCache
	blockNotifSvc  *notifSvcV1.BlockNotifierSvc
	apiServer      *http.Server
}

// Stop gracefully shuts down all the App's services and servers.
func (app *App) Stop() error {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownCancel()

	if err := app.apiServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("apiServer: %w", err)
	}
	if err := app.blockNotifSvc.Close(); err != nil {
		return fmt.Errorf("blockNotifSvc: %w", err)
	}
	if err := app.blockCacheSvc.Close(); err != nil {
		return fmt.Errorf("blockCacheSvc: %w", err)
	}
	if err := app.blockReaderSvc.Close(); err != nil {
		return fmt.Errorf("blockReaderSvc: %w", err)
	}
	if err := app.blockProvider.Close(); err != nil {
		return fmt.Errorf("blockProvider: %w", err)
	}

	return nil
}

// NewApp initializes all the dependencies and starts the API server.
func NewApp(logger zerolog.Logger) (*App, error) {
	// Providers
	blockProvider, err := cloudflare.NewCloudflareBlockProvider()
	if err != nil {
		return nil, fmt.Errorf("BlockProvider (CloudFlare): %w", err)
	}

	// Services
	baseSvc := service.NewBaseSvc(logger)

	blockReaderSvc, err := readerSvcV1.NewFallbackBlockReaderSvc(baseSvc, blockProvider)
	if err != nil {
		return nil, fmt.Errorf("BlockReader (v1): %w", err)
	}

	blockCacheSvc, err := cacheSvcV1.NewCacheSvc(baseSvc, blockReaderSvc)
	if err != nil {
		return nil, fmt.Errorf("BlockCache (v1): %w", err)
	}

	blockNotifSvc, err := notifSvcV1.NewBlockNotifierSvc(baseSvc, blockReaderSvc, blockCacheSvc.HandleNewBlockEvent, blockCacheSvc.HandleChainForkedEvent)
	if err != nil {
		return nil, fmt.Errorf("BlockNotifierSvc (v1): %w", err)
	}

	// API servers
	router := rest.NewRouter(logger, blockCacheSvc)

	srvAddr := net.JoinHostPort(viper.GetString(cfgKeyHost), viper.GetString(cfgKeyServicePort))
	srv := &http.Server{
		Addr:    srvAddr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal().Err(err).Msg("API server")
			}
		}
	}()
	logger.Info().Msgf("API server started: %s", srvAddr)

	return &App{
		blockProvider:  blockProvider,
		blockReaderSvc: blockReaderSvc,
		blockCacheSvc:  blockCacheSvc,
		blockNotifSvc:  blockNotifSvc,
		apiServer:      srv,
	}, nil
}
