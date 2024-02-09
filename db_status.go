package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AllStatusV2(client *mongo.Client) ([]Status, error) {
	collection := client.Database(*flagDBName).Collection("status")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "order", Value: -1}})
	var status []Status
	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return status, err
	}
	err = cursor.All(ctx, &status)
	if err != nil {
		return status, err
	}
	return status, nil
}

func GetInitStatusIDV2(client *mongo.Client) (string, error) {
	collection := client.Database(*flagDBName).Collection("status")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := Status{}
	n, err := collection.CountDocuments(ctx, bson.M{"initstatus": true})
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", errors.New("아이템 생성시 초기 Status로 사용할 Status 설정이 필요합니다")
	}
	if n != 1 {
		return "", errors.New("초기 상태 설정값이 1개가 아닙니다")
	}
	err = collection.FindOne(ctx, bson.M{"initstatus": true}).Decode(&s)
	if err != nil {
		return "", err
	}
	return s.ID, nil
}

func AddStatusV2(client *mongo.Client, s Status) error {
	collection := client.Database(*flagDBName).Collection("status")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	num, err := collection.CountDocuments(ctx, bson.M{"id": s.ID})
	if err != nil {
		return err
	}
	if num != 0 {
		err = errors.New(s.ID + " Status가 이미 존재합니다")
		return err
	}
	_, err = collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func RmStatusV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("status")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := collection.FindOneAndDelete(ctx, bson.M{"id": id})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return mongo.ErrNoDocuments
		}
		return result.Err()
	}
	return nil

}

func GetStatusV2(client *mongo.Client, id string) (Status, error) {
	collection := client.Database(*flagDBName).Collection("status")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := Status{}
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func SetStatusV2(client *mongo.Client, s Status) error {
	collection := client.Database(*flagDBName).Collection("status")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": s.ID}
	update := bson.D{{Key: "$set", Value: s}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id" + s.ID)
	}
	return nil
}
