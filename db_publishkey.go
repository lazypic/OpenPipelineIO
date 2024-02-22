package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HasPublishKeyV2(client *mongo.Client, id string) bool {
	collection := client.Database(*flagDBName).Collection("publishkey")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := collection.CountDocuments(ctx, bson.M{"id": id})
	if err != nil {
		return false
	}
	if n == 0 {
		return false
	}
	return true
}

func addTaskPublishV2(client *mongo.Client, id, task, key string, p Publish) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}

	result, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$push": bson.M{fmt.Sprintf("tasks.%s.publishes.%s", task, key): p}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id" + id)
	}
	return nil
}
