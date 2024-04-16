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
	State string          `json:"state"`
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
var TikChan = make(chan Message)
var ClearChan = make(chan bool)
var PlayerChan = make(chan Message)

var turn = "X"
var PM = make(chan int)
var count = 0

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Failed: ", err)
		return
	}
	var client Client
	client.Conn = conn
	clients[&client] = true
	PM <- 1

	defer func() {
		PM <- -1
		delete(clients, &client)
		conn.Close()
	}()

	conn.WriteJSON(NowTikStateMongo(mongoClient))

	for {
		var m Message
		conn.ReadJSON(&m)
		switch m.Type {
		case "user":
			client.State = m.Data.(string)
			PlayerChan <- Message{Type: "player", Data: client.State}
		case "board":
			if client.State == turn {
				TikChan <- m
			}
		case "clear":
			DeleteTikMongo(mongoClient)
			InitTikMongo(mongoClient)
			ClearChan <- true
		}
	}
}

func HandleTik() {
	mongoClient = MongoOpen()

	DeleteTikMongo(mongoClient)
	InitTikMongo(mongoClient)

	for {
		var win Message
		m := <-TikChan
		win.Type = UpdateOneTikMongo(mongoClient, m, turn)

		for client := range clients {
			m.Turn = turn
			client.Conn.WriteJSON(m)
			if win.Type != "" {
				client.Conn.WriteJSON(win)
				client.State = ""
			}
		}

		if turn == "X" {
			turn = "O"
		} else {
			turn = "X"
		}
	}
}

func ClearTik() {
	<-ClearChan
	turn = "X"
	for client := range clients {
		client.State = ""
		client.Conn.WriteJSON(NowTikStateMongo(mongoClient))
		client.Conn.WriteJSON(Message{Type: "clear"})
	}
}

func UserCount() {
	for {
		now := <-PM
		count += now
		for client := range clients {
			err := client.Conn.WriteJSON(Message{Type: "userCount", Data: count})
			if err != nil {
				fmt.Println("Failed to send message to client: ", err)
				delete(clients, client)
				client.Conn.Close()
			}
		}
	}
}

func PlayerTik() {
	for {
		m := <-PlayerChan
		for client := range clients {
			client.Conn.WriteJSON(m)
		}
	}
}
