package main

import (
    "context"
    "log"
)

var existPositions = map[string]bool{}

func GetAllPosition() {
    data, err := Client.NewUtaBybitServiceWithParams(map[string]interface{}{
        "category":   "linear",
        "settleCoin": "USDT",
    }).GetPositionList(context.Background())
    if err != nil {
        log.Println("Error getting position list:", err)
    }

    for _, v := range data.Result.(map[string]interface{})["list"].([]interface{}) {
        existPositions[v.(map[string]interface{})["symbol"].(string)] = true
    }
}
