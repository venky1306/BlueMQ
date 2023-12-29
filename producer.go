package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/slog"
)

type Producer interface {
	Start() error
}

type ProducerClient struct {
	listenAddr string
	producech  chan<- Message
}

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

func NewProducer(listenAddr string, producerch chan Message) Producer {
	return &ProducerClient{
		listenAddr: listenAddr,
		producech:  producerch,
	}
}

func (p *ProducerClient) Start() error {
	slog.Info(fmt.Sprintf("Listening on %s", p.listenAddr))
	return http.ListenAndServe(p.listenAddr, p)
}
