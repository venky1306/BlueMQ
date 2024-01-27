package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

var upgrader = websocket.Upgrader{}

// Consumer is an interface for a consumer.
type Consumer interface {
	Start() error
}

// ConsumerClient is a client for a consumer.
type ConsumerClient struct {
	listenAddr string
	server     *Server
}

// NewConsumer creates a new consumer.
func NewConsumer(listenAddr string, server *Server) Consumer {
	return &ConsumerClient{
		listenAddr: listenAddr,
		server:     server,
	}
}

// CMessage is a message from a consumer.
type CMessage struct {
	Action string   `json:"action"`
	Topic  []string `json:"topic"`
}

// Start starts the consumer.
func (c *ConsumerClient) Start() error {
	slog.Info(fmt.Sprintf("Listening on %s", c.listenAddr))
	return http.ListenAndServe(c.listenAddr, c)
}

// ServeHTTP serves the consumer. Upgrades the connection to a websocket.
func (c *ConsumerClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	// defer conn.Close()
	if err != nil {
		slog.Error("Error upgrading connection", err.Error())
		return
	}
	p := NewPeer(conn, c.server)
	c.server.AddConn(p)
}
