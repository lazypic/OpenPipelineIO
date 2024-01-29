package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// allDivisionsV2 함수는 DB에서 전체 Divisions 정보를 가지고오는 함수입니다.
func allDivisionsV2(client *mongo.Client) ([]Division, error) {
	var results []Division
	collection := client.Database(*flagDBName).Collection("divisions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item Division
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// allDepartmentsV2 함수는 DB에서 전체 Departments 정보를 가지고오는 함수입니다.
func allDepartmentsV2(client *mongo.Client) ([]Department, error) {
	var results []Department
	collection := client.Database(*flagDBName).Collection("departments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item Department
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
