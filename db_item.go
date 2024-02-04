package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func addItemV2(client *mongo.Client, i Item) error {
	err := i.CheckError()
	if err != nil {
		return err
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"id": i.ID})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 이름을 가진 데이터가 있습니다")
	}
	_, err = collection.InsertOne(ctx, i)
	if err != nil {
		return err
	}
	return nil
}

// DistinctDdline 함수는 프로젝트, dict key를 받아서 key에 사용되는 모든 마감일을 반환한다. 예) 태그
func DistinctDdlineV2(client *mongo.Client, project string, key string) ([]string, error) {
	var results []string
	if project == "" || key == "" {
		return results, nil
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "project", Value: project}}
	values, err := collection.Distinct(ctx, key, filter)
	if err != nil {
		return results, err
	}
	for _, value := range values {
		results = append(results, fmt.Sprintf("%v", value))
	}
	sort.Strings(results)

	if *flagDebug {
		fmt.Println("DB에서 가지고온 마감일 리스트")
		fmt.Println(results)
		fmt.Println()
	}
	var before string
	var datelist []string
	for _, r := range results {
		if r != "" {
			date := ToNormalTime(r)
			if date == before {
				break
			} else {
				datelist = append(datelist, date)
			}
			before = date
		}
	}
	sort.Strings(datelist) //기존 OpenPipelineIO 2.0의 4자리 수를 위하여 정렬한다. 추후 이 줄은 사라진다.
	if *flagDebug {
		fmt.Println("마감일을 Tag형태로 바꾼 리스트")
		fmt.Println(datelist)
		fmt.Println()
	}
	return datelist, nil
}

func DistinctV2(client *mongo.Client, project string, key string) ([]string, error) {
	var results []string
	if project == "" || key == "" {
		return results, nil
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "project", Value: project}}
	values, err := collection.Distinct(ctx, key, filter)
	if err != nil {
		return results, err
	}
	for _, value := range values {
		results = append(results, fmt.Sprintf("%v", value))
	}
	sort.Strings(results)
	return results, nil
}

func AllAssetsV2(client *mongo.Client, project string) ([]string, error) {
	var results []string
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{
		"project": project,
		"type":    "asset",
	}

	values, err := collection.Distinct(ctx, "name", filter)
	if err != nil {
		return results, err
	}
	for _, value := range values {
		if name, ok := value.(string); ok {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results, nil
}

func TotalnumV2(client *mongo.Client, project string) (Infobarnum, error) {
	if project == "" {
		return Infobarnum{}, nil
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var results Infobarnum

	filter := bson.M{"$or": []bson.M{{"project": project, "type": "org"}, {"project": project, "type": "left"}}}
	num, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return Infobarnum{}, err
	}
	results.Total = int(num)
	return results, nil
}
