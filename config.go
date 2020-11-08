package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func setDefaults() {
	viper.SetDefault("database-file", "./cfc_suggestions.db")
	viper.SetDefault("ignore-auth", false)
	viper.SetDefault("port", "8080")
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
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
