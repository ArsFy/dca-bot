package main

import (
    "context"
    "dca/utils"
    "errors"
    "math/big"
    "strconv"
)

type DCA struct {
    Symbol      string     `json:"symbol"`
    Amount      *big.Float `json:"amount"`
    Leverage    int        `json:"leverage"`
    PriceGap    *big.Float `json:"price_gap"`
    PriceScale  *big.Float `json:"price_scale"`
    AmountScale *big.Float `json:"amount_scale"`
    OrderNumber int        `json:"order_number"`
}

// GetPrice 獲取價格
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

// SetLeverage 設定杠杆
func (dca *DCA) SetLeverage() {
    _, _ = Client.NewUtaBybitServiceWithParams(map[string]interface{}{
        "category":     "linear",
        "symbol":       dca.Symbol,
        "buyLeverage":  strconv.Itoa(dca.Leverage),
        "sellLeverage": strconv.Itoa(dca.Leverage),
    }).SetPositionLeverage(context.Background())
}

// PlaceMarketOrder 下市價單
func (dca *DCA) PlaceMarketOrder(qty *big.Float) error {
    data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
        "category":    "linear",
        "symbol":      dca.Symbol,
        "side":        "Buy",
        "orderType":   "Market",
        "qty":         qty.String(),
        "positionIdx": "1",
        "reduceOnly":  false,
    }).PlaceOrder(context.Background())

    if err != nil {
        return err
    }
    if data.RetCode != 0 {
        return err
    }

    return nil
}

// PlaceOrder 下限價單
func (dca *DCA) PlaceOrder(price, qty *big.Float) error {
    data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
        "category":    "linear",
        "symbol":      dca.Symbol,
        "side":        "Buy",
        "orderType":   "Limit",
        "price":       price.String(),
        "qty":         qty.String(),
        "positionIdx": "1",
        "reduceOnly":  false,
    }).PlaceOrder(context.Background())
    if err != nil {
        return err
    }
    if data.RetCode != 0 {
        return errors.New(data.RetMsg)
    }
    return nil
}

func (dca *DCA) Run() error {
    amount := big.NewFloat(0).Mul(dca.Amount, big.NewFloat(float64(dca.Leverage)))

    // 當前價格
    nowPrice, err := dca.GetPrice()
    if err != nil {
        return err
    }
    // 第二單加倉價格
    var seEnterPrice *big.Float
    {
        priceDiff := big.NewFloat(0).Mul(nowPrice, dca.PriceGap)
        seEnterPrice = big.NewFloat(0).Sub(nowPrice, priceDiff)
        seEnterPrice = utils.QuantizePrice(seEnterPrice, InstrumentMap[dca.Symbol].TickSize)
    }
    // 價格列表
    var priceList = []*big.Float{
        nowPrice, seEnterPrice,
    }
    for index := 1; index < dca.OrderNumber-1; index++ {
        lastPrice := priceList[len(priceList)-1]

        priceScalePow := big.NewFloat(1)
        for i := 0; i < index; i++ {
            priceScalePow.Mul(priceScalePow, dca.PriceScale)
        }
        priceDiff := big.NewFloat(0).Mul(lastPrice, big.NewFloat(0).Mul(dca.PriceGap, priceScalePow))
        newPrice := big.NewFloat(0).Sub(lastPrice, priceDiff)
        newPrice = utils.QuantizePrice(newPrice, InstrumentMap[dca.Symbol].TickSize)
        priceList = append(priceList, newPrice)
    }

    totalRatio := big.NewFloat(0)
    amountScalePow := big.NewFloat(1)
    for i := 0; i < dca.OrderNumber; i++ {
        totalRatio.Add(totalRatio, amountScalePow)
        amountScalePow.Mul(amountScalePow, dca.AmountScale)
    }

    baseAmount := big.NewFloat(0).Quo(amount, totalRatio)

    var amountList []*big.Float
    amountScalePow.SetFloat64(1)
    for i := 0; i < dca.OrderNumber; i++ {
        thisAmount := big.NewFloat(0).Mul(baseAmount, amountScalePow)
        thisAmount.Quo(thisAmount, priceList[i])
        thisAmount = utils.QuantizePrice(thisAmount, InstrumentMap[dca.Symbol].QtySize)
        amountList = append(amountList, thisAmount)
        amountScalePow.Mul(amountScalePow, dca.AmountScale)
    }

    dca.SetLeverage()

    if err = dca.PlaceMarketOrder(amountList[0]); err != nil {
        return err
    }

    for i := 1; i < dca.OrderNumber; i++ {
        if err = dca.PlaceOrder(priceList[i], amountList[i]); err != nil {
            return err
        }
    }

    return nil
}
