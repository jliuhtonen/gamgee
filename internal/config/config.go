package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Port          int
	UpstreamDns   string
	BlocklistURIs []string
}

const configFile = "./config.toml"

func ReadConfig() (*Config, error) {
	var config Config

	_, err := toml.DecodeFile(configFile, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
