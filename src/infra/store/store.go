// Copyright (c) 2024 Berk Kirtay

package store

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func StoreInstance() *mongo.Client {
	if client == nil {
		client = connect()
		log.Println("MongoDB connection is successfully created.")
	}
	return client
}

func connect() (b *mongo.Client) {
	var connectOnce sync.Once
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
