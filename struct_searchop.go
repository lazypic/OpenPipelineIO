package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

// SearchOption 은 웹 검색창의 옵션 자료구조이다.
type SearchOption struct {
	Project    string   `json:"project"`    // 선택한 프로젝트
	Searchword string   `json:"searchword"` // 검색어
	Sortkey    string   `json:"sortkey"`    // 정렬방식
	Task       string   `json:"task"`       // Task명
	TrueStatus []string `json:"truestatus"` // true 상태리스트
	Shot       bool     `json:"shot"`
	Assets     bool     `json:"assets"`
	Type3d     bool     `json:"type3d"`
	Type2d     bool     `json:"type2d"`
	Page       int      `json:"page"`
}

// SearchOption과 관련된 메소드

func (op *SearchOption) setStatusAll() error {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return err
	}
	defer session.Close()
	status, err := AllStatus(session)
	if err != nil {
		return err
	}
	for _, s := range status {
		op.TrueStatus = append(op.TrueStatus, s.ID)
	}
	return nil
}

func (op *SearchOption) setStatusDefaultV1() error {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return err
	}
	defer session.Close()
	status, err := AllStatus(session)
	if err != nil {
		return err
	}
	for _, s := range status {
		if !s.DefaultOn {
			continue
		}
		op.TrueStatus = append(op.TrueStatus, s.ID)
	}
	return nil
}

func (op *SearchOption) setStatusDefaultV2() error {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		return err
	}
	defer session.Close()
	status, err := AllStatus(session)
	if err != nil {
		return err
	}
	for _, s := range status {
		if !s.DefaultOn {
			continue
		}
		op.TrueStatus = append(op.TrueStatus, s.ID)
	}
	return nil
}

func (op *SearchOption) setStatusNone() {
	op.TrueStatus = []string{}
}

func handleRequestToSearchOption(r *http.Request) SearchOption {
	q := r.URL.Query()
	op := SearchOption{}
	op.Project = q.Get("project")
	op.Searchword = q.Get("searchword")
	op.Sortkey = q.Get("sortkey")
	op.Task = q.Get("task")
	// 페이지를 구한다.
	page, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		op.Page = 1 // 에러가 발생하면 1페이지로 이동한다.
	} else {
		op.Page = page
	}
	for _, s := range strings.Split(q.Get("truestatus"), ",") {
		if s == "" {
			continue
		}
		op.TrueStatus = append(op.TrueStatus, s)
	}
	return op
}

// LoadCookie 메소드는 request에 이미 설정된 쿠키값을을 SearchOption 자료구조에 추가한다.
func (op *SearchOption) LoadCookie(session *mgo.Session, r *http.Request) error {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "Project" {
			op.Project = cookie.Value
		}
		if cookie.Name == "Task" {
			op.Task = cookie.Value
		}
		if cookie.Name == "Searchword" {
			cookieByte, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err != nil {
				log.Println(err)
			}
			op.Searchword = string(cookieByte)
		}
	}
	if op.Project == "" {
		plist, err := Projectlist(session)
		if err != nil {
			return err
		}
		op.Project = plist[0] // 프로젝트가 빈 문자열이면 첫번째 프로젝트를 설정합니다.
	}
	return nil
}

// LoadCookie 메소드는 request에 이미 설정된 쿠키값을을 SearchOption 자료구조에 추가한다.
func (op *SearchOption) LoadCookieV2(client *mongo.Client, r *http.Request) error {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "Project" {
			op.Project = cookie.Value
		}
		if cookie.Name == "Task" {
			op.Task = cookie.Value
		}
		if cookie.Name == "Searchword" {
			cookieByte, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err != nil {
				log.Println(err)
			}
			op.Searchword = string(cookieByte)
		}
	}
	if op.Project == "" {
		plist, err := ProjectlistV2(client)
		if err != nil {
			return err
		}
		op.Project = plist[0] // 프로젝트가 빈 문자열이면 첫번째 프로젝트를 설정합니다.
	}
	return nil
}
