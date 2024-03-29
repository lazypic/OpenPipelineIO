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

func getShotTaskSettingV2(client *mongo.Client) ([]Tasksetting, error) {
	results := []Tasksetting{}
	collection := client.Database(*flagDBName).Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Find()
	opts.SetSort(bson.M{"name": 1})
	cursor, err := collection.Find(ctx, bson.M{"type": "shot"}, opts)
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

func AddTaskSettingV2(client *mongo.Client, t Tasksetting) error {
	collection := client.Database(*flagDBName).Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"id": t.ID})
	if err != nil {
		return err
	}

	if num > 0 {
		return errors.New("이미 Tasksetting이 존재합니다")
	}

	_, err = collection.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	return nil
}

func getTaskSettingV2(client *mongo.Client, id string) (Tasksetting, error) {
	collection := client.Database(*flagDBName).Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := Tasksetting{}

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func SetTaskSettingV2(client *mongo.Client, t Tasksetting) error {
	collection := client.Database(*flagDBName).Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": t.ID}
	update := bson.D{{Key: "$set", Value: t}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found")
	}
	return nil
}

func RmTaskSetting(client *mongo.Client, name, typ string) error {
	collection := client.Database(*flagDBName).Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"name": name, "type": typ}

	result := collection.FindOneAndDelete(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return mongo.ErrNoDocuments
		}
		return result.Err()
	}
	return nil
}

func TasksettingNamesByExcelOrderV2(client *mongo.Client) ([]string, error) {
	var tasksettings []Tasksetting
	var results []string

	collection := client.Database(*flagDBName).Collection("tasksetting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.M{"excelorder": 1})
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &tasksettings)
	if err != nil {
		return results, err
	}
	for _, t := range tasksettings {
		results = append(results, t.Name)
	}
	return Unique(results), nil
}
