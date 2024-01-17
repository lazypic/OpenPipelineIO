package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AllTaskSettingsV2(client *mongo.Client) ([]Tasksetting, error) {
	results := []Tasksetting{}
	collection := client.Database(*flagDBName).Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Find()
	opts.SetSort(bson.M{"id": 1})
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func TaskSettingNamesV2(client *mongo.Client) ([]string, error) {
	var results []string
	collection := client.Database(*flagDBName).Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	values, err := collection.Distinct(ctx, "name", bson.D{})
	if err != nil {
		return results, err
	}
	for _, value := range values {
		results = append(results, fmt.Sprintf("%v", value))
	}
	return results, nil
}

// GetAdminSettingV2 함수는 adminsetting을 DB에서 가지고 온다.
func GetAdminSettingV2(client *mongo.Client) (Setting, error) {
	s := Setting{}
	collection := client.Database(*flagDBName).Collection("admin")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"id": "admin"}).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.ID = "admin"
			_, err = collection.InsertOne(ctx, s)
			if err != nil {
				return s, err
			}
			return s, nil
		}
		return s, err
	}
	return s, nil
}
