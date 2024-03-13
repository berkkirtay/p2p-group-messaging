package infrastructure

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatastore struct {
	Session *mongo.Client
}

func NewDatastore() *mongo.Client {
	client := connect()
	log.Println("MongoDB connection is successfully created...")
	return client
}

func connect() (b *mongo.Client) {
	var connectOnce sync.Once
	var client *mongo.Client
	var err error
	connectOnce.Do(func() {
		client, err = mongo.Connect(context.TODO(),
			options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
	})
	return client
}
