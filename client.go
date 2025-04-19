package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	// WebSocket connection
	conn *websocket.Conn

	// Outgoing messages to be sent
	send chan []byte
}

// readPump reads messages from the WebSocket and sends them to the hub
func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var payload map[string]string
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Println("json decode error:", err)
			continue
		}

		formatted := []byte(payload["username"] + ": " + payload["message"])
		hub.broadcast <- formatted

	}
}

// writePump sends messages from the hub to the client
func (c *Client) writePump() {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
