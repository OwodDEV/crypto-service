package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/OwodDEV/crypto-service/internal/config"
	"github.com/OwodDEV/crypto-service/internal/external"
	"github.com/OwodDEV/crypto-service/internal/metrics"
	"github.com/OwodDEV/crypto-service/internal/service"
	"github.com/OwodDEV/crypto-service/internal/storages"
	"github.com/OwodDEV/crypto-service/internal/transport/http"
	"github.com/OwodDEV/crypto-service/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	logger.NewLogger(cfg)
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err := runApp(ctx, cfg)
	done()

	if err != nil {
		slog.Error("application shutdown with an error", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("application shutdown successfully")
	os.Exit(0)
}

func runApp(ctx context.Context, cfg *config.Config) error {
	metrics, err := metrics.NewMetrics(cfg)
	if err != nil {
		return err
	}

	external, err := external.NewExternal(cfg)
	if err != nil {
		return err
	}

	storages, err := storages.NewStorages(cfg)
	if err != nil {
		return err
	}

	srv, err := service.NewService(external, storages, cfg)
	if err != nil {
		return err
	}

	httpServer, err := http.NewServer(srv, metrics, cfg)
	if err != nil {
		return err
	}

	// Running
	errCh := make(chan error, 1)

	err = external.Ethereum.Connect()
	if err != nil {
		return err
	}
	defer external.Ethereum.Shutdown()

	err = external.Tron.Connect()
	if err != nil {
		return err
	}
	defer external.Tron.Shutdown()

	err = storages.Cache.Connect()
	if err != nil {
		return err
	}
	defer storages.Cache.Shutdown()

	go httpServer.Run(errCh)
	defer httpServer.Shutdown()

	// Waiting for stop
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}
