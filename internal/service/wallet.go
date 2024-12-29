package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/OwodDEV/crypto-service/internal/models"
	"github.com/OwodDEV/crypto-service/pkg/utils"
)

func (s *Service) GetWallet(ctx context.Context, address string) (resp models.GetWalletResp, err error) {
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
		slog.String("func", "service.GetBalance()"),
	)

	// check for cached balance
	balance, err := s.Cache.GetWalletBalance(ctx, address)
	if err != nil {
		return
	}
	if balance != "" {
		resp.Balance = balance
		return
	}

	// get realtime balance
	network, err := utils.DetectNetworkByAddr(address)
	if err != nil {
		logger.Warn(err.Error(), slog.String("address", address))
		return
	}

	switch network {
	case "ERC20":
		balance, err = s.External.Ethereum.GetBalance(ctx, address, "USDT")
		if err != nil {
			return
		}
	case "TRC20":
		balance, err = s.External.Tron.GetBalance(ctx, address, "USDT")
		if err != nil {
			return
		}
	default:
		err = errors.New("unsupported network")
		logger.Warn(err.Error(), slog.String("network", network))
		if err != nil {
			return
		}
	}

	// save and response
	_ = s.Cache.SaveWalletBalance(ctx, address, balance)
	resp.Balance = balance
	return
}
