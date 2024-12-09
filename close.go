package main

import (
    "context"
    "dca/utils"
    "log"
    "math/big"
    "time"
)

type Position struct {
    AvgPrice *big.Float `json:"avg_price"`
    TP       *big.Float `json:"tp"`
}

func ScanPosition() {
    var PositionMap = map[string]Position{}
    var PreviousPositionMap = map[string]bool{}

    for {
        data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
            "category":   "linear",
            "settleCoin": "USDT",
        }).GetPositionList(context.Background())
        if err != nil {
            log.Println("Error getting position list:", err)
            time.Sleep(5 * time.Second)
            continue
        }

        currentPositions := make(map[string]bool)
        for _, v := range data.Result.(map[string]interface{})["list"].([]interface{}) {
            position := v.(map[string]interface{})

            symbol := position["symbol"].(string)

            ap, _ := big.NewFloat(0).SetString(position["avgPrice"].(string))
            tp := big.NewFloat(0).Mul(
                ap,
                big.NewFloat(0).Add(
                    big.NewFloat(1),
                    big.NewFloat(0).Mul(
                        big.NewFloat(SymbolConfig[symbol].TP),
                        big.NewFloat(0.01),
                    ),
                ),
            )
            tp = utils.QuantizePrice(tp, InstrumentMap[symbol].TickSize)

            movePrice := big.NewFloat(0).Mul(
                ap,
                big.NewFloat(0).Mul(
                    big.NewFloat(0).Quo(
                        big.NewFloat(SymbolConfig[symbol].MoveTP),
                        big.NewFloat(float64(SymbolConfig[symbol].Leverage)),
                    ),
                    big.NewFloat(0.01),
                ),
            )
            movePrice = utils.QuantizePrice(movePrice, InstrumentMap[symbol].TickSize)

            currentPositions[symbol] = true

            if pv, ok := PositionMap[symbol]; ok {
                if pv.TP.Cmp(tp) != 0 {
                    PositionMap[symbol] = Position{
                        AvgPrice: ap,
                        TP:       tp,
                    }

                    if err = SetTP(symbol, tp, movePrice); err != nil {
                        log.Println("Error setting TP:", err)
                    }
                }
            } else {
                PositionMap[symbol] = Position{
                    AvgPrice: ap,
                    TP:       tp,
                }

                if err = SetTP(symbol, tp, movePrice); err != nil {
                    log.Println("Error setting TP:", err)
                }
            }
        }

        for symbol := range PreviousPositionMap {
            if _, ok := currentPositions[symbol]; !ok {
                if config, ok := SymbolConfig[symbol]; ok {
                    _ = CancelOrder(symbol)

                    bot := &DCA{
                        Symbol:      symbol,
                        Leverage:    config.Leverage,
                        Amount:      big.NewFloat(float64(config.Amount)),
                        PriceGap:    big.NewFloat(config.PriceGap * 0.01),
                        PriceScale:  big.NewFloat(config.PriceScale),
                        AmountScale: big.NewFloat(config.AmountScale),
                        OrderNumber: config.OrderNumber,
                    }
                    if err = bot.Run(); err != nil {
                        log.Println("Error running bot:", err)
                    }
                }
            }
        }
        PreviousPositionMap = currentPositions

        time.Sleep(5 * time.Second)
    }
}

func SetTP(symbol string, tp, movePrice *big.Float) error {
    _, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
        "category":     "linear",
        "symbol":       symbol,
        "trailingStop": movePrice.String(),
        "activePrice":  tp.String(),
        "positionIdx":  "1",
        "tpslMode":     "Full",
    }).SetPositionTradingStop(context.Background())
    if err != nil {
        return err
    }

    return nil
}

func CancelOrder(symbol string) error {
    data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
        "category": "linear",
        "symbol":   symbol,
        "baseCoin": "USDT",
    }).GetOpenOrders(context.Background())
    if err != nil {
        return err
    }

    var list = make([]map[string]interface{}, 0)

    orderList := data.Result.(map[string]interface{})["list"].([]interface{})
    for _, v := range orderList {
        list = append(list, map[string]interface{}{
            "symbol":  symbol,
            "orderId": v.(map[string]interface{})["orderId"].(string),
        })
    }

    _, err = Client.NewUtaBybitServiceWithParams(map[string]interface{}{
        "category": "linear",
        "request":  list,
    }).CancelBatchOrder(context.Background())
    if err != nil {
        return err
    }

    return nil
}
