package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// searchUsersV2 함수는 검색을 입력받고 해당 검색어가 있는 사용자 정보를 가지고 옵니다.
func searchUsersV2(client *mongo.Client, words []string) ([]User, error) {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var searchwords []string
	// 사람 이름을 가지고 검색을 자주한다.
	for _, word := range words {
		if isASCII(word) {
			searchwords = append(searchwords, word)
			continue
		}
		if strings.HasPrefix(word, "tag:") {
			searchwords = append(searchwords, word)
			continue
		}
		r := []rune(word)
		if len(r) == 2 { // 이름이 2자리 일경우 "김웅"
			searchwords = append(searchwords, string(r[0])) // 성을 추가한다.
			searchwords = append(searchwords, string(r[1])) // 이름을 추가한다.
			continue
		} else if len(r) == 3 { // 이름일 확률이 높다.
			searchwords = append(searchwords, string(r[0]))  // 성을 추가한다.
			searchwords = append(searchwords, string(r[1:])) // 이름을 추가한다.
			continue
		} else if len(r) == 4 { // 이름이 4자리 일경우가 있다. 예)독고영재
			searchwords = append(searchwords, string(r[2:])) // 이름이 고영재" 또는 "영재" 일 수 있다. 이름을 위주로 검색시킨다.
			continue
		}
		searchwords = append(searchwords, word)
	}

	allQueries := bson.A{}
	if *flagDebug {
		log.Println(searchwords)
	}
	for _, word := range searchwords {
		orQueries := bson.A{}
		if strings.HasPrefix(word, "tag:") {
			orQueries = append(orQueries, bson.M{"tags": strings.TrimPrefix(word, "tag:")})
		} else if strings.HasPrefix(word, "id:") {
			orQueries = append(orQueries, bson.M{"id": strings.TrimPrefix(word, "id:")})
		} else {
			orQueries = append(orQueries, bson.M{"id": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"employeenumber": primitive.Regex{Pattern: word, Options: "i"}})
			orQueries = append(orQueries, bson.M{"firstnamekor": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"lastnamekor": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"firstnameeng": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"lastnameeng": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"firstnamechn": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"lastnamechn": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"email": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"emailexternal": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"phone": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"hotline": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"location": primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"tags": primitive.Regex{Pattern: word, Options: "i"}})
			orQueries = append(orQueries, bson.M{"lastip": primitive.Regex{Pattern: word}})
		}
		allQueries = append(allQueries, bson.M{"$or": orQueries})
	}

	q := bson.M{"$and": allQueries}

	var results []User

	cursor, err := collection.Find(ctx, q)
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func UserTagsV2(client *mongo.Client) ([]string, error) {
	var tags []string

	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	values, err := collection.Distinct(ctx, "tags", filter)
	if err != nil {
		return tags, err
	}
	for _, value := range values {
		tags = append(tags, fmt.Sprintf("%v", value))
	}
	sort.Strings(tags)
	return tags, nil
}

func getTokenV2(client *mongo.Client, id string) (Token, error) {
	collection := client.Database(*flagDBName).Collection("token")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t := Token{}
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func setTokenV2(client *mongo.Client, t Token) error {
	collection := client.Database(*flagDBName).Collection("token")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": t.ID}
	update := bson.D{{Key: "$set", Value: t}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id" + t.ID)
	}
	return nil
}

func updatePasswordUserV2(client *mongo.Client, id, pw, newPw string) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	// 과거의 패스워드로 로그인가능했는지 체크한다.
	err := vaildUserV2(client, id, pw)
	if err != nil {
		return err
	}
	// 새로운 패스워드로 업데이트 한다.
	encryptPass, err := Encrypt(newPw)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"password":   encryptPass,
			"updatetime": time.Now().Format(time.RFC3339),
			"token":      base64.StdEncoding.EncodeToString([]byte(encryptPass)),
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func initPassUser(client *mongo.Client, id, initpassword string) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}

	num, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New("해당 유저가 존재하지 않습니다")
	}

	encryptPass, err := Encrypt(initpassword)
	if err != nil {
		log.Println(err)
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"password":        encryptPass,
			"passwordattempt": 0,
			"updatetime":      time.Now().Format(time.RFC3339),
			"token":           base64.StdEncoding.EncodeToString([]byte(encryptPass)),
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}

	return nil
}

func rmUserV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func ReplaceTagsV2(client *mongo.Client, old, new string) error {
	collection := client.Database(*flagDBName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []User
	filter := bson.M{"tags": primitive.Regex{Pattern: old}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		return err
	}

	for _, u := range users {
		var newTags []string
		for _, tag := range u.Tags {
			if tag == old {
				newTags = append(newTags, new)
				continue
			}
			newTags = append(newTags, tag)
		}
		u.Tags = newTags
		update := bson.D{{Key: "$set", Value: u}}
		_, err = collection.UpdateOne(ctx, bson.M{"id": u.ID}, update)
		if err != nil {
			return err
		}

	}
	// 각 유저를 체크하면서 태그이름을 변경한다.
	return nil
}
