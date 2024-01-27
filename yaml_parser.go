package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	ProducerHttpPort string `yaml:"producer_port"`
	ConsumerWsPort   string `yaml:"consumer_port"`
	WalPath          string `yaml:"wal_path"`
}

// ParseServerConfig parses the config file.
func ParseServerConfig(filename string) (ServerConfig, error) {
	var serverConfig ServerConfig
	fname, err := filepath.Abs(filename)

	if err != nil {
		return serverConfig, err
	}

	data, err := os.ReadFile(fname)
	if err != nil {
		return serverConfig, err
	}

	err = yaml.Unmarshal(data, &serverConfig)

	if err != nil {
		return serverConfig, err
	}

	return serverConfig, nil

}
