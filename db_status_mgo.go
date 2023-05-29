package main

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetStatus 함수는 Status를 DB에 저장한다.
func SetStatus(client *mongo.Client, s Status) error {
	collection := client.Database(*flagDBName).Collection("status")

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: s.ID}}
	update := bson.D{{Key: "$set", Value: s}}

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = collection.InsertOne(context.TODO(), s)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// GetStatus 함수는 Status를 DB에서 가지고 온다.
func GetStatus(client *mongo.Client, id string) (Status, error) {
	collection := client.Database(*flagDBName).Collection("status")

	var s Status
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&s)
	if err != nil {
		return s, err
	}

	return s, nil
}

// GetInitStatusID 함수는 초기 Status를 DB에서 가지고 온다.
func GetInitStatusID(client *mongo.Client) (string, error) {
	collection := client.Database(*flagDBName).Collection("status")

	count, err := collection.CountDocuments(context.Background(), bson.M{"initstatus": true})
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", errors.New("A status needs to be configured as the initial status for item creation")
	}

	if count != 1 {
		return "", errors.New("The number of initial status values is not equal to 1")
	}

	var s Status
	err = collection.FindOne(context.Background(), bson.M{"initstatus": true}).Decode(&s)
	if err != nil {
		return "", err
	}

	return s.ID, nil
}

// AddStatus 함수는 tasksetting을 DB에 추가한다.
func AddStatus(client *mongo.Client, s Status) error {
	collection := client.Database(*flagDBName).Collection("status")

	count, err := collection.CountDocuments(context.Background(), bson.M{"id": s.ID})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New(s.ID + " Status already exists")
	}

	_, err = collection.InsertOne(context.Background(), s)
	if err != nil {
		return err
	}

	return nil
}

// RmStatus 함수는 Status를 DB에서 삭제한다.
func RmStatus(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("status")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// AllStatus 함수는 모든 Status값을 DB에서 가지고 온다.
func AllStatus(client *mongo.Client) ([]Status, error) {
	collection := client.Database(*flagDBName).Collection("status")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "order", Value: -1}})

	cur, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var results []Status
	for cur.Next(context.Background()) {
		var status Status
		err := cur.Decode(&status)
		if err != nil {
			log.Println(err)
			continue
		}

		results = append(results, status)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return results, nil
}
