package utils

import "math/big"

func QuantizePrice(price, tickSize *big.Float) *big.Float {
    quantizedPrice := big.NewFloat(0).Quo(price, tickSize)
    intPart, _ := quantizedPrice.Int(nil)
    quantizedPrice = big.NewFloat(0).SetInt(intPart)
    quantizedPrice = big.NewFloat(0).Mul(quantizedPrice, tickSize)
    return quantizedPrice
}
