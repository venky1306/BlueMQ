package main

type Producer interface {
	Start() error
}

type Consumer interface {
	Start() error
}

type ProducerClient struct {
	listenAddr string
	producech  chan<- Message
}

func (p *ProducerClient) Start() error {
	return nil
}
