package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetPublishKeyV2(client *mongo.Client, key PublishKey) error {
	collection := client.Database(*flagDBName).Collection("publishkey")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": key.ID}

	update := bson.M{"$set": key}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		// If no document was matched, insert the key
		_, err := collection.InsertOne(ctx, key)
		if err != nil {
			return err
		}
	}

	return nil
}


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

func GetPublishKeyV2(client *mongo.Client, id string) (PublishKey, error) {
	collection := client.Database(*flagDBName).Collection("publishkey")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	key := PublishKey{}

	err := collection.FindOne(ctx, filter).Decode(&key)
	if err != nil {
		return key, err
	}

	return key, nil
}

func addTaskPublishV2(client *mongo.Client, id, task, key string, p Publish) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}

	// Ensure the task structure exists
	initUpdate := bson.M{
		"$set": bson.M{
			fmt.Sprintf("tasks.%s.publishes", task): bson.M{},
		},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"id": id}, initUpdate)
	if err != nil {
		return err
	}

	// Add the publish key
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{"$push": bson.M{fmt.Sprintf("tasks.%s.publishes.%s", task, key): p}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id " + id)
	}
	return nil
}

func AddPublishKeyV2(client *mongo.Client, key PublishKey) error {
	collection := client.Database(*flagDBName).Collection("publishkey")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n, err := collection.CountDocuments(ctx, bson.M{"id": key.ID})
	if err != nil {
		return err
	}
	
	if n > 0 {
		return errors.New(key.ID + " PublishKey가 이미 존재합니다")
	}
	_, err = collection.InsertOne(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func AllPublishKeysV2(client *mongo.Client) ([]PublishKey, error) {
	collection := client.Database(*flagDBName).Collection("publishkey")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var keys []PublishKey
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "id", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &keys); err != nil {
		return nil, err
	}

	return keys, nil
}

func RmPublishKeyV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("publishkey")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
