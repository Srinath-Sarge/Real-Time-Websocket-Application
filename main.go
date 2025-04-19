package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Web-Socket Upgrade Error", err.Error())
		return
	}
	defer conn.Close()

	client := &Client{conn: conn, send: make(chan []byte)}
	hub.register <- client

	go client.writePump()
	client.readPump()

}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go hub.run()

	fmt.Println("Server Started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
