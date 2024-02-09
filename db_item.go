package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
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

func AddTagV2(client *mongo.Client, id, inputTag string) (string, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	i, err := getItemV2(client, id)
	if err != nil {
		return id, err
	}
	rmspaceTag := strings.Replace(inputTag, " ", "", -1) // 태그는 공백을 제거한다.
	for _, tag := range i.Tag {
		if rmspaceTag == tag {
			return id, errors.New(inputTag + "태그는 이미 존재하고 있습니다 추가할 수 없습니다")
		}
	}
	newTags := append(i.Tag, rmspaceTag)

	result, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"tag": newTags, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return i.Name, err
	}
	if result.MatchedCount == 0 {
		return i.Name, errors.New("no document found with id" + id)
	}
	return i.Name, nil
}

func getItemV2(client *mongo.Client, id string) (Item, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := Item{}
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func rmItemIDV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

func RmTag(client *mongo.Client, id, inputTag string, isContain bool) (string, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	i, err := getItemV2(client, id)
	if err != nil {
		return "", err
	}
	var newTags []string
	for _, tag := range i.Tag {
		if isContain {
			if strings.Contains(tag, inputTag) {
				continue
			}
		}
		if inputTag == tag {
			continue
		}
		newTags = append(newTags, tag)
	}
	i.Tag = newTags
	// 만약 태그에 권정보가 없더라도 권관련 태그는 날아가면 안된다. setItem을 이용한다.

	filter := bson.M{"id": id}
	update := bson.D{{Key: "$set", Value: i}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return i.ID, err
	}
	if result.MatchedCount == 0 {
		return i.ID, errors.New("no document found with id" + i.ID)
	}

	return i.ID, nil
}

func AddTaskV2(client *mongo.Client, id, task, status string) error {
	item, err := getItemV2(client, id)
	if err != nil {
		return err
	}

	taskname := strings.ToLower(task)
	// 기존에 Task가 없다면 추가한다.
	if _, found := item.Tasks[task]; !found {
		t := Task{}
		t.Title = taskname
		t.StatusV2 = status
		item.Tasks[task] = t
	} else {
		return fmt.Errorf("이미 %s 에 %s Task가 존재합니다", id, taskname)
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item.Updatetime = time.Now().Format(time.RFC3339)
	globalStatus, err := AllStatusV2(client)
	if err != nil {
		return err
	}
	item.updateStatusV2(globalStatus)

	filter := bson.M{"id": item.ID}
	update := bson.D{{Key: "$set", Value: item}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + item.ID)
	}
	return nil

}

func RmTaskV2(client *mongo.Client, id, taskname string) error {
	item, err := getItemV2(client, id)
	if err != nil {
		return err
	}

	delete(item.Tasks, taskname)

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item.Updatetime = time.Now().Format(time.RFC3339)

	status, err := AllStatusV2(client)
	if err != nil {
		return err
	}
	item.updateStatusV2(status)

	filter := bson.M{"id": item.ID}
	update := bson.D{{Key: "$set", Value: item}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + item.ID)
	}
	return nil
}
