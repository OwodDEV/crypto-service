package external

import (
	"github.com/OwodDEV/crypto-service/internal/config"
	"github.com/OwodDEV/crypto-service/internal/external/ethereum"
	"github.com/OwodDEV/crypto-service/internal/external/tron"
)

type External struct {
	Ethereum *ethereum.Ethereum
	Tron     *tron.Tron
}

func NewExternal(cfg *config.Config) (external *External, err error) {
	external = &External{}
	external.Ethereum, err = ethereum.NewEthereumService(cfg)
	if err != nil {
		return
	}

	external.Tron, err = tron.NewTronService(cfg)
	if err != nil {
		return
	}

	return
}
