package utils

import (
	"fmt"
	"math/big"
)

func FormatCurrency(value *big.Int, tokenDecimals int) string {
	valueFloat := new(big.Float).SetInt(value)
	decimalFactor := new(big.Float).SetInt(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(int64(tokenDecimals)), nil))
	valueFloat = valueFloat.Quo(valueFloat, decimalFactor)
	format := fmt.Sprintf("%%.%df", tokenDecimals)
	return fmt.Sprintf(format, valueFloat)
}
