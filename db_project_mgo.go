package main

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Projectlist 함수는 프로젝트 리스트를 출력하는 함수입니다.
func Projectlist(session *mgo.Session) ([]string, error) {
	var results []string
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("project")
	projects := []Project{}
	err := c.Find(bson.M{}).Sort("id").All(&projects)
	if err != nil {
		return nil, err
	}
	for _, p := range projects {
		results = append(results, p.ID)
	}
	return results, nil
}

// 프로젝트를 추가하는 함수입니다.
func addProject(session *mgo.Session, p Project) error {
	if p.ID == "" {
		return errors.New("빈 문자열입니다. 프로젝트를 생성할 수 없습니다")
	}
	if p.ID == "user" {
		return errors.New("user 이름으로 프로젝트를 생성할 수 없습니다")
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("project")
	num, err := c.Find(bson.M{"id": p.ID}).Count()
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 프로젝트가 존재해서 프로젝트를 생성할 수 없습니다")
	}
	err = c.Insert(p)
	if err != nil {
		return err
	}
	return nil
}

func getProject(session *mgo.Session, id string) (Project, error) {
	if id == "" {
		return Project{}, errors.New("프로젝트 이름이 빈 문자열 입니다")
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("project")
	var p Project
	err := c.Find(bson.M{"id": id}).One(&p)
	if err != nil {
		if err == mgo.ErrNotFound {
			return p, errors.New(id + " 프로젝트가 존재하지 않습니다.")
		}
		p := Project{}
		p.ID = id
		return p, err
	}
	return p, nil
}

// 전체 프로젝트 정보를 가지고오는 함수입니다.
func getProjects(session *mgo.Session) ([]Project, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("project")
	results := []Project{}
	err := c.Find(bson.M{}).Sort("id").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

// getStatusProjects 함수는 상태를 받아서 프로젝트 정보를 가지고오는 함수입니다.
func getStatusProjects(session *mgo.Session, status ProjectStatus) ([]Project, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("project")
	results := []Project{}
	err := c.Find(bson.M{"status": status}).Sort("id").All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func rmProject(session *mgo.Session, project string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("project")
	err := c.Remove(bson.M{"id": project})
	if err != nil {
		return err
	}
	return nil
}

// setProject 함수는 프로젝트 정보를 수정하는 함수입니다.
func setProject(session *mgo.Session, p Project) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("project")
	p.Updatetime = time.Now().Format(time.RFC3339)
	err := c.Update(bson.M{"id": p.ID}, p)
	if err != nil {
		return err
	}
	return nil
}

// HasProject 함수는 프로젝트가 존재하는지 체크한다.
func HasProject(session *mgo.Session, project string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("project")
	num, err := c.Find(bson.M{"id": project}).Count()
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New(project + " 프로젝트가 존재하지 않습니다.")
	}
	return nil
}
