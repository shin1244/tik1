package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Type  string          `json:"type"`
	Conn  *websocket.Conn `json:"conn"`
	state bool            `json:"state"`
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Turn string      `json:"turn"`
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*Client]bool)
var broadcastTik = make(chan Message)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Failed: ", err)
		return
	}
	var client Client
	client.Conn = conn
	client.Type = "user"
	clients[&client] = true
	defer delete(clients, &client)

	conn.WriteJSON(client)

	for {
		var m Message
		conn.ReadJSON(&m)
		switch m.Type {
		case "user":
			client.state = m.Data.(bool)
			fmt.Println(client)
		case "board":
			if client.state {
				broadcastTik <- m
			}
		}
	}
}

func HandleTik() {
	for {
		m := <-broadcastTik
		for client := range clients {
			err := client.Conn.WriteJSON(m)
			if err != nil {
				fmt.Println("Failed to send message to client: ", err)
				delete(clients, client)
				client.Conn.Close()
			}
		}
	}
}
