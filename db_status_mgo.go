package main

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
)

// SetStatus 함수는 Status를 DB에 저장한다.
func SetStatus(client *mongo.Client, s Status) error {
	collection := client.Database(*flagDBName).Collection("status")

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: s.ID}}
	update := bson.D{{Key: "$set", Value: s}}

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = collection.InsertOne(context.TODO(), s)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// GetStatus 함수는 Status를 DB에서 가지고 온다.
func GetStatus(session *mgo.Session, id string) (Status, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("status")
	s := Status{}
	err := c.Find(bson.M{"id": id}).One(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// GetInitStatusID 함수는 초기 Status를 DB에서 가지고 온다.
func GetInitStatusID(session *mgo.Session) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("status")
	s := Status{}
	n, err := c.Find(bson.M{"initstatus": true}).Count()
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", errors.New("아이템 생성시 초기 Status로 사용할 Status 설정이 필요합니다")
	}
	if n != 1 {
		return "", errors.New("초기 상태 설정값이 1개가 아닙니다")
	}
	err = c.Find(bson.M{"initstatus": true}).One(&s)
	if err != nil {
		return "", err
	}
	return s.ID, nil
}

// AddStatus 함수는 tasksetting을 DB에 추가한다.
func AddStatus(session *mgo.Session, s Status) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("status")
	n, err := c.Find(bson.M{"id": s.ID}).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New(s.ID + " Status가 이미 존재합니다")
	}
	err = c.Insert(s)
	if err != nil {
		return err
	}
	return nil
}

// RmStatus 함수는 Status를 DB에서 삭제한다.
func RmStatus(session *mgo.Session, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("status")
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

// AllStatus 함수는 모든 Status값을 DB에서 가지고 온다.
func AllStatus(client *mongo.Client) ([]Status, error) {
	collection := client.Database(*flagDBName).Collection("status")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "order", Value: -1}})

	cur, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var results []Status
	for cur.Next(context.Background()) {
		var status Status
		err := cur.Decode(&status)
		if err != nil {
			log.Println(err)
			continue
		}

		results = append(results, status)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return results, nil
}
