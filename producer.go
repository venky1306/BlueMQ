package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/slog"
)

// Producer is an interface for a producer.
type Producer interface {
	Start() error
}

// ProducerClient is a client for a producer.
type ProducerClient struct {
	listenAddr string
	producech  chan<- Message
}

// ServeHTTP is the http handler for a producer.
func (p *ProducerClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

	if r.Method == "GET" {
	}

	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error(err.Error())
		}
		if len(parts) != 2 {
			fmt.Println("invalid action")
			return
		}
		p.producech <- Message{
			Value: body,
			Topic: parts[1],
		}
	}
}

// NewProducer creates a new producer. It listens on listenAddr and sends messages to producech.
func NewProducer(listenAddr string, producerch chan Message) Producer {
	return &ProducerClient{
		listenAddr: listenAddr,
		producech:  producerch,
	}
}

// Start starts the producer. It listens on listenAddr.
func (p *ProducerClient) Start() error {
	slog.Info(fmt.Sprintf("Listening on %s", p.listenAddr))
	return http.ListenAndServe(p.listenAddr, p)
}
