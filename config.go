package main

import (
    "encoding/json"
    "log"
    "os"
)

type SymbolType struct {
    Amount      int64   `json:"amount"`
    Leverage    int     `json:"leverage"`
    PriceGap    float64 `json:"price_gap"`
    PriceScale  float64 `json:"price_scale"`
    AmountScale float64 `json:"amount_scale"`
    OrderNumber int     `json:"order_number"`
    TP          float64 `json:"tp"`
    MoveTP      float64 `json:"move_tp"`
}

var Config struct {
    ApiKey    string `json:"api_key"`
    ApiSecret string `json:"api_secret"`
    Demo      bool   `json:"demo"`
}

var SymbolConfig = map[string]SymbolType{}

func init() {
    // Config
    file, err := os.ReadFile("./config/config.json")
    if err != nil {
        log.Panicln("Error reading config file:", err)
    }
    err = json.Unmarshal(file, &Config)
    if err != nil {
        log.Panicln("Error parsing config file:", err)
    }

    // SymbolConfig
    file, err = os.ReadFile("./config/symbol.json")
    if err != nil {
        log.Panicln("Error reading symbol file:", err)
    }
    err = json.Unmarshal(file, &SymbolConfig)
    if err != nil {
        log.Panicln("Error parsing symbol file:", err)
    }
}
