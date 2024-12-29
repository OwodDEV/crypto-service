package cache

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func walletBalanceKey(address string) string {
	return "wallet_balance:" + address
}

func (s *Storage) SaveWalletBalance(ctx context.Context, address, balance string) (err error) {
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
		slog.String("func", "storage.cache.SaveWalletBalance()"),
		slog.String("address", address),
	)

	err = s.client.Set(ctx, walletBalanceKey(address), balance, s.walletBalanceTTL).Err()
	if err != nil {
		logger.Error("failed to save wallet balance to cache", slog.Any("error", err))
		return
	}

	logger.Info("successfully saved wallet balance to cache")
	return
}

func (s *Storage) GetWalletBalance(ctx context.Context, address string) (balance string, err error) {
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
		slog.String("func", "storage.cache.GetWalletBalance()"),
		slog.String("address", address),
	)

	balance, err = s.client.Get(ctx, walletBalanceKey(address)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		logger.Error("failed to get wallet balance from cache", slog.Any("error", err))
		return
	}

	logger.Info("successfully loaded wallet balance from cache")
	return
}
