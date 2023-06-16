package main

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetImageSize 함수는 해당 샷의 이미지 사이즈를 설정한다. // legacy
// key 설정값 : platesize, undistortionsize, rendersize
func SetImageSize(session *mgo.Session, project, id, key, size string) error {
	if !(key == "platesize" || key == "undistortionsize" || key == "rendersize") {
		return errors.New("잘못된 key값입니다")
	}
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	if key == "undistortionsize" {
		err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"undistortionsize": size, "updatetime": time.Now().Format(time.RFC3339)}})
		if err != nil {
			return err
		}
	} else {
		err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{key: size, "updatetime": time.Now().Format(time.RFC3339)}})
		if err != nil {
			return err
		}
	}
	return nil
}

// SetTaskUser 함수는 item에 task의 user 값을 셋팅한다.
func SetTaskUser(session *mgo.Session, project, name, task, user string) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", err
	}
	id := name + "_" + typ
	err = HasTask(session, project, id, task)
	if err != nil {
		return id, err
	}
	item, err := getItem(session, project, id)
	if err != nil {
		return id, err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": item.ID}, bson.M{"$set": bson.M{"tasks." + task + ".user": user, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}
