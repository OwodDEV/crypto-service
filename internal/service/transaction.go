package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/OwodDEV/crypto-service/internal/models"
	"github.com/OwodDEV/crypto-service/pkg/utils"
)

func (s *Service) GetTransaction(ctx context.Context, hash string) (resp models.GetTransactionResp, err error) {
	logger := slog.With(
		slog.String("request_id", ctx.Value("request_id").(string)),
		slog.String("func", "service.GetTransaction()"),
	)

	network, err := utils.DetectNetworkByHash(hash)
	if err != nil {
		logger.Warn(err.Error(), slog.String("hash", hash))
		return
	}

	var trxData models.Transaction
	switch network {
	case "ERC20":
		trxData, err = s.External.Ethereum.GetTransaction(ctx, hash, "USDT")
		if err != nil {
			return resp, err
		}
	case "TRC20":
		trxData, err = s.External.Tron.GetTransaction(ctx, hash, "USDT")
		if err != nil {
			return resp, err
		}
	default:
		err = errors.New("unsupported network")
		logger.Warn(err.Error(), slog.String("network", network))
		if err != nil {
			return
		}
	}

	resp = models.GetTransactionResp{
		From:   trxData.From,
		To:     trxData.To,
		Amount: trxData.Amount,
	}

	return
}
