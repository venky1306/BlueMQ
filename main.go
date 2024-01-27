package main

import (
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := &Config{
		ListenAddrHttp: ":5001",
		ListenAddrWS:   ":5002",
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

// planning on leveraging the linux filesystem copy to actually benefit from not
// copying data from kernel space -> user space and back to kernel space? That's
// the reason for kafka 3mil tx /sec.
