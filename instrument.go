package main

import (
    "context"
    "math/big"
)

type InstrumentType struct {
    TickSize *big.Float `json:"tick_size"`
    QtySize  *big.Float `json:"qty_size"`
}

var InstrumentMap = map[string]InstrumentType{}

func GetInstrumentInfo() error {
    data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
        "category": "linear",
    }).GetInstrumentInfo(context.Background())
    if err != nil {
        return err
    }

    symbols := data.Result.(map[string]interface{})["list"].([]interface{})

    for _, v := range symbols {
        symbol := v.(map[string]interface{})

        tickSizeData := symbol["priceFilter"].(map[string]interface{})["tickSize"].(string)
        qtySizeData := symbol["lotSizeFilter"].(map[string]interface{})["qtyStep"].(string)

        tickSize, _ := big.NewFloat(0).SetString(tickSizeData)
        qtySize, _ := big.NewFloat(0).SetString(qtySizeData)

        InstrumentMap[symbol["symbol"].(string)] = InstrumentType{
            TickSize: tickSize,
            QtySize:  qtySize,
        }
    }

    return nil
}
