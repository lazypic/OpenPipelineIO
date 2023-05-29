package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func connectToMongoDB(dbIP string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbIP))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func disconnectFromMongoDB(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func getServerVersion(client *mongo.Client) (string, error) {
	buildInfo, err := client.Database("admin").RunCommand(context.Background(), bson.D{{Key: "buildInfo", Value: 1}}).DecodeBytes()
	if err != nil {
		return "", err
	}

	var result struct {
		Version string `bson:"version"`
	}

	err = bson.Unmarshal(buildInfo, &result)
	if err != nil {
		return "", err
	}

	return result.Version, nil
}
