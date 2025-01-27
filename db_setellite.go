package main

import (
	"context"
	"time"
	"fmt"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func addSetelliteV2(client *mongo.Client, project string, item Setellite, overwrite bool) error {
	collection := client.Database(*flagDBName).Collection("setellite")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": item.ID}

	// 데이터 존재 여부 확인
	exists, err := documentExists(ctx, collection, filter)
	if err != nil {
		return fmt.Errorf("failed to check document existence: %w", err)
	}

	switch {
	case exists && !overwrite:
		return errors.New("데이터가 이미 존재합니다")
	case exists && overwrite:
		if _, err := collection.UpdateOne(ctx, filter, bson.M{"$set": item}); err != nil {
			return fmt.Errorf("failed to update document: %w", err)
		}
		return nil
	default:
		if _, err := collection.InsertOne(ctx, item); err != nil {
			return fmt.Errorf("failed to insert document: %w", err)
		}
		return nil
	}
}

func documentExists(ctx context.Context, collection *mongo.Collection, filter bson.M) (bool, error) {
	err := collection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
