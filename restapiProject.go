package main

import (
	"context"
	"encoding/json"
	"net/http"
)

// handleAPIAddproject 함수는 프로젝트를 추가한다.
func handleAPIAddproject(w http.ResponseWriter, r *http.Request) {
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
	id := r.FormValue("id")
	p := *NewProject(id)
	err = addProjectV2(client, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project, err := getProjectV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIProject 함수는 프로젝트 정보를 불러온다.
func handleAPIProject(w http.ResponseWriter, r *http.Request) {
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
	project, err := getProjectV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIProjectTags 함수는 프로젝트에 사용되는 태그리스트를 불러온다.
func handleAPIProjectTags(w http.ResponseWriter, r *http.Request) {
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
	project := q.Get("project")
	_, err = getProjectV2(client, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tags, err := DistinctV2(client, project, "tag")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 태그중에 빈 문자열을 제거한다.
	var filteredTags []string
	for _, t := range tags {
		if t == "" {
			continue
		}
		filteredTags = append(filteredTags, t)
	}
	data, err := json.Marshal(filteredTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIProjectAssetTags 함수는 프로젝트에 사용되는 에셋 태그리스트를 불러온다.
func handleAPIProjectAssetTags(w http.ResponseWriter, r *http.Request) {
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
	project := q.Get("project")
	_, err = getProjectV2(client, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	assettags, err := DistinctV2(client, project, "assettags")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 에셋태그 중 빈 문자열을 제거한다.
	var filteredTags []string
	for _, t := range assettags {
		if t == "" {
			continue
		}
		filteredTags = append(filteredTags, t)
	}
	data, err := json.Marshal(filteredTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleAPIProjects 함수는 프로젝트 리스트를 반환한다.
func handleAPI2Projects(w http.ResponseWriter, r *http.Request) {
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
	projects, err := getProjectsV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
