package main

import (
	"encoding/json"
	"log"
	"os"
)

type suggestionsConfig struct {
	SuggestionsWebhook        string `json:"suggestions_webhook"`
	Database                  string `json:"database"`
	SuggestionsLoggingWebhook string `json:"suggestions_logging_webhook"`
	AuthToken                 string `json:"auth_token"`
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

	if config.AuthToken == "" {
		log.Fatal("auth_token not set in config")
	}

	return &config
}
