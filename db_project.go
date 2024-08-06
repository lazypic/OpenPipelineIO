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
		return errors.New("The string is empty. Unable to create the project")
	}
	if p.ID == "user" {
		return errors.New("Unable to create a project with the name 'user'")
	}
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	num, err := collection.CountDocuments(ctx, bson.M{"id": p.ID})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("A project with the same name already exists, so it cannot be created")
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

// OnProjectlistV2 함수는 준비중, 진행중, 백업중인 상태의 프로젝트 리스트만 출력하는 함수입니다.
func OnProjectlistV2(client *mongo.Client) ([]string, error) {
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []string
	projects := []Project{}
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &projects)
	if err != nil {
		return results, err
	}

	for _, p := range projects {
		if p.Status == TestProjectStatus || p.Status == PreProjectStatus || p.Status == PostProjectStatus || p.Status == BackupProjectStatus {
			results = append(results, p.ID)
		}
	}
	return results, nil
}

func getProjectV2(client *mongo.Client, id string) (Project, error) {
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var p Project
	result := collection.FindOne(ctx, bson.M{"id": id})
	if result.Err() == mongo.ErrNoDocuments {
		return p, mongo.ErrNoDocuments
	}
	err := result.Decode(&p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// setProjectV2 함수는 프로젝트 정보를 업데이트하는 함수이다.
func setProjectV2(client *mongo.Client, p Project) error {
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p.Updatetime = time.Now().Format(time.RFC3339)
	filter := bson.M{"id": p.ID}
	update := bson.D{{Key: "$set", Value: p}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id" + p.ID)
	}
	return nil
}

func rmProjectV2(client *mongo.Client, project string) error {
	collection := client.Database(*flagDBName).Collection("project")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"id": project})
	if err != nil {
		return err
	}
	return nil
}
