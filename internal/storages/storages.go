package storages

import (
	"github.com/OwodDEV/crypto-service/internal/config"
	"github.com/OwodDEV/crypto-service/internal/storages/cache"
)

type Storages struct {
	Cache *cache.Storage
}

func NewStorages(cfg *config.Config) (storages *Storages, err error) {
	storages = &Storages{}

	storages.Cache, err = cache.NewStorage(cfg)
	if err != nil {
		return
	}

	return
}
