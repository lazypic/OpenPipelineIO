package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// getUser 함수는 사용자를 가지고오는 함수이다.
func getUserV2(client *mongo.Client, id string) (User, error) {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := User{}
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// addUserV2 함수는 사용자를 추가하는 함수이다.
func addUserV2(client *mongo.Client, u User) error {
	if u.ID == "" {
		err := errors.New("ID가 빈 문자열입니다. 유저를 생성할 수 없습니다")
		return err
	}
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"id": u.ID})
	if err != nil {
		return err
	}

	if num != 0 {
		err = errors.New(u.ID + " ID를 가진 사용자가 이미 DB에 존재합니다.")
		return err
	}
	u.Createtime = time.Now().Format(time.RFC3339)
	_, err = collection.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

// addTokenV2 함수는 사용자정보로 token을 추가하는 함수이다.
func addTokenV2(client *mongo.Client, u User) error {
	collection := client.Database(*flagDBName).Collection("token")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	num, err := collection.CountDocuments(ctx, bson.M{"token": u.Token})
	if err != nil {
		return err
	}
	if num != 0 {
		err = errors.New(u.Token + " 키가 이미 DB에 존재합니다.")
		return err
	}
	t := Token{
		Token:       u.Token,
		AccessLevel: u.AccessLevel,
		ID:          u.ID,
	}
	_, err = collection.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	return nil
}
