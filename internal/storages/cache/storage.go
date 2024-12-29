package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/OwodDEV/crypto-service/internal/config"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Config           *config.Config
	client           *redis.Client
	walletBalanceTTL time.Duration
}

func NewStorage(cfg *config.Config) (storage *Storage, err error) {
	storage = &Storage{
		Config: cfg,
	}
	storage.walletBalanceTTL = time.Duration(cfg.Storages.Cache.WalletBalanceTTL) * time.Second
	return
}

func (s *Storage) Connect() (err error) {
	slog.Info("initializing Cache storage connection...")
	s.client = redis.NewClient(&redis.Options{
		Addr:     s.Config.Storages.Cache.Host + ":" + s.Config.Storages.Cache.Port,
		Password: s.Config.Storages.Cache.Password,
		DB:       s.Config.Storages.Cache.DBIndex,
	})

	ctx := context.Background()
	_, err = s.client.Ping(ctx).Result()
	if err != nil {
		slog.Error("unable to ping Cache storage", slog.Any("error", err))
		return
	}
	return
}

func (s *Storage) Shutdown() {
	slog.Info("shutting down Cache storage...")
	s.client.Close()
}
