package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addScenario(client *mongo.Client, s Scenario) error {
	if s.Script == "" {
		return errors.New("nead scenario script")
	}
	collection := client.Database("OpenPipelineIO").Collection("scenario")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"script": s.Script})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("The same script of data exist. Change the script of data")
	}
	_, err = collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func getScenario(client *mongo.Client, id string) (Scenario, error) {
	collection := client.Database("OpenPipelineIO").Collection("scenario")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := Scenario{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return s, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func rmScenario(client *mongo.Client, id string) error {
	collection := client.Database("OpenPipelineIO").Collection("scenario")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

func setScenario(client *mongo.Client, s Scenario) error {
	collection := client.Database("OpenPipelineIO").Collection("scenario")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": s.ID},
		bson.D{{Key: "$set", Value: s}},
	)
	if err != nil {
		return err
	}
	return nil
}

func allScenarios(client *mongo.Client) ([]Scenario, error) {
	collection := client.Database("OpenPipelineIO").Collection("scenario")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []Scenario
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}
