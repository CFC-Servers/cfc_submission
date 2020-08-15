package main

import (
	"encoding/json"
	"os"
)

type suggestionsConfig struct {
	SuggestionsWebhook string `json:"suggestions_webhook"`
	Database string `json:"database"`
}

func loadConfig(filename string) *suggestionsConfig {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var config suggestionsConfig
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&config); err != nil {
		panic(err)
	}
	return &config
}
