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
	fmt.Println(err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.ID = "admin"
			_, err = collection.InsertOne(ctx, s)
			return s, err
		}
		return s, err
	}
	return s, nil
}

// SetAdminsettingV2 함수는 사용자 정보를 업데이트하는 함수이다.
func SetAdminSettingV2(client *mongo.Client, s Setting) error {
	collection := client.Database(*flagDBName).Collection("admin")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": "admin"}
	update := bson.D{{Key: "$set", Value: s}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found")
	}
	return nil
}
