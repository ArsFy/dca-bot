package main

import (
    "fmt"
    bybit "github.com/wuhewuhe/bybit.go.api"
    "log"
    "math/big"
)

var Client *bybit.Client

func main() {
    Client = bybit.NewBybitHttpClient(
        Config.ApiKey, Config.ApiSecret,
        bybit.WithBaseURL(bybit.DEMO_ENV),
    )

    // Get Instrument Info
    err := GetInstrumentInfo()
    if err != nil {
        log.Panicln("Error getting instrument info:", err)
    }

    // Get All Position
    GetAllPosition()

    // Scan
    go ScanPosition()

    // Set Enter
    for symbol, config := range SymbolConfig {
        if _, ok := existPositions[symbol]; ok {
            log.Printf("Position already exists for %s, skipping...\n", symbol)
            continue
        }

        bot := &DCA{
            Symbol:      symbol,
            Leverage:    config.Leverage,
            Amount:      big.NewFloat(float64(config.Amount)),
            PriceGap:    big.NewFloat(config.PriceGap * 0.01),
            PriceScale:  big.NewFloat(config.PriceScale),
            AmountScale: big.NewFloat(config.AmountScale),
            OrderNumber: config.OrderNumber,
        }
        if err := bot.Run(); err != nil {
            log.Println("Error running bot:", err)
        }
    }

    fmt.Println("Bot started.")
    select {}
}
