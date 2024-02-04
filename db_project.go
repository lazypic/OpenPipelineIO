package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProjectlistV2(client *mongo.Client) ([]string, error) {
	var results []string
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	values, err := collection.Distinct(ctx, "id", bson.D{})
	if err != nil {
		return results, err
	}
	for _, value := range values {
		results = append(results, fmt.Sprintf("%v", value))
	}
	return results, nil
}

// 프로젝트를 추가하는 함수입니다.
func addProjectV2(client *mongo.Client, p Project) error {
	if p.ID == "" {
		return errors.New("빈 문자열입니다. 프로젝트를 생성할 수 없습니다")
	}
	if p.ID == "user" {
		return errors.New("user 이름으로 프로젝트를 생성할 수 없습니다")
	}
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	num, err := collection.CountDocuments(ctx, bson.M{"id": p.ID})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 프로젝트가 존재해서 프로젝트를 생성할 수 없습니다")
	}
	_, err = collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// getStatusProjectsV2 함수는 상태를 받아서 프로젝트 정보를 가지고오는 함수입니다.
func getStatusProjectsV2(client *mongo.Client, status ProjectStatus) ([]Project, error) {
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := []Project{}

	cursor, err := collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// getProjectsV2 함수는 전체 프로젝트 정보를 가지고오는 함수입니다.
func getProjectsV2(client *mongo.Client) ([]Project, error) {
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	results := []Project{}
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}
