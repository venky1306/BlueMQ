package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Action string   `json:"action"`
	Topic  []string `json:"topic"`
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:5002", nil)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	msg := WSMessage{
		Action: "subscribe",
		Topic:  []string{"test"},
	}
	b, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(b))
	}
}
