package main

type Consumer interface {
	Start() error
}

type ConsumerClient struct {
	listenAddr string
	consumech  chan Request
}

func NewConsumer(listenAddr string, consumech chan Request) Consumer {
	return &ConsumerClient{
		listenAddr: listenAddr,
		consumech:  consumech,
	}
}

func (c *ConsumerClient) Start() error {
	return nil
}
