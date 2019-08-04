package dao

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var db *mongo.Database

func init() {
	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		uri = "mongodb://127.0.0.1:27017"
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}

	db = client.Database("auth_manager")

	createIndexes()
}

func createIndexes() {
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := db.Collection(UserCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bsonx.Doc{{
			Key:   "email",
			Value: bsonx.Int32(1),
		}},
		Options: options.Index().SetUnique(true),
	}, opts)

	if err != nil {
		log.Panicf("Error while creating indexes on database: %v", err)
	}
}
