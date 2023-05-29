package main

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Projectlist 함수는 프로젝트 리스트를 출력하는 함수입니다.
func Projectlist(client *mongo.Client) ([]string, error) {
	var results []string
	collection := client.Database(*flagDBName).Collection("project")

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var project Project
		err := cur.Decode(&project)
		if err != nil {
			log.Println(err)
			continue
		}

		results = append(results, project.ID)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// OnProjectlist 함수는 준비중, 진행중, 백업중인 상태의 프로젝트 리스트만 출력하는 함수입니다.
func OnProjectlist(client *mongo.Client) ([]string, error) {
	collection := client.Database(*flagDBName).Collection("project")

	filter := bson.M{}
	sort := bson.D{{"id", 1}}

	cur, err := collection.Find(context.Background(), filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var projects []Project
	if err := cur.All(context.Background(), &projects); err != nil {
		return nil, err
	}

	var results []string
	for _, p := range projects {
		if p.Status == TestProjectStatus || p.Status == PreProjectStatus || p.Status == PostProjectStatus || p.Status == BackupProjectStatus {
			results = append(results, p.ID)
		}
	}

	return results, nil
}

// 프로젝트를 추가하는 함수입니다.
func addProject(client *mongo.Client, p Project) error {
	if p.ID == "" {
		return errors.New("빈 문자열입니다. 프로젝트를 생성할 수 없습니다")
	}
	if p.ID == "user" {
		return errors.New("user 이름으로 프로젝트를 생성할 수 없습니다")
	}

	collection := client.Database(*flagDBName).Collection("project")

	filter := bson.M{"id": p.ID}
	num, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 프로젝트가 존재해서 프로젝트를 생성할 수 없습니다")
	}

	_, err = collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}

	return nil
}

func getProject(client *mongo.Client, id string) (Project, error) {
	if id == "" {
		return Project{}, errors.New("프로젝트 이름이 빈 문자열입니다")
	}
	collection := client.Database(*flagDBName).Collection("project")

	filter := bson.M{"id": id}

	var p Project
	err := collection.FindOne(context.Background(), filter).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return p, errors.New(id + " 프로젝트가 존재하지 않습니다.")
		}
		p := Project{}
		p.ID = id
		return p, err
	}
	return p, nil
}

// 전체 프로젝트 정보를 가지고오는 함수입니다.
func getProjects(client *mongo.Client) ([]Project, error) {
	collection := client.Database(*flagDBName).Collection("project")

	sort := bson.D{{"id", 1}}

	cur, err := collection.Find(context.Background(), bson.M{}, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var results []Project
	if err := cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// getStatusProjects 함수는 상태를 받아서 프로젝트 정보를 가지고오는 함수입니다.
func getStatusProjects(client *mongo.Client, status ProjectStatus) ([]Project, error) {
	collection := client.Database(*flagDBName).Collection("project")

	filter := bson.M{"status": status}
	sort := bson.D{{"id", 1}}

	cur, err := collection.Find(context.Background(), filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var results []Project
	if err := cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func rmProject(client *mongo.Client, project string) error {
	collection := client.Database(*flagDBName).Collection("project")

	filter := bson.M{"id": project}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

// setProject 함수는 프로젝트 정보를 수정하는 함수입니다.
func setProject(client *mongo.Client, p Project) error {
	collection := client.Database(*flagDBName).Collection("project")

	p.Updatetime = time.Now().Format(time.RFC3339)

	filter := bson.M{"id": p.ID}
	update := bson.M{"$set": p}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// HasProject 함수는 프로젝트가 존재하는지 체크한다.
func HasProject(client *mongo.Client, project string) error {
	collection := client.Database(*flagDBName).Collection("project")

	filter := bson.M{"id": project}
	count, err := collection.CountDocuments(context.Background(), filter, options.Count().SetLimit(1))
	if err != nil {
		return err
	}

	if count != 1 {
		return errors.New(project + " 프로젝트가 존재하지 않습니다.")
	}

	return nil
}
