package io

import (
	"encoding/json"
	"os"
)

type Config struct {
	Extensions []string `json:"extensions"`
}

const configFile = "config.json"

func Default() Config {
	return Config{Extensions: []string{".txt", ".md"}}
}

func LoadConfig() (Config, error) {
	var cfg Config

	file, err := os.ReadFile(configFile)
	if err != nil {
		// If file doesn't exist, return default config
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return cfg, err
	}

	err = json.Unmarshal(file, &cfg)
	return cfg, err
}

func SaveConfig(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}
