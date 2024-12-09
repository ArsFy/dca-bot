package main

import (
	"encoding/json"
	"log"
	"os"
)

var Config struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	Demo      bool   `json:"demo"`
}

func init() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Panicln("Error reading config file:", err)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Panicln("Error parsing config file:", err)
	}
}
