package main

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

// Peer is an interface for a peer.
type Peer interface {
	Send([]byte) error
	Close() error
}

// PeerClient is a client for a peer.
type PeerClient struct {
	conn   *websocket.Conn
	server *Server
	mu     sync.Mutex
}

// NewPeer creates a new peer.
func NewPeer(conn *websocket.Conn, s *Server) *PeerClient {
	p := &PeerClient{
		conn:   conn,
		server: s,
	}
	go p.readLoop()
	return p
}

// CMessage is a message from a consumer.
func (p *PeerClient) readLoop() {
	var msg CMessage
	for {
		_, b, err := p.conn.ReadMessage()
		if err != nil {
			if err == websocket.ErrCloseSent {
				slog.Info("Peer closed connection")
				return
			}
			slog.Error("Error reading message", "err", err.Error())
			return
		}
		err = json.Unmarshal(b, &msg)
		if err != nil {
			slog.Error("Error unmarshalling message", "err", err.Error())
			continue
		}
		err = p.handleMsg(msg)
		if err != nil {
			slog.Error("Error handling message", "err", err.Error())
			continue
		}
	}
}

// handleMsg handles a message from a consumer.
// this is where we will handle the message from the consumer.
func (p *PeerClient) handleMsg(msg CMessage) error {
	switch msg.Action {
	case "subscribe":
		slog.Info("Got subscribe message", "topic", msg.Topic)
		for _, t := range msg.Topic {
			p.server.AddSubscriber(t, p)
		}
	case "unsubscribe":
		slog.Info("Got unsubscribe message", "topic", msg.Topic)
		for _, t := range msg.Topic {
			p.server.RemoveSubscriber(t, p)
		}
	case "close":
		slog.Info("Got close message")
		p.server.RemoveConn(p)
	default:
		slog.Error("Got unknown message", "action", msg.Action)
	}
	return nil
}

func (p *PeerClient) Close() error {
	return p.conn.Close()
}

// Send sends a message to a peer.
func (p *PeerClient) Send(msg []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conn.WriteMessage(websocket.TextMessage, msg)
}
