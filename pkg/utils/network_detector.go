package utils

import (
	"errors"
	"strings"
)

func DetectNetworkByAddr(address string) (network string, err error) {
	if strings.HasPrefix(address, "0x") && len(address) == 42 {
		return "ERC20", nil
	}

	if strings.HasPrefix(address, "T") && len(address) == 34 {
		return "TRC20", nil
	}

	err = errors.New("unable to detect the network by address")
	return
}

func DetectNetworkByHash(hash string) (network string, err error) {
	if strings.HasPrefix(hash, "0x") && len(hash) == 66 {
		return "ERC20", nil
	}

	if len(hash) == 64 {
		return "TRC20", nil
	}

	err = errors.New("unable to detect the network by hash")
	return
}
