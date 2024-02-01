package main

import (
	"fmt"
	"sync"

	"golang.org/x/exp/slog"
)

// Message is a message.
type Message struct {
	Topic string
	Value []byte
}

type Config struct {
	ProducerListenAddrHttp string
	ConsumerListenAddrWS   string
	WalPath                string
	StoreProducerFunc      StoreProducerFunc
}

// Server struct is the main struct for the server.
type Server struct {
	*Config

	topics map[string]Storer

	mu    sync.Mutex
	peers map[Peer]bool

	subscribers map[string][]Peer

	consumers  []Consumer
	producers  []Producer
	producech  chan Message
	quitsignal chan struct{}
}

// NewServer creates a new server.
func NewServer(c *Config) (*Server, error) {
	producerch := make(chan Message)
	s := &Server{
		Config:      c,
		topics:      make(map[string]Storer), // topics mapping to corresponding Storer obj.
		peers:       make(map[Peer]bool),
		subscribers: make(map[string][]Peer),
		producech:   producerch,
		quitsignal:  make(chan struct{}),
		consumers:   []Consumer{},
		producers: []Producer{
			NewProducer(c.ProducerListenAddrHttp, producerch),
		},
	}
	s.consumers = append(s.consumers, NewConsumer(c.ConsumerListenAddrWS, s))
	return s, nil
}

// Start starts the server.
func (s *Server) Start() error {
	for _, p := range s.producers {
		// fmt.Println("starting producer")
		go p.Start()
	}
	for _, c := range s.consumers {
		go c.Start()
	}
	for {
		select {
		case <-s.quitsignal:
			return nil
		case m := <-s.producech:
			index, err := s.publish(m)
			if err != nil {
				slog.Error("Error publishing message: %v", err)
			} else {
				slog.Info(fmt.Sprintf("Published message to topic %s with index %d", m.Topic, index))
			}
		}
	}
}

// Publish publishes a message.
func (s *Server) publish(msg Message) (int, error) {
	storer, ok := s.topics[msg.Topic]
	if !ok {
		storer = s.StoreProducerFunc()
		s.topics[msg.Topic] = storer
	}
	val, err := storer.Put(msg.Value)
	if err != nil {
		return 0, err
	}
	for _, p := range s.subscribers[msg.Topic] {
		go p.Send(msg.Value)
	}
	return val, nil
}

// Stop stops the server.
func (s *Server) Stop() {
	close(s.quitsignal)
}

// AddConn adds a connection.
func (s *Server) AddConn(p Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Info("Adding connection", "peer", fmt.Sprint(p))
	s.peers[p] = true
}

// RemoveConn removes a connection.
func (s *Server) RemoveConn(p Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Info("Removing connection", "peer", fmt.Sprint(p))
	delete(s.peers, p)
	for topic := range s.subscribers {
		go s.RemoveSubscriber(topic, p)
	}
	err := p.Close()
	slog.Error("Err Closing connection.", fmt.Sprint(err.Error()))
}

// AddSubscriber adds a subscriber.
func (s *Server) AddSubscriber(topic string, p Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Info("Adding subscriber", "topic", topic, "peer", fmt.Sprint(p))
	s.subscribers[topic] = append(s.subscribers[topic], p)
}

// RemoveSubscriber removes a subscriber.
func (s *Server) RemoveSubscriber(topic string, p Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Info("Removing subscriber", "topic", topic, "peer", fmt.Sprint(p))
	for i, peer := range s.subscribers[topic] {
		if peer == p {
			s.subscribers[topic] = append(s.subscribers[topic][:i], s.subscribers[topic][i+1:]...)
			return
		}
	}
}
