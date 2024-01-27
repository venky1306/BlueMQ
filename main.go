package main

import (
	"flag"
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var configFile string
	flag.StringVar(&configFile, "config", "config.yaml", "Path to config file")
	flag.Parse()

	serverConfig, err := initServerConfig(configFile)
	if err != nil {
		panic(err)
	}

	cfg := &Config{
		ProducerListenAddrHttp: ":" + serverConfig.ProducerHttpPort,
		ConsumerListenAddrWS:   ":" + serverConfig.ConsumerWsPort,
		WalPath:                serverConfig.WalPath,
		StoreProducerFunc: func() Storer {
			return NewMemory()
		},
	}

	s, err := NewServer(cfg)
	if err != nil {
		slog.Error("Error creating server", err.Error())
		return
	}
	if err := s.Start(); err != nil {
		slog.Error("Error starting server", err.Error())
		return
	}
}

func initServerConfig(configFile string) (ServerConfig, error) {
	serverConfig, err := ParseServerConfig(configFile)
	if err != nil {
		slog.Error("Error parsing config file", err.Error())
		return serverConfig, err
	}

	if serverConfig.ProducerHttpPort == "" {
		serverConfig.ProducerHttpPort = "5001"
	}

	if serverConfig.ConsumerWsPort == "" {
		serverConfig.ConsumerWsPort = "5002"
	}

	if serverConfig.WalPath == "" {
		serverConfig.WalPath = "wal.aof"
	}

	return serverConfig, nil
}

// Todo:
// planning on leveraging the linux filesystem copy to actually benefit from not
// copying data from kernel space -> user space and back to kernel space? That's
// the reason for kafka 3mil tx /sec.
