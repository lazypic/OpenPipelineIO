package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

// validTokenV2 함수는 Token이 유효한지 체크한다.
func validTokenV2(client *mongo.Client, token string) (Token, error) {
	collection := client.Database(*flagDBName).Collection("token")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := Token{}
	err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&t)
	if err != nil {
		return t, errors.New("authorization failed")
	}
	return t, nil
}

// vaildUserV2 함수는 사용자의 id, pw를 받아서 유효한 사용자인지 체크한다.
func vaildUserV2(client *mongo.Client, id, pw string) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}
	u := User{}
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&u)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if err != nil {
		return err
	}
	return nil
}

// addPasswordAttemptV2 함수는 사용자의 id를 받아서 패스워드 시도횟수를 추가한다.
func addPasswordAttemptV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$inc": bson.M{"passwordattempt": 1}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id" + id)
	}
	return nil
}

// setUserV2 함수는 사용자 정보를 업데이트하는 함수이다.
func setUserV2(client *mongo.Client, u User) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u.Updatetime = time.Now().Format(time.RFC3339)
	filter := bson.M{"id": u.ID}
	update := bson.D{{Key: "$set", Value: u}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id" + u.ID)
	}
	return nil
}

func rmTokenV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("token")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := collection.FindOneAndDelete(ctx, bson.M{"id": id})
	result.Err()
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return mongo.ErrNoDocuments
		}
		return result.Err()
	}
	return nil
}

// allUsersV2 함수는 전체 프로젝트 정보를 가지고오는 함수입니다.
func allUsersV2(client *mongo.Client) ([]User, error) {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	results := []User{}
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
