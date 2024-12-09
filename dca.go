package main

import (
	"context"
	"math/big"
	"strconv"
)

type DCA struct {
	Symbol      string     `json:"symbol"`
	Amount      *big.Float `json:"amount"`
	PriceGap    *big.Float `json:"price_gap"`
	PriceScale  *big.Float `json:"price_scale"`
	AmountScale *big.Float `json:"amount_scale"`
	OrderNumber int        `json:"order_number"`

	TickSize *big.Float `json:"tick_size"`
	QtySize  *big.Float `json:"qty_size"`
}

// 獲取價格
func (dca *DCA) GetPrice() (*big.Float, error) {
	data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
		"category": "linear",
		"symbol":   dca.Symbol,
	}).GetMarketTickers(context.Background())
	if err != nil {
		return nil, err
	}
	price, _ := strconv.ParseFloat(data.Result.(map[string]interface{})["list"].([]interface{})[0].(map[string]interface{})["lastPrice"].(string), 64)

	return big.NewFloat(price), nil
}

// 獲取合約
func (dca *DCA) GetInstrumentsInfo() error {
	data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
		"category": "linear",
		"symbol":   dca.Symbol,
	}).GetInstrumentInfo(context.Background())
	if err != nil {
		return err
	}

	info := data.Result.(map[string]interface{})["list"].([]interface{})[0].(map[string]interface{})

	tickSizeData := info["priceFilter"].(map[string]interface{})["tickSize"].(string)
	qtySizeData := info["lotSizeFilter"].(map[string]interface{})["qtyStep"].(string)

	dca.TickSize, _ = big.NewFloat(0).SetString(tickSizeData)
	dca.QtySize, _ = big.NewFloat(0).SetString(qtySizeData)

	return nil
}
