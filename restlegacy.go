package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"gopkg.in/mgo.v2"
)

// handleAPISetThummov 함수는 아이템의 Thummov 값을 설정한다.
func handleAPISetThummov(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project  string `json:"project"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		UserID   string `json:"userid"`
		Error    string `json:"error"`
		Protocol string `json:"protocol"`
	}
	rcp := Recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	for key, values := range r.PostForm {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			rcp.Project = v
		case "name":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Name = v
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		case "path":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Path = v
		}
	}
	err = SetThummov(session, rcp.Project, rcp.Name, rcp.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set Thumbnail: %s\nProject: %s, Name: %s, Author: %s", rcp.Path, rcp.Project, rcp.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIRenderSize 함수는 아이템에 RenderSize를 설정한다.
func handleAPISetRenderSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		ID      string `json:"id"`
		Size    string `json:"size"`
		UserID  string `json:"userid"`
		Error   string `json:"error"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	for key, values := range r.PostForm {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Project = v
		case "id":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.ID = v
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		case "size":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if !regexpImageSize.MatchString(v) {
				http.Error(w, "2048x1152 형태로 입력해주세요", http.StatusBadRequest)
				return
			}
			rcp.Size = v
		}
	}
	err = SetImageSize(session, rcp.Project, rcp.ID, "rendersize", rcp.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set Rendersize: %s\nID: %s, Author: %s", rcp.Size, rcp.ID, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetRnum 함수는 아이템에 롤넘버를 설정합니다.
func handleAPISetRnum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		ID      string `json:"id"`
		Rnum    string `json:"rnum"`
		UserID  string `json:"userid"`
		Error   string `json:"error"`
	}
	rcp := Recipe{}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	rcp.UserID, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	for key, values := range r.PostForm {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Project = v
		case "name":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Name = v
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		case "rnum":
			if len(values) == 0 {
				rcp.Rnum = ""
			} else {
				rcp.Rnum = values[0]
			}
		}
	}
	if rcp.Rnum != "" && !regexpRnum.MatchString(rcp.Rnum) {
		http.Error(w, rcp.Rnum+"값은 A0001 형식이 아닙니다.", http.StatusBadRequest)
		return
	}
	typ, err := Type(session, rcp.Project, rcp.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = SetRnum(session, rcp.Project, rcp.Name+"_"+typ, rcp.Rnum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.ID = rcp.Name + "_" + typ
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set Rnum: %s\nProject: %s, Name: %s, Author: %s", rcp.Rnum, rcp.Project, rcp.Name, rcp.UserID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, _ := json.Marshal(rcp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
