package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	clients      []*websocket.Conn
	clientsMutex sync.RWMutex
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(clients)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientsMutex.Lock()
	clients = append(clients, conn)
	clientsMutex.Unlock()

	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("Received message:", string(p))

		for i := range clients {
			err = clients[i].WriteMessage(messageType, p)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.ListenAndServe(":8080", nil)
}
