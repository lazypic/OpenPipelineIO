package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// SetImageSize 함수는 해당 샷의 이미지 사이즈를 설정한다. // legacy
// key 설정값 : platesize, undistortionsize, rendersize
func SetImageSize(client *mongo.Client, project, name, key, size string) (string, error) {
	if !(key == "platesize" || key == "dsize" || key == "undistortionsize" || key == "rendersize") {
		return "", errors.New("잘못된 key값입니다")
	}

	err := HasProject(client, project)
	if err != nil {
		return "", err
	}

	typ, err := Type(client, project, name)
	if err != nil {
		return "", err
	}

	id := name + "_" + typ
	collection := client.Database(*flagDBName).Collection("items")

	update := bson.M{"updatetime": time.Now().Format(time.RFC3339)}
	if key == "dsize" || key == "undistortionsize" {
		update["dsize"] = size
		update["undistortionsize"] = size
	} else {
		update[key] = size
	}

	filter := bson.M{"id": id}
	updateStmt := bson.M{"$set": update}

	_, err = collection.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return id, err
	}

	return id, nil
}

// SetTaskUser 함수는 item에 task의 user 값을 셋팅한다.
func SetTaskUser(client *mongo.Client, project, name, task, user string) (string, error) {
	err := HasProject(client, project)
	if err != nil {
		return "", err
	}

	typ, err := Type(client, project, name)
	if err != nil {
		return "", err
	}

	id := name + "_" + typ

	err = HasTask(client, project, id, task)
	if err != nil {
		return id, err
	}

	item, err := getItem(client, project, id)
	if err != nil {
		return id, err
	}

	collection := client.Database(*flagDBName).Collection("items")

	update := bson.M{
		"tasks." + task + ".user": user,
		"updatetime":              time.Now().Format(time.RFC3339),
	}
	filter := bson.M{"id": item.ID}
	updateStmt := bson.M{"$set": update}

	_, err = collection.UpdateOne(context.Background(), filter, updateStmt)
	if err != nil {
		return id, err
	}

	return id, nil
}
