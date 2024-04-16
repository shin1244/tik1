package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient mongo.Client

type Board struct {
	Type string `json:"type"`
	O    []int  `json:"o"`
	X    []int  `json:"x"`
}

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

	conn.WriteJSON(client)
	conn.WriteJSON(NowTikStateMongo(mongoClient))

	for {
		var m Message
		conn.ReadJSON(&m)
		switch m.Type {
		case "user":
			client.state = m.Data.(bool)
		case "board":
			if client.state {
				broadcastTik <- m
			}
		}
	}
}

func HandleTik() {
	mongoClient = MongoOpen()
	defer MongoClose(mongoClient)

	DeleteTikMongo(mongoClient)
	InitTikMongo(mongoClient)

	turn := "X"
	for {
		var win Message
		m := <-broadcastTik
		win.Type = UpdateOneTikMongo(mongoClient, m, turn)

		for client := range clients {
			m.Turn = turn
			client.Conn.WriteJSON(m)
			if win.Type != "" {
				client.Conn.WriteJSON(win)
				client.state = false
			}
		}

		if turn == "X" {
			turn = "O"
		} else {
			turn = "X"
		}
	}
}
