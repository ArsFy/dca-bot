package main

import (
	"fmt"
	"math/big"

	bybit "github.com/wuhewuhe/bybit.go.api"
)

var Client *bybit.Client

func main() {
	Client = bybit.NewBybitHttpClient(
		Config.ApiKey, Config.ApiSecret,
		bybit.WithBaseURL(bybit.DEMO_ENV),
	)

	bot1 := &DCA{
		Symbol:      "DOGEUSDT",
		Amount:      big.NewFloat(100.0),
		PriceGap:    big.NewFloat(2.4 * 0.01),
		PriceScale:  big.NewFloat(1.7),
		AmountScale: big.NewFloat(1.5),
		OrderNumber: 6,
	}

	bot1.GetPrice()
	bot1.GetInstrumentsInfo()

	fmt.Println(bot1.TickSize.String(), bot1.QtySize)
}
