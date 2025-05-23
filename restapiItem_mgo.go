package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

// handleAPI2Item 함수는 아이템 자료구조를 불러온다.
func handleAPI2GetItem(w http.ResponseWriter, r *http.Request) {
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	if id == "" {
		http.Error(w, errors.New("need id").Error(), http.StatusInternalServerError)
		return
	}
	item, err := getItemV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPITimeinfo(w http.ResponseWriter, r *http.Request) {
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}

	item, err := getItemV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipe struct {
		ScanIn          int    `json:"scanin"`
		ScanOut         int    `json:"scanout"`
		ScanFrame       int    `json:"scanframe"`
		ScanTimecodeIn  string `json:"scantimecodein"`
		ScanTimecodeOut string `json:"scantimecodeout"`
		PlateIn         int    `json:"platein"`
		PlateOut        int    `json:"plateout"`
		HandleIn        int    `json:"handlein"`
		HandleOut       int    `json:"handleout"`
		JustIn          int    `json:"justin"`
		JustOut         int    `json:"justout"`
		JustTimecodeIn  string `json:"justtimecodein"`
		JustTimecodeOut string `json:"justtimecodeout"`
	}
	rcp := recipe{}
	rcp.ScanIn = item.ScanIn
	rcp.ScanOut = item.ScanOut
	rcp.ScanFrame = item.ScanFrame
	rcp.ScanTimecodeIn = item.ScanTimecodeIn
	rcp.ScanTimecodeOut = item.ScanTimecodeOut
	rcp.PlateIn = item.PlateIn
	rcp.PlateOut = item.PlateOut
	rcp.HandleIn = item.HandleIn
	rcp.HandleOut = item.HandleOut
	rcp.JustIn = item.JustIn
	rcp.JustOut = item.JustOut
	rcp.JustTimecodeIn = item.JustTimecodeIn
	rcp.JustTimecodeOut = item.JustTimecodeOut

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

// handleAPIRmItemID 함수는 아이템을 삭제한다.
func handleAPIRmItemID(w http.ResponseWriter, r *http.Request) {
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	userID, level, err := TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if PmAccessLevel > level {
		http.Error(w, "need permission", http.StatusUnauthorized)
		return
	}
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	type Recipe struct {
		ID     string `json:"id"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	rcp.UserID = userID
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	err = rmItemIDV2(client, id)
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

// handleAPI3Items 함수는 아이템을 검색한다.
func handleAPI3Items(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}

	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q, err := URLUnescape(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if q.Get("project") == "" {
		http.Error(w, "The project string is empty", http.StatusInternalServerError)
		return
	}
	op := SearchOption{
		Project:    q.Get("project"),
		Searchword: q.Get("searchword"),
		Sortkey:    "id",
		Type3d:     str2bool(q.Get("type3d")),
		Type2d:     str2bool(q.Get("type2d")),
		TrueStatus: strings.Split(q.Get("truestatus"), ","),
	}
	if q.Get("sortkey") != "" {
		op.Sortkey = "id"
	}
	items, err := SearchV2(client, op)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIShot 함수는 project, name을 받아서 shot을 반환한다.
func handleAPIGetShot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	name := q.Get("name")
	type recipe struct {
		Data Item `json:"data"`
	}
	rcp := recipe{}
	result, err := Shot(session, project, name)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	rcp.Data = result
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIAsset 함수는 project, name을 받아서 asset을 반환한다.
func handleAPIGetAsset(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	if project == "" {
		http.Error(w, "project 가 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	name := q.Get("name")
	if name == "" {
		http.Error(w, "name 이 빈 문자열 입니다", http.StatusBadRequest)
		return
	}
	item, err := Asset(session, project, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISeqs 함수는 프로젝트의 시퀀스를 가져온다.
func handleAPISeqs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	if project == "" {
		fmt.Fprintln(w, "{\"error\":\"project 정보가 없습니다.\"}")
		return
	}
	seqs, err := Seqs(session, project)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	type recipe struct {
		Data []string `json:"data"`
	}
	rcp := recipe{}
	rcp.Data = seqs
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIAssets 함수는 project 정보를 입력받아서 모든 에셋 정보를 출력한다.
func handleAPIAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	if project == "" {
		http.Error(w, "project 가 빈 문자열입니다", http.StatusBadRequest)
		return
	}
	assets, err := SearchAllAsset(session, project, "name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(assets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPI2Shots 함수는 project, seq를 입력받아서 cut 이름을 출력합니다.
func handleAPI2Shots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	if project == "" {
		http.Error(w, "project 정보가 없습니다", http.StatusBadRequest)
		return
	}
	seq := q.Get("seq")
	if seq == "" {
		http.Error(w, "seq 정보가 없습니다", http.StatusBadRequest)
		return
	}
	shots, err := Shots(session, project, seq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	data, err := json.Marshal(shots)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIAllShots 함수는 project를 받아서 모든 샷을 출력한다.
func handleAPIAllShots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	project := q.Get("project")
	if project == "" {
		http.Error(w, "project 정보가 없습니다.", http.StatusUnauthorized)
		return
	}
	items, err := SearchAllShot(session, project, "name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIUseTypes 함수는 project를 받아서 모든 샷을 출력한다.
func handleAPIUseTypes(w http.ResponseWriter, r *http.Request) {
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()

	id := q.Get("id")
	if id == "" {
		http.Error(w, "need id", http.StatusUnauthorized)
		return
	}
	useType, types, err := UseTypesV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Recipe struct {
		Types   []string `json:"types"`
		Usetype string   `json:"usetype"`
	}
	rcp := Recipe{}
	rcp.Types = types
	rcp.Usetype = useType
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPI2SetTaskMov 함수는 Task에 mov를 설정한다.
func handleAPI2SetTaskMov(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Task     string `json:"task"`
		Mov      string `json:"mov"`
		UserID   string `json:"userid"`
		Error    string `json:"error"`
		Protocol string `json:"protocol"`
	}
	rcp := Recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	rcp.Mov = r.FormValue("mov")
	err = setTaskMovV2(client, rcp.ID, rcp.Task, rcp.Mov)
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

// handleAPISetTaskExpectDay 함수는 Task에 예상일을 설정한다.
func handleAPISetTaskExpectDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project   string `json:"project"`
		ID        string `json:"id"`
		Task      string `json:"task"`
		ExpectDay int    `json:"expectday"`
		UserID    string `json:"userid"`
		Error     string `json:"error"`
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
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "task를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	expectday := r.FormValue("expectday")
	if expectday == "" {
		http.Error(w, "expectday를 설정해주세요", http.StatusBadRequest)
		return
	}
	num, err := strconv.Atoi(expectday)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.ExpectDay = num
	err = setTaskExpectDay(session, rcp.Project, rcp.ID, rcp.Task, rcp.ExpectDay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set ExpectDay: %s %d\nProject: %s, ID: %s, Author: %s", rcp.Task, rcp.ExpectDay, rcp.Project, rcp.ID, rcp.UserID))
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

// handleAPISetTaskUserComment 함수는 Task에 UserComment를 설정한다.
func handleAPISetTaskUserComment(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID          string `json:"id"`
		Task        string `json:"task"`
		UserComment string `json:"usercomment"`
		UserID      string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	rcp.UserComment = r.FormValue("usercomment")
	err = setTaskUserCommentV2(client, rcp.ID, rcp.Task, rcp.UserComment)
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

// handleAPISetTaskResultDay 함수는 Task에 실제 작업일을 설정한다.
func handleAPISetTaskResultDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project   string `json:"project"`
		ID        string `json:"id"`
		Task      string `json:"task"`
		ResultDay int    `json:"resultday"`
		UserID    string `json:"userid"`
		Error     string `json:"error"`
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
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "task를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	resultday := r.FormValue("resultday")
	if resultday == "" {
		http.Error(w, "resultday를 설정해주세요", http.StatusBadRequest)
		return
	}
	num, err := strconv.Atoi(resultday)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.ResultDay = num
	err = setTaskResultDay(session, rcp.Project, rcp.ID, rcp.Task, rcp.ResultDay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set ResultDay: %s %d\nProject: %s, ID: %s, Author: %s", rcp.Task, rcp.ResultDay, rcp.Project, rcp.ID, rcp.UserID))
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

// handleAPIUnDistortionSize 함수는 아이템의 DistortionSize를 설정한다.
func handleAPISetUnDistortionSize(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Size   string `json:"size"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	size := r.FormValue("size")
	if !regexpImageSize.MatchString(size) {
		http.Error(w, "Please enter in the format of 2048x1152", http.StatusBadRequest)
		return
	}
	rcp.Size = size

	err = SetImageSizeV2(client, rcp.ID, "undistortionsize", rcp.Size)
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

// handleAPISetJustIn 함수는 아이템에 JustIn 값을 설정한다.
func handleAPISetJustIn(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Frame  int    `json:"frame"`
		Error  string `json:"error"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "justin", rcp.Frame)
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

// handleAPISetPlateIn 함수는 아이템에 PlateIn 값을 설정한다.
func handleAPISetPlateIn(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Frame int    `json:"frame"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "platein", rcp.Frame)
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

// handleAPISetPlateOut 함수는 아이템에 PlateOut 값을 설정한다.
func handleAPISetPlateOut(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Frame int    `json:"frame"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "plateout", rcp.Frame)
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

// handleAPISetScanIn 함수는 아이템에 ScanIn 값을 설정한다.
func handleAPISetScanIn(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Frame int    `json:"frame"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "scanin", rcp.Frame)
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

// handleAPISetScanOut 함수는 아이템에 ScanOut 값을 설정한다.
func handleAPISetScanOut(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Frame int    `json:"frame"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "scanout", rcp.Frame)
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

// handleAPISetScanFrame 함수는 아이템에 ScanFrame 값을 설정한다.
func handleAPISetScanFrame(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Frame int    `json:"frame"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "scanframe", rcp.Frame)
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

// handleAPISetHandleIn 함수는 아이템에 HandleIn 값을 설정한다.
func handleAPISetHandleIn(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Frame int    `json:"frame"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "handlein", rcp.Frame)
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

// handleAPISetJustOut 함수는 아이템에 JustOut 값을 설정한다.
func handleAPISetJustOut(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Frame int    `json:"frame"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "justout", rcp.Frame)
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

// handleAPISetHandleOut 함수는 아이템에 HandleOut 값을 설정한다.
func handleAPISetHandleOut(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Frame int    `json:"frame"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	frame := r.FormValue("frame")
	n, err := strconv.Atoi(frame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rcp.Frame = n

	err = SetFrameV2(client, rcp.ID, "handleout", rcp.Frame)
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

// handleAPIPlateSize 함수는 아이템의 PlateSize를 설정한다.
func handleAPISetPlateSize(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Size   string `json:"size"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	size := r.FormValue("size")
	if !regexpImageSize.MatchString(size) {
		http.Error(w, "Please enter in the format of 2048x1152", http.StatusBadRequest)
		return
	}
	rcp.Size = size

	err = SetImageSizeV2(client, rcp.ID, "platesize", rcp.Size)
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

// PostFormValueInList 는 PostForm 쿼리시 Value값이 1개라면 값을 리턴한다.
func PostFormValueInList(key string, values []string) (string, error) {
	if len(values) != 1 {
		return "", errors.New(key + "값이 여러개 입니다")
	}
	if key == "startdate" && values[0] == "" { // Task 시작일은 빈 문자를 허용한다.
		return "", nil
	}
	if key == "predate" && values[0] == "" { // 1차마감일은 빈 문자를 허용한다.
		return "", nil
	}
	if key == "date" && values[0] == "" { // 2차마감일은 빈 문자를 허용한다.
		return "", nil
	}
	if key == "shottype" && values[0] == "" { // 샷타입은 빈 문자를 허용한다.
		return "", nil
	}
	if values[0] == "" {
		return "", errors.New(key + "값이 빈 문자입니다")
	}
	return values[0], nil
}

// handleAPISetCameraPubPath 함수는 아이템의 Camera PubPath를 설정한다.
func handleAPISetCameraPubPath(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Path     string `json:"path"`
		UserID   string `json:"userid"`
		Protocol string `json:"protocol"`
	}
	rcp := Recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "path를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Path = path
	err = SetCameraPubPathV2(client, rcp.ID, rcp.Path)
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

// handleAPISetCameraPubTask 함수는 아이템의 Camera PubTask를 설정한다.
func handleAPISetCameraPubTask(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Task   string `json:"task"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	err = SetCameraPubTaskV2(client, rcp.ID, rcp.Task)
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

// handleAPISetCameraLensmm 함수는 아이템의 Camera Lensmm를 설정한다.
func handleAPISetCameraLensmm(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Lensmm string `json:"lensmm"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	lensmm := r.FormValue("lensmm")
	if lensmm == "" {
		http.Error(w, "need lensmm", http.StatusBadRequest)
		return
	}
	rcp.Lensmm = lensmm
	err = SetCameraLensmmV2(client, rcp.ID, rcp.Lensmm)
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

// handleAPISetCameraProjection 함수는 아이템의 Camera Projection 여부를 설정한다.
func handleAPISetCameraProjection(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID         string `json:"id"`
		Projection bool   `json:"projection"`
		UserID     string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	projection := r.FormValue("projection")
	rcp.Projection = str2bool(projection)
	err = SetCameraProjectionV2(client, rcp.ID, rcp.Projection)
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

// handleAPISetObjectID 함수는 아이템의 ObjectID 값을 설정한다.
func handleAPISetObjectID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		In      int    `json:"in"`
		Out     int    `json:"out"`
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
		case "name":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Name = v
		case "in":
			if len(values) == 1 {
				rcp.In, err = strconv.Atoi(values[0])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				rcp.In = 0
			}
		case "out":
			if len(values) == 1 {
				rcp.Out, err = strconv.Atoi(values[0])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				rcp.Out = 0
			}
		case "userid":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if rcp.UserID == "unknown" && v != "" {
				rcp.UserID = v
			}
		}
	}
	err = SetObjectID(session, rcp.Project, rcp.Name, rcp.In, rcp.Out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("ObjectID: %d - %d\nProject: %s, Name: %s, Author: %s", rcp.In, rcp.Out, rcp.Project, rcp.Name, rcp.UserID))
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

// handleAPISetSeq 함수는 아이템의 seq 값을 설정한다.
func handleAPISetSeq(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Seq    string `json:"seq"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Seq = r.FormValue("seq")

	err = SetSeq(client, rcp.ID, rcp.Seq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetScene 함수는 아이템의 scene 값을 설정한다.
func handleAPISetScene(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Scene    string `json:"scene"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Scene = r.FormValue("scene")

	err = SetScene(client, rcp.ID, rcp.Scene)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetCut 함수는 아이템의 cut 값을 설정한다.
func handleAPISetCut(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Cut    string `json:"cut"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Cut = r.FormValue("cut")

	err = SetCut(client, rcp.ID, rcp.Cut)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetPlatePath 함수는 아이템의 PlatePath 값을 설정한다.
func handleAPISetPlatePath(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Path   string `json:"path"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Path = r.FormValue("path")

	err = SetPlatePathV2(client, rcp.ID, rcp.Path)
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

// handleAPI2SetThummov 함수는 아이템의 Thummov 값을 설정한다.
func handleAPI2SetThummov(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Path   string `json:"path"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	rcp.Path = path

	err = SetThummovV2(client, rcp.ID, rcp.Path)
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

// handleAPISetBeforemov 함수는 아이템의 Before mov 값을 설정한다.
func handleAPISetBeforemov(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Path   string `json:"path"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	rcp.Path = path

	err = SetBeforemovV2(client, rcp.ID, rcp.Path)
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

// handleAPISetAftermov 함수는 아이템의 After mov 값을 설정한다.
func handleAPISetAftermov(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Path   string `json:"path"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "path를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Path = path

	err = SetAftermovV2(client, rcp.ID, rcp.Path)
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

// handleAPISetEditmov 함수는 아이템의 Edit mov 값을 설정한다.
func handleAPISetEditmov(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Path   string `json:"path"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	rcp.Path = path

	err = SetEditmovV2(client, rcp.ID, rcp.Path)
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

// handleAPI2SetTaskStatus 함수는 아이템의 task에 대한 상태를 설정한다.
func handleAPI2SetTaskStatus(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Task   string `json:"task"`
		Status string `json:"status"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Task = r.FormValue("task")
	if rcp.Task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Status = r.FormValue("status")
	if rcp.Status == "" {
		http.Error(w, "need status", http.StatusBadRequest)
		return
	}

	// task가 존재하는지 체크한다.
	err = HasTaskV2(client, rcp.ID, rcp.Task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = SetTaskStatusV3(client, rcp.ID, rcp.Task, rcp.Status)
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

// handleAPIRmTask 함수는 아이템의 task를 제거한다.
func handleAPIRmTask(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Task   string `json:"task"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		if err != nil {
			http.Error(w, "need id", http.StatusBadRequest)
			return
		}
	}
	rcp.ID = id
	rcp.Task = r.FormValue("task")
	if rcp.Task == "" {
		if err != nil {
			http.Error(w, "task가 빈 문자열 입니다", http.StatusBadRequest)
			return
		}
	}
	err = RmTaskV2(client, rcp.ID, rcp.Task)
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

// handleAPIAddTask 함수는 아이템에 task를 추가한다.
func handleAPIAddTask(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Task   string `json:"task"`
		Status string `json:"status"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, s := range status {
		if s.InitStatus {
			rcp.Status = s.ID
		}
	}
	err = AddTaskV2(client, rcp.ID, rcp.Task, rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

// handleAPI2SetTaskUser 함수는 아이템의 task에 대한 유저를 설정한다.
func handleAPISetTaskUser(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID              string `json:"id"`
		Task            string `json:"task"`
		Username        string `json:"username"`
		UsernameAndTeam string `json:"usernameandteam"`
		UserID          string `json:"userid"`
		Error           string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	rcp.Username = r.FormValue("user")
	if rcp.Username != "" {
		rcp.UserID = onlyID(rcp.Username)            // id(name,team) 문자열중 id만 추출한다.
		rcp.UsernameAndTeam = userInfo(rcp.Username) // id(name,team) 문자열을 name,team으로 바꾼다. 웹에서 보기좋게 하기 위함.
		err = SetTaskUserIDV2(client, rcp.ID, rcp.Task, rcp.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = SetTaskUserV3(client, rcp.ID, rcp.Task, rcp.Username)
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

// handleAPISetTaskStartdate 함수는 아이템의 task에 대한 시작일을 설정한다.
func handleAPISetTaskStartdate(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Date   string `json:"date"`
		Task   string `json:"task"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	rcp.ID = r.FormValue("id")
	if rcp.ID == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return

	}
	rcp.Task = r.FormValue("task")
	if rcp.Task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Date = r.FormValue("date") // 마감일이 빈 문자열이 될 수 있다.

	err = HasTaskV2(client, rcp.ID, rcp.Task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = SetTaskStart(client, rcp.ID, rcp.Task, rcp.Date)
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

// handleAPISetTaskUserNote 함수는 아이템의 task에 대한 시작일을 설정한다.
func handleAPISetTaskUserNote(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Task     string `json:"task"`
		UserNote string `json:"usernote"`
		UserID   string `json:"userid"`
		Error    string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	userNote := r.FormValue("usernote")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.UserNote = userNote

	err = SetTaskUserNoteV2(client, rcp.ID, rcp.Task, rcp.UserNote)
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

// handleAPISetDeadline2D 함수는 아이템의 2D 마감일을 설정한다.
func handleAPISetDeadline2D(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID        string `json:"id"`
		Date      string `json:"date"`
		ShortDate string `json:"shortdate"`
		UserID    string `json:"userid"`
		Error     string `json:"error"`
	}
	rcp := Recipe{}

	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	date := r.FormValue("date")
	rcp.Date = date
	err = SetDeadline2DV2(client, rcp.ID, rcp.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// json 으로 결과 전송
	rcp.ShortDate = ToShortTime(rcp.Date) // 웹사이트에 렌더링시 사용한다.
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetDeadline3D 함수는 아이템의 3D 마감일을 설정한다.
func handleAPISetDeadline3D(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID        string `json:"id"`
		Date      string `json:"date"`
		ShortDate string `json:"shortdate"`
		UserID    string `json:"userid"`
		Error     string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	date := r.FormValue("date")
	rcp.Date = date

	err = SetDeadline3DV2(client, rcp.ID, rcp.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json 으로 결과 전송
	rcp.ShortDate = ToShortTime(rcp.Date) // 웹사이트에 렌더링시 사용한다.
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetTaskPredate 함수는 아이템의 task에 대한 1차마감일을 설정한다.
func handleAPISetTaskEnd(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID        string `json:"id"`
		Date      string `json:"date"`
		ShortDate string `json:"shortdate"`
		Task      string `json:"task"`
		UserID    string `json:"userid"`
		Error     string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	rcp.ID = r.FormValue("id")
	if rcp.ID == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return

	}
	rcp.Task = r.FormValue("task")
	if rcp.Task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Date = r.FormValue("date") // 마감일이 빈 문자열이 될 수 있다.

	err = HasTaskV2(client, rcp.ID, rcp.Task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = SetTaskEnd(client, rcp.ID, rcp.Task, rcp.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json 으로 결과 전송
	rcp.ShortDate = ToShortTime(rcp.Date)
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetTaskDate 함수는 아이템의 task에 대한 최종마감일을 설정한다.
func handleAPISetTaskDate(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID        string `json:"id"`
		Date      string `json:"date"`
		ShortDate string `json:"shortdate"`
		Task      string `json:"task"`
		UserID    string `json:"userid"`
		Error     string `json:"error"`
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	rcp.Date = r.FormValue("date")

	err = HasTask(session, rcp.ID, rcp.Task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = SetTaskDate(session, rcp.ID, rcp.Task, rcp.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json 으로 결과 전송
	rcp.ShortDate = ToShortTime(rcp.Date)
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetTaskDuration 함수는 아이템의 task에 대한 시작일, 종료일을 설정한다.
func handleAPISetTaskDuration(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID         string `json:"id"`
		Start      string `json:"start"`
		End        string `json:"end"`
		ShortStart string `json:"shortstart"`
		ShortEnd   string `json:"shortend"`
		Task       string `json:"task"`
		UserID     string `json:"userid"`
		Error      string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	start := r.FormValue("start")
	if start == "" {
		http.Error(w, "need start", http.StatusBadRequest)
		return
	}
	rcp.Start = start
	rcp.ShortStart = ToShortTime(rcp.Start)
	end := r.FormValue("end")
	if end == "" {
		http.Error(w, "need end", http.StatusBadRequest)
		return
	}
	rcp.End = end
	rcp.ShortEnd = ToShortTime(rcp.End)

	err = HasTaskV2(client, rcp.ID, rcp.Task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = SetTaskDurationV2(client, rcp.ID, rcp.Task, rcp.Start, rcp.End)
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

// handleAPISetShotType 함수는 아이템의 shot type을 설정한다.
func handleAPISetShotType(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Type = r.FormValue("shottype")

	err = SetShotTypeV2(client, rcp.ID, rcp.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json 으로 결과 전송
	if rcp.Type == "" { // template에서 렌더링시에는 빈 문자열이면 눈에 보이지 않기 때문에 none으로 반환한다.
		rcp.Type = "none"
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetUseType 함수는 아이템의 Usetype을 설정한다.
func handleAPISetUseType(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	typ := r.FormValue("type")
	if typ == "" {
		http.Error(w, "need type", http.StatusBadRequest)
		return
	}
	rcp.Type = typ
	err = SetUseTypeV2(client, rcp.ID, rcp.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetOutputName 함수는 아이템의 shot의 아웃풋 이름을 설정합니다.
func handleAPISetOutputName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var name string
	var outputname string
	args := r.PostForm
	for key, value := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "name":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			name = v
		case "outputname":
			v, err := PostFormValueInList(key, value)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			outputname = v
		}
	}
	err = SetOutputName(session, project, name, outputname)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	fmt.Fprintf(w, "{\"error\":\"\"}\n")
}

// handleAPISetRetimePlate 함수는 아이템의 retimeplate 값을 설정합니다.
func handleAPISetRetimePlate(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Path     string `json:"path"`
		UserID   string `json:"userid"`
		Error    string `json:"error"`
		Protocol string `json:"protocol"`
	}
	rcp := Recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	rcp.Path = path

	err = SetRetimePlateV2(client, rcp.ID, rcp.Path)
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

// handleAPISetScanname 함수는 아이템의 Scanname 값을 설정합니다.
func handleAPISetScanname(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Scanname string `json:"scanname"`
		UserID   string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Scanname = r.FormValue("scanname")
	err = UpdateItem(client, rcp.ID, "scanname", rcp.Scanname)
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

// handleAPISetAssetType 함수는 아이템의 shot type을 설정한다.
func handleAPISetAssetType(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		OldType string `json:"oldtype"`
		UserID  string `json:"userid"`
		Error   string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	assetType := r.FormValue("assettype")
	if assetType == "" {
		http.Error(w, "need assettype", http.StatusBadRequest)
		return
	}
	rcp.Type = assetType

	beforeType, _, err := SetAssetTypeV2(client, rcp.ID, rcp.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json 으로 결과 전송
	rcp.OldType = beforeType // 브라우저에 기존에 드로잉된 에셋태그를 제거하기 위해서 사용한다.
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPI2SetRnum 함수는 아이템에 롤넘버를 설정합니다.
func handleAPI2SetRnum(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Rnum   string `json:"rnum"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rnum := r.FormValue("rnum")
	if id == "" {
		http.Error(w, "need rnum", http.StatusBadRequest)
		return
	}
	rcp.Rnum = rnum
	if rcp.Rnum != "" && !regexpRnum.MatchString(rcp.Rnum) {
		http.Error(w, rcp.Rnum+"값은 A0001 형식이 아닙니다.", http.StatusBadRequest)
		return
	}
	err = SetRnumV2(client, rcp.ID, rcp.Rnum)
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

// handleAPISetScanTimecodeIn 함수는 아이템에 Scan TimecodeIn 값을 설정한다.
func handleAPISetScanTimecodeIn(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Timecode string `json:"timecode"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Timecode = r.FormValue("timecode")
	err = SetScanTimecodeInV2(client, rcp.ID, rcp.Timecode)
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

// handleAPISetScanTimecodeOut 함수는 아이템에 Scan TimecodeOut 값을 설정한다.
func handleAPISetScanTimecodeOut(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Timecode string `json:"timecode"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Timecode = r.FormValue("timecode")
	err = SetScanTimecodeInV2(client, rcp.ID, rcp.Timecode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = SetScanTimecodeOutV2(client, rcp.ID, rcp.Timecode)
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

// handleAPISetJustTimecodeIn 함수는 아이템에 Just TimecodeIn 값을 설정한다.
func handleAPISetJustTimecodeIn(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Timecode string `json:"timecode"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Timecode = r.FormValue("timecode")

	err = SetJustTimecodeInV2(client, rcp.ID, rcp.Timecode)
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

// handleAPISetJustTimecodeOut 함수는 아이템에 Just TimecodeOut 값을 설정한다.
func handleAPISetJustTimecodeOut(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Timecode string `json:"timecode"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Timecode = r.FormValue("timecode")

	err = SetJustTimecodeOutV2(client, rcp.ID, rcp.Timecode)
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

// handleAPISetFinver 함수는 아이템에 파이널 버전값을 설정한다.
func handleAPISetFinver(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID    string `json:"id"`
		Version string `json:"version"`
		UserID  string `json:"userid"`
		Error   string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need item id", http.StatusBadRequest)
		return
	}
	rcp.ID = id 
	version := r.FormValue("version")
	if version == "" {
		http.Error(w, "need version", http.StatusBadRequest)
		return
	}
	rcp.Version = version 
	err = SetFinverV2(client, rcp.ID, rcp.Version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetFindate 함수는 데이터가 최종으로 나간 날짜를 설정한다.
func handleAPISetFindate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		Date    string `json:"date"`
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
		case "date":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rcp.Date = v
		}
	}
	err = SetFindate(session, rcp.Project, rcp.Name, rcp.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set FinDate: %s\nProject: %s, Name: %s, Author: %s", rcp.Date, rcp.Project, rcp.Name, rcp.UserID))
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

// handleAPISetCrowdAsset 함수는 CrowdAsset을 설정한다.
func handleAPISetCrowdAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project    string `json:"project"`
		Name       string `json:"name"`
		ID         string `json:"id"`
		Crowdasset bool   `json:"crowdasset"`
		UserID     string `json:"userid"`
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
		}
	}
	id, crowdType, err := SetCrowdAsset(session, rcp.Project, rcp.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Crowdasset = crowdType
	rcp.ID = id
	// slack log
	err = slacklog(session, rcp.Project, fmt.Sprintf("Set CrowdType: %t\nProject: %s, Name: %s, Author: %s", rcp.Crowdasset, rcp.Project, rcp.Name, rcp.UserID))
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

// handleAPIAddTag 함수는 아이템에 태그를 설정합니다.
func handleAPIAddTag(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Tag    string `json:"tag"`
		UserID string `json:"userid"`
		Error  string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	tag := r.FormValue("tag")
	if tag == "" {
		http.Error(w, "need tag", http.StatusBadRequest)
		return
	}
	if !regexpTag.MatchString(tag) {
		http.Error(w, "invalid tag rule", http.StatusBadRequest)
		return
	}
	rcp.Tag = tag
	err = AddTagV2(client, rcp.ID, rcp.Tag)
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

// handleAPIAddAssetTag 함수는 아이템에 태그를 설정합니다.
func handleAPIAddAssetTag(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Assettag string `json:"assettag"`
		UserID   string `json:"userid"`
		Error    string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	assettag := r.FormValue("assettag")
	if assettag == "" {
		http.Error(w, "need assettag", http.StatusBadRequest)
		return
	}
	if !regexpTag.MatchString(assettag) {
		http.Error(w, "invalid assettag rule", http.StatusBadRequest)
		return
	}
	rcp.Assettag = assettag
	err = AddAssetTagV2(client, rcp.ID, rcp.Assettag)
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

// handleAPIRenameTag 함수는 아이템의 태그를 변경합니다.
func handleAPIRenameTag(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		Project string `json:"project"`
		Before  string `json:"before"`
		After   string `json:"after"`
		UserID  string `json:"userid"`
	}
	rcp := Recipe{}

	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "need project", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	before := r.FormValue("before")
	if before == "" {
		http.Error(w, "need current tag", http.StatusBadRequest)
		return
	}
	rcp.Before = strings.Replace(before, " ", "", -1) // 빈 공백을 제거한다.
	after := r.FormValue("after")
	if after == "" {
		http.Error(w, "need change tag", http.StatusBadRequest)
		return
	}
	rcp.After = strings.Replace(after, " ", "", -1) // 빈 공백을 제거한다.
	err = RenameTagV2(client, rcp.Project, rcp.Before, rcp.After)
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

// handleAPIRmTag 함수는 아이템에 태그를 삭제합니다.
func handleAPIRmTag(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID        string `json:"id"`
		Tag       string `json:"tag"`
		UserID    string `json:"userid"`
		IsContain bool   `json:"iscontain"`
		Error     string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	tag := r.FormValue("tag")
	if tag == "" {
		http.Error(w, "tag를 설정해주세요", http.StatusBadRequest)
		return
	}
	if !regexpTag.MatchString(tag) {
		http.Error(w, "tag 규칙이 아닙니다", http.StatusBadRequest)
		return
	}
	rcp.Tag = tag
	rcp.IsContain = str2bool(r.FormValue("iscontain"))
	rcp.ID, err = RmTag(client, rcp.ID, rcp.Tag, rcp.IsContain)
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

// handleAPIRmAssetTag 함수는 아이템에 에셋태그를 삭제합니다.
func handleAPIRmAssetTag(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID        string `json:"id"`
		Assettag  string `json:"assettag"`
		UserID    string `json:"userid"`
		IsContain bool   `json:"iscontain"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	assettag := r.FormValue("assettag")
	if assettag == "" {
		http.Error(w, "assettag를 설정해주세요", http.StatusBadRequest)
		return
	}
	if !regexpTag.MatchString(assettag) {
		http.Error(w, "assettag 규칙이 아닙니다", http.StatusBadRequest)
		return
	}
	rcp.Assettag = assettag
	rcp.IsContain = str2bool(r.FormValue("iscontain"))
	err = RmAssetTagV2(client, rcp.ID, rcp.Assettag, rcp.IsContain)
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

// handleAPISetNote 함수는 아이템에 작업내용을 설정합니다.
func handleAPISetNote(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID        string `json:"id"`
		Text      string `json:"text"`
		Overwrite bool   `json:"overwrite"`
		UserID    string `json:"userid"`
		Error     string `json:"error"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Text = r.FormValue("text")
	rcp.Overwrite = str2bool(r.FormValue("overwrite"))

	err = SetNoteV2(client, rcp.ID, rcp.UserID, rcp.Text, rcp.Overwrite)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIAddComment 함수는 아이템에 수정사항을 추가합니다.
func handleAPIAddComment(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID         string `json:"id"`
		Date       string `json:"date"`
		Text       string `json:"text"`
		Media      string `json:"media"`
		MediaTitle string `json:"mediatitle"`
		UserID     string `json:"userid"`
		AuthorName string `json:"authorname"`
		Error      string `json:"error"`
		Protocol   string `json:"protocol"`
	}
	rcp := Recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// 사용자의 이름을 구한다.
	u, err := getUserV2(client, rcp.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) // 사용자가 존재하지 않으면 당연히 Comment를 작성하면 안된다.
		return
	}
	rcp.AuthorName = u.LastNameKor + u.FirstNameKor
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "id을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "text를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Text = text
	rcp.Media = r.FormValue("media")
	rcp.MediaTitle = r.FormValue("mediatitle")
	rcp.Date = time.Now().Format(time.RFC3339)
	err = AddCommentV2(client, rcp.ID, rcp.UserID, rcp.AuthorName, rcp.Date, rcp.Text, rcp.Media, rcp.MediaTitle)
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

// handleAPIEditComment 함수는 아이템에 수정사항을 수정합니다.
func handleAPIEditComment(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID         string `json:"id"`
		Time       string `json:"time"`
		Text       string `json:"text"`
		MediaTitle string `json:"mediatitle"`
		Media      string `json:"media"`
		UserID     string `json:"userid"`
		AuthorName string `json:"authorname"`
		Protocol   string `json:"protocol"`
	}
	rcp := Recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// 사용자의 이름을 구한다.
	u, err := getUserV2(client, rcp.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) // 사용자가 존재하지 않으면 당연히 Comment를 작성하면 안된다.
		return
	}
	rcp.AuthorName = u.LastNameKor + u.FirstNameKor
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	time := r.FormValue("time")
	if time == "" {
		http.Error(w, "time을 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Time = time
	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "text를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Text = text
	rcp.Media = r.FormValue("media")
	rcp.MediaTitle = r.FormValue("mediatitle")
	err = EditCommentV2(client, rcp.ID, rcp.Time, rcp.AuthorName, rcp.Text, rcp.MediaTitle, rcp.Media)
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

// handleAPIRmComment 함수는 아이템에서 수정사항을 삭제합니다.
func handleAPIRmComment(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Date   string `json:"date"`
		Text   string `json:"text"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	date := r.FormValue("date")
	if date == "" {
		http.Error(w, "date를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Date = date
	rcp.ID, rcp.Text, err = RmCommentV2(client, rcp.ID, rcp.UserID, rcp.Date)
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

// handleAPIAddSource 함수는 아이템에 소스를 추가합니다.
func handleAPIAddSource(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Path     string `json:"path"`
		UserID   string `json:"userid"`
		Error    string `json:"error"`
		Protocol string `json:"protocol"`
	}
	rcp := Recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "need title", http.StatusBadRequest)
		return
	}
	rcp.Title = title

	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	rcp.Path = path

	err = AddSourceV2(client, rcp.ID, rcp.UserID, rcp.Title, rcp.Path)
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

// handleAPIRmSource 함수는 아이템에서 링크소스를 삭제합니다.
func handleAPIRmSource(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "need title", http.StatusBadRequest)
		return
	}
	rcp.Title = title

	err = RmSourceV2(client, rcp.ID, rcp.Title)
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

// handleAPIAddReference 함수는 아이템에 레퍼런스를 추가합니다.
func handleAPIAddReference(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Path     string `json:"path"`
		Protocol string `json:"protocol"`
		UserID   string `json:"userid"`
	}
	rcp := Recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "need title", http.StatusBadRequest)
		return
	}
	rcp.Title = title

	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	rcp.Path = path

	err = AddReferenceV2(client, rcp.ID, rcp.UserID, rcp.Title, rcp.Path)
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

// handleAPIRmReference 함수는 아이템에서 레퍼런스를 삭제합니다.
func handleAPIRmReference(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "need title", http.StatusBadRequest)
		return
	}
	rcp.Title = title

	err = RmReferenceV2(client, rcp.ID, rcp.Title)
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

// handleAPISearch 함수는 아이템을 검색합니다.
func handleAPISearch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	var searchword string
	var sortkey string
	args := r.PostForm
	for key, values := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		case "word", "searchword":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			searchword = v
		case "sort", "sortkey":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			sortkey = v
		}
	}

	type recipe struct {
		Data []Item `json:"data"`
	}
	rcp := recipe{}
	searchOp := SearchOption{
		Project:    project,
		Searchword: searchword,
		Sortkey:    sortkey,
	}
	items, err := Search(session, searchOp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	rcp.Data = items
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIDeadline2D 함수는 아이템을 검색합니다.
func handleAPIDeadline2D(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	args := r.PostForm
	for key, values := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		}
	}
	type recipe struct {
		Data []string `json:"data"`
	}
	rcp := recipe{}

	dates, err := DistinctDdline(session, project, "ddline2d")
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	rcp.Data = dates
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPIDeadline3D 함수는 아이템을 검색합니다.
func handleAPIDeadline3D(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	defer session.Close()
	_, _, err = TokenHandler(r, session)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	r.ParseForm() // 받은 문자를 파싱합니다. 파싱되면 map이 됩니다.
	var project string
	args := r.PostForm
	for key, values := range args {
		switch key {
		case "project":
			v, err := PostFormValueInList(key, values)
			if err != nil {
				fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
				return
			}
			project = v
		}
	}
	type recipe struct {
		Data []string `json:"data"`
	}
	rcp := recipe{}

	dates, err := DistinctDdline(session, project, "ddline3d")
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
	rcp.Data = dates
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\"%v\"}\n", err)
		return
	}
}

// handleAPITask 함수는 Task정보를 가지고온다.
func handleAPITask(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID          string `json:"id"`
		UserID      string `json:"userid"`
		RequestTask string `json:"requesttask"`
		Task        `json:"task"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.RequestTask = task

	t, err := GetTaskV2(client, rcp.ID, rcp.RequestTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Task = t
	// 웹에 표시를 위해서 FullTime을 NormalTime으로 변경
	rcp.Task.Start = ToNormalTime(rcp.Task.Start)
	rcp.Task.End = ToNormalTime(rcp.Task.End)
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

// handleAPIShottype 함수는 Shottype 정보를 가지고온다.
func handleAPIShottype(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID       string `json:"id"`
		UserID   string `json:"userid"`
		Shottype string `json:"shottype"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	i, err := getItemV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Shottype = i.Shottype

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

// handleAPIMailInfo 함수는 Email을 전송할 때 필요한 정보를 가지고온다.
func handleAPIMailInfo(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		Project string   `json:"project"`
		ID      string   `json:"id"` // SS_0010_org 형태
		Title   string   `json:"title"`
		Header  string   `json:"header"`
		Mails   []string `json:"mails"`
		Cc      []string `json:"cc"`
		UserID  string   `json:"userid"`
		Lang    string   `json:"lang"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	rcp.Lang = r.FormValue("lang")
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	p, err := getProjectV2(client, rcp.Project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// PM 이메일이 프로젝트 정보에 기입되어있다면 PM에 이메일을 보낼 때 참조한다.
	if regexpEmail.MatchString(p.PmEmail) {
		// 사용자가 아니라 그룹 메일이 설정되어 있을 수 있다.

		rcp.Cc = append(rcp.Cc, fmt.Sprintf("%s<%s>", p.PmEmail, p.PmEmail)) // 한글이름
	} else if regexpUserInfo.MatchString(p.PmEmail) {
		// User가 설정되어 있을 수 있다.
		id := strings.Split(p.PmEmail, "(")[0]
		u, err := getUserV2(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !regexpEmail.MatchString(u.Email) {
			http.Error(w, fmt.Sprintf("%s 사용자는 E-mail 구조를 띄지 않습니다", u.ID), http.StatusBadRequest)
			return
		}
		rcp.Cc = append(rcp.Cc, u.emailString(rcp.Lang))
	}
	// 메일헤더가 빈 문자열이면 프로젝트 id를 메일해더로 사용한다.
	if p.MailHead == "" {
		rcp.Header = rcp.Project
	} else {
		rcp.Header = p.MailHead
	}
	i, err := getItemV2(client, rcp.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Title = i.Name
	for key := range i.Tasks {
		if regexpUserInfo.MatchString(i.Tasks[key].User) { // "khw7096(김한웅,2D1)" 패턴이라면, 앞 문자열이 ID이다.
			id := strings.Split(i.Tasks[key].User, "(")[0]
			if id == rcp.UserID { // 자기 자신에게 이메일을 보내지 않는다.
				continue
			}
			u, err := getUserV2(client, id)
			if err != nil {
				continue
			}
			if !regexpEmail.MatchString(u.Email) {
				continue
			}
			// Task 아티스트의 메일을 메일리스트에 넣는다.
			rcp.Mails = append(rcp.Mails, u.emailString(rcp.Lang))
			// Task 아티스트의 팀명을 찾는다.
			var teamName string
			for _, o := range u.Organizations {
				if o.Primary {
					teamName = o.Team.Name
					break
				}
				teamName = o.Team.Name

			}
			// 만약 팀명이 선언되어있지 않다면, 팀장리스트를 구하지 않는다.
			if teamName == "" {
				continue
			}
			// Task 아티스트의 팀장을 구한다.
			leaderlist1, err := searchUsersV2(client, []string{teamName, "팀장"})
			if err != nil {
				continue
			}
			leaderlist2, err := searchUsersV2(client, []string{teamName, "Lead"})
			if err != nil {
				continue
			}
			// 팀장의 이메일을 참조에 추가한다. 만약 기존 메일리스트에 메일값이 중복되어 있다면, 제거한다.
			for _, leader := range append(leaderlist1, leaderlist2...) {
				if leader.IsLeave { // 퇴사자는 제거한다.
					continue
				}
				has := false
				for _, email := range append(rcp.Mails, rcp.Cc...) {
					if email == leader.Email {
						has = true
					}
				}
				if !has {
					rcp.Cc = append(rcp.Cc, leader.emailString(rcp.Lang))
				}
			}
		}
	}
	// 혹시나 중복된 데이터가 있수 있다 중복을 제거한다.
	rcp.Mails = UniqueSlice(rcp.Mails)
	rcp.Cc = UniqueSlice(rcp.Cc)
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

// handleAPITaskStatusNum 함수는 project, task를 입력받아서 각 status의 갯수를 반환한다.
func handleAPITaskStatusNum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Task    string `json:"task"`
		UserID  string `json:"userid"`
		Infobarnum
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
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "task를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	statusnum, err := TotalTaskStatusnum(session, project, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Infobarnum = statusnum
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPITaskAndUserStatusNum 함수는 project, task, user를 입력받아서 각 status의 갯수를 반환한다.
func handleAPITaskAndUserStatusNum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		Task    string `json:"task"`
		User    string `json:"user"`
		UserID  string `json:"userid"`
		Infobarnum
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
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "task를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	user := r.FormValue("user")
	if task == "" {
		http.Error(w, "task를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.User = user
	rcp.Infobarnum, err = TotalTaskAndUserStatusnum(session, project, task, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIUserStatusNum 함수는 project, user를 입력받아서 각 status의 갯수를 반환한다.
func handleAPIUserStatusNum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		User    string `json:"user"`
		UserID  string `json:"userid"`
		Infobarnum
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
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	user := r.FormValue("user")
	rcp.User = user
	statusnum, err := TotalUserStatusnum(session, project, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Infobarnum = statusnum
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIStatusNum 함수는 project를 입력받아서 각 status의 갯수를 반환한다.
func handleAPIStatusNum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	type Recipe struct {
		Project string `json:"project"`
		UserID  string `json:"userid"`
		Infobarnum
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
	r.ParseForm()
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "project를 설정해주세요", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	statusnum, err := TotalStatusnum(session, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Infobarnum = statusnum
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIAddTaskPublish 함수는 task의 publish 정보를 기록하는 핸들러이다.
func handleAPIAddTaskPublish(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID             string `json:"id"`
		Task           string `json:"task"`
		Key            string `json:"key"`          // Primary Key
		SecondaryKey   string `json:"secondarykey"` // Secondary Key
		Path           string `json:"path"`
		MainVersion    string `json:"mainversion"`
		SubVersion     string `json:"subversion"`
		Subject        string `json:"subject"`
		FileType       string `json:"filetype"`
		KindOfUSD      string `json:"kindofusd"`
		Status         string `json:"status"`
		Createtime     string `json:"createtime"`
		UserID         string `json:"userid"`
		TaskToUse      string `json:"tasktouse"`
		IsOutput       bool   `json:"isoutput"`
		OutputDataPath string `json:"outputdatapath"`
		AuthorNameKor  string `json:"authornamekor"`
	}
	rcp := Recipe{}
	if err := json.NewDecoder(r.Body).Decode(&rcp); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	rcp.UserID, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if rcp.AuthorNameKor == "" {
		user, err := getUserV2(client, rcp.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.AuthorNameKor = user.LastNameKor + user.FirstNameKor
	}
	_, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if rcp.ID == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	if rcp.Task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	if rcp.Key == "" {
		http.Error(w, "need key", http.StatusBadRequest)
		return
	}

	if !HasPublishKeyV2(client, rcp.Key) {
		http.Error(w, rcp.Key+" key는 등록된 키가 아닙니다. 사용할 수 없습니다", http.StatusBadRequest)
		return
	}
	if rcp.Path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	if !(rcp.Status == "usethis" || rcp.Status == "notuse" || rcp.Status == "working") {
		http.Error(w, "status는 usethis, notuse, working 문자열만 사용할 수 있습니다", http.StatusBadRequest)
		return
	}
	// 사용자가 설정한 시간이 있다면 해당시간으로 설정한다.
	_, err = time.Parse(time.RFC3339, rcp.Createtime)
	if err != nil {
		// 시간포멧이 다르다면 현재시간을 입력한다.
		rcp.Createtime = time.Now().Format(time.RFC3339)
	}
	p := Publish{
		SecondaryKey:   rcp.SecondaryKey,
		MainVersion:    rcp.MainVersion,
		SubVersion:     rcp.SubVersion,
		Path:           rcp.Path,
		Subject:        rcp.Subject,
		FileType:       rcp.FileType,
		KindOfUSD:      rcp.KindOfUSD,
		Status:         rcp.Status,
		Createtime:     rcp.Createtime,
		TaskToUse:      rcp.TaskToUse,
		IsOutput:       rcp.IsOutput,
		AuthorNameKor:  rcp.AuthorNameKor,
		OutputDataPath: rcp.OutputDataPath,
	}
	fmt.Println(rcp)
	err = addTaskPublishV2(client, rcp.ID, rcp.Task, rcp.Key, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("ep")
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIRmTaskPublishKey 함수는 task의 publish key정보를 삭제하는 핸들러이다.
func handleAPIRmTaskPublishKey(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID     string `json:"id"`
		Task   string `json:"task"`
		Key    string `json:"key"`
		UserID string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "need key", http.StatusBadRequest)
		return
	}
	rcp.Key = key
	err = rmTaskPublishKeyV2(client, id, task, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIRmTaskPublish 함수는 task의 publish 정보중 하나를 삭제한다.
func handleAPIRmTaskPublish(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID         string `json:"id"`
		Task       string `json:"task"`
		Key        string `json:"key"`
		Createtime string `json:"createtime"`
		Path       string `json:"path"`
		UserID     string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "need key", http.StatusBadRequest)
		return
	}
	rcp.Key = key
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	rcp.Path = path
	createtime := r.FormValue("createtime")
	if createtime == "" {
		http.Error(w, "need createtime", http.StatusBadRequest)
		return
	}
	rcp.Createtime = createtime

	// Item 가져오기
	item, err := getItemV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// task가 존재하는지 체크
	hasTask := false
	for _, t := range item.Tasks {
		if t.Title == task {
			hasTask = true
		}
	}
	if !hasTask {
		http.Error(w, task+" Task가 존재하지 않습니다", http.StatusBadRequest)
		return
	}
	// Publish Primary key가 존재하는지 체크
	// path가 존재하는지 체크
	// createtime이 존재하는지 체크
	hasKey := false
	hasTime := false
	hasPath := false
	for k, pubList := range item.Tasks[task].Publishes {
		if key != k {
			continue
		}
		hasKey = true
		for _, p := range pubList {
			if p.Path != path {
				continue
			}
			hasPath = true
			if p.Createtime == createtime {
				hasTime = true
			}
		}
	}
	if !hasKey {
		http.Error(w, key+" key로 Publish한 데이터가 존재하지 않습니다", http.StatusInternalServerError)
		return
	}
	if !hasPath {
		http.Error(w, path+" path로 Publish한 데이터가 존재하지 않습니다", http.StatusInternalServerError)
		return
	}
	if !hasTime {
		http.Error(w, createtime+" 시간으로 Publish한 데이터가 존재하지 않습니다", http.StatusInternalServerError)
		return
	}
	// 에러처리가 끝나면 해당 publish를 지운다.
	err = rmTaskPublishV2(client, id, task, key, createtime, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPITaskPublishStatus 함수는 task > publish > status 정보를 변경한다.
func handleAPISetTaskPublishStatus(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		Project    string `json:"project"`
		ID         string `json:"id"`
		Task       string `json:"task"`
		Key        string `json:"key"`
		Status     string `json:"status"`
		Path       string `json:"path"`
		Createtime string `json:"createtime"`
		UserID     string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, "need project", http.StatusBadRequest)
		return
	}
	rcp.Project = project
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	task := r.FormValue("task")
	if task == "" {
		http.Error(w, "need task", http.StatusBadRequest)
		return
	}
	rcp.Task = task
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "need key", http.StatusBadRequest)
		return
	}
	rcp.Key = key
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path", http.StatusBadRequest)
		return
	}
	rcp.Path = path
	createtime := r.FormValue("createtime")
	if createtime == "" {
		http.Error(w, "need createtime", http.StatusBadRequest)
		return
	}
	rcp.Createtime = createtime
	status := r.FormValue("status")
	if !(status == "usethis" || status == "notuse" || status == "working") {
		http.Error(w, "status는 usethis, notuse, working 문자열만 사용할 수 있습니다", http.StatusBadRequest)
		return
	}
	rcp.Status = status
	i, err := getItemV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for n, p := range i.Tasks[task].Publishes[key] {
		if p.Createtime == rcp.Createtime && p.Path == rcp.Path {
			i.Tasks[task].Publishes[key][n].Status = rcp.Status
		}
	}
	err = setItemV2(client, i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPISetEpisode 함수는 아이템의 episode 값을 설정한다.
func handleAPISetEpisode(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		ID      string `json:"id"`
		Episode string `json:"episode"`
		UserID  string `json:"userid"`
	}
	rcp := Recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.UserID, _, err = TokenHandlerV2(r, client)
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

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id
	rcp.Episode = r.FormValue("episode")

	err = SetEpisodeV2(client, rcp.ID, rcp.Episode)
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
