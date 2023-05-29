package main

import (
	"context"
	"errors"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2"
)

// SetAdminSetting 함수는 adminsetting을 DB에 저장한다.
func SetAdminSetting(client *mongo.Client, s Setting) error {
	collection := client.Database(*flagDBName).Collection("admin")

	filter := bson.D{{Key: "id", Value: "admin"}}
	opts := options.Update().SetUpsert(true)

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

// GetAdminSetting 함수는 adminsetting을 DB에서 가지고 온다.
func GetAdminSetting(client *mongo.Client) (Setting, error) {
	collection := client.Database(*flagDBName).Collection("admin")

	filter := bson.D{{Key: "id", Value: "admin"}}
	opts := options.FindOne()

	var s Setting
	err := collection.FindOne(context.TODO(), filter, opts).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.ID = "admin"
			_, err = collection.InsertOne(context.TODO(), s)
			if err != nil {
				return s, err
			}
			return s, nil
		}
		return s, err
	}

	return s, nil
}

// AddTaskSetting 함수는 tasksetting을 DB에 추가한다.
func AddTaskSetting(session *mgo.Session, t Tasksetting) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	num, err := c.Find(bson.M{"id": t.ID}).Count()
	if err != nil {
		return err
	}
	if num > 0 {
		return errors.New("이미 Tasksetting이 존재합니다")
	}
	err = c.Insert(t)
	if err != nil {
		return err
	}
	return nil
}

// RmTaskSetting 함수는 tasksetting을 DB에 추가한다.
func RmTaskSetting(session *mgo.Session, name, typ string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	err := c.Remove(bson.M{"name": name, "type": typ})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskSetting 함수는 Tasksetting 값을 바꾼다.
func SetTaskSetting(session *mgo.Session, t Tasksetting) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	err := c.Update(bson.M{"id": t.ID}, t)
	if err != nil {
		return err
	}
	return nil
}

// AllTaskSettings 함수는 모든 tasksetting값을 가지고 온다.
func AllTaskSettings(session *mgo.Session) ([]Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	results := []Tasksetting{}
	err := c.Find(bson.M{}).Sort("order").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// getTaskSetting 함수는 id를 입력받아서 tasksetting값을 가지고 온다.
func getTaskSetting(session *mgo.Session, id string) (Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	result := Tasksetting{}
	err := c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// getShotTaskSetting 함수는 type이 shot인 tasksetting값을 가지고 온다.
func getShotTaskSetting(session *mgo.Session) ([]Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	results := []Tasksetting{}
	err := c.Find(bson.M{"type": "shot"}).Sort("name").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// getAssetTaskSetting 함수는 type이 asset인 tasksetting값을 가지고 온다.
func getAssetTaskSetting(session *mgo.Session) ([]Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	results := []Tasksetting{}
	err := c.Find(bson.M{"type": "asset"}).Sort("name").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// getCategoryTaskSettings 함수는 type이 asset인 tasksetting값을 가지고 온다.
func getCategoryTaskSettings(session *mgo.Session, category string) ([]Tasksetting, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	results := []Tasksetting{}
	err := c.Find(bson.M{"category": category}).Sort("name").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// TasksettingNames 함수는 Tasksetting 이름을 수집하여 반환한다.
func TasksettingNames(session *mgo.Session) ([]string, error) {
	var results []string
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	err := c.Find(bson.M{}).Distinct("name", &results)
	if err != nil {
		return nil, err
	}
	sort.Strings(results)
	return results, nil
}

// Unique 함수는 리스트에서 중복되는 문자열을 제거한다.
func Unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// TasksettingNamesByExcelOrder 함수는 Tasksetting 이름을 ExcelOrder순으로 반환한다.
func TasksettingNamesByExcelOrder(session *mgo.Session) ([]string, error) {
	var tasksettings []Tasksetting
	var results []string
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("tasksetting")
	err := c.Find(bson.M{}).Sort("excelorder").All(&tasksettings)
	if err != nil {
		return nil, err
	}
	for _, t := range tasksettings {
		results = append(results, t.Name)
	}
	return Unique(results), nil
}
