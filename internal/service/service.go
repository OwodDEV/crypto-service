package service

import (
	"context"

	"github.com/OwodDEV/crypto-service/internal/config"
	"github.com/OwodDEV/crypto-service/internal/external"
	"github.com/OwodDEV/crypto-service/internal/storages"
)

type Service struct {
	Config   *config.Config
	External *external.External

	Cache Cache
}

type Cache interface {
	SaveWalletBalance(ctx context.Context, address, balance string) (err error)
	GetWalletBalance(ctx context.Context, address string) (balance string, err error)
}

func NewService(external *external.External, storages *storages.Storages, cfg *config.Config) (service *Service, err error) {
	service = &Service{
		Config:   cfg,
		External: external,
		Cache:    storages.Cache,
	}
	return
}
