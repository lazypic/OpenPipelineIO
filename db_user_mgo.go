// 이 코드는 사용자와 관련된 DBapi가 모여있는 파일입니다.

package main

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// legacy
// validToken 함수는 Token이 유효한지 체크한다.
func validToken(session *mgo.Session, token string) (Token, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("token")
	t := Token{}
	err := c.Find(bson.M{"token": token}).One(&t)
	if err != nil {
		return t, errors.New("authorization failed")
	}
	return t, nil
}

// getUser 함수는 사용자를 가지고오는 함수이다.
func getUser(session *mgo.Session, id string) (User, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("users")
	u := User{}
	err := c.Find(bson.M{"id": id}).One(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}
