package main

import (
	"fmt"

	"golang.org/x/exp/slog"
)

type Message struct {
	Topic string
	Value []byte
}

type Request struct {
	Topic string
	Index int
}

type Config struct {
	ListenAddr        string
	StoreProducerFunc StoreProducerFunc
}

type Server struct {
	*Config

	topics map[string]Storer

	// consumers []Consumer
	producers  []Producer
	producech  chan Message
	consumech  chan Request
	quitsignal chan struct{}
}

func NewServer(c *Config) (*Server, error) {
	producerch := make(chan Message)
	consumech := make(chan Request)
	return &Server{
		Config:     c,
		topics:     make(map[string]Storer),
		producech:  producerch,
		consumech:  consumech,
		quitsignal: make(chan struct{}),
		// consumers: []Consumer{
		// 	NewConsumer(c.ListenAddr, consumech),
		// },
		producers: []Producer{
			NewProducer(c.ListenAddr, producerch),
		},
	}, nil
}

func (s *Server) Start() error {
	for _, p := range s.producers {
		// fmt.Println("starting producer")
		go p.Start()
	}
	// // for _, c := range s.consumers {
	// // 	go c.Start()
	// // }
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
			// case index := <-s.consumech:
		}
	}
}

func (s *Server) publish(msg Message) (int, error) {
	storer, ok := s.topics[msg.Topic]
	if !ok {
		storer = s.StoreProducerFunc()
		s.topics[msg.Topic] = storer
	}
	return storer.Put(msg.Value)
}

func (s *Server) consume(index int) (Message, error) {
	return Message{}, nil
}
