package app

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
mongoDB 클라이언트 연결
*/
func MongoOpen() mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// MongoDB 연결
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return *client
}

/*
mongoDB 클라이언트 종료
*/
func MongoClose(client mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func InitTikMongo(client mongo.Client) {
	collection := client.Database("monGo").Collection("tik")

	board := bson.D{
		{"X", bson.A{}},
		{"O", bson.A{}},
	}

	_, err := collection.InsertOne(context.TODO(), board)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateOneTikMongo(client mongo.Client, tik Message, turn string) string {
	data, _ := strconv.Atoi(tik.Data.(string))
	filter := bson.D{}
	update := bson.D{{"$addToSet", bson.D{{turn, data}}}}
	collection := client.Database("monGo").Collection("tik")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	nowBoard := collection.FindOne(context.TODO(), filter)
	var board Board
	board.Type = "nowBoard"

	nowBoard.Decode(&board)
	return checkWinner(board)
}

func DeleteTikMongo(client mongo.Client) {
	filter := bson.D{}
	collection := client.Database("monGo").Collection("tik")

	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
}

func NowTikStateMongo(client mongo.Client) Board {
	filter := bson.D{}
	collection := client.Database("monGo").Collection("tik")

	nowBoard := collection.FindOne(context.TODO(), filter)

	var board Board
	board.Type = "nowBoard"

	nowBoard.Decode(&board)

	return board
}
