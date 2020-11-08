package main

import (
	"fmt"
	"github.com/spf13/viper"
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

func setDefaults() {
	viper.SetDefault("database-file", "./cfc_suggestions.db")
	viper.SetDefault("ignore-auth", false)
}

func loadConfig() {
	viper.SetConfigName("suggestions_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/cfc_suggestions/")
	viper.AddConfigPath(".")

	setDefaults()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error reading config %w", err))
		}
	}
	viper.AutomaticEnv()
}
