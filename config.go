package main

import (
	"encoding/json"
	"log"
	"os"
)

type suggestionsConfig struct {
	SuggestionsChannel        string `json:"suggestions_channel"`
	Database                  string `json:"database"`
	SuggestionsLoggingChannel string `json:"suggestions_logging_channel"`
	AuthToken                 string `json:"auth_token"`
	BotToken                  string `json:"bot_token"`
	IgnoreAuth                bool   `json:"ignore_auth"`
	LogSql                    bool   `json:"log_sql"`
	SentryDSN                 string `json:"sentry_dsn"`
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
