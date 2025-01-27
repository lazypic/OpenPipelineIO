package main

import (
	"net/http"
	"context"
)

// handleAddPublishKey 함수는 PublishKey를 추가하는 페이지이다.
func handleAddPublishKey(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}

	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	
	type recipe struct {
		User User
		SearchOption
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookieV2(client, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addpublishkey", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddPublishKeySubmit 함수는 PublishKey를 추가합니다.
func handleAddPublishKeySubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	
	key := PublishKey{
		ID:          r.FormValue("id"),
		Description: r.FormValue("description"),
	}
	err = key.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = AddPublishKeyV2(client, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/publishkey", http.StatusSeeOther)
}

// handlePublishKey 함수는 PublishKey 를 보는 페이지이다.
func handlePublishKey(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}

	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	
	type recipe struct {
		User        User
		PublishKeys []PublishKey
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.PublishKeys, err = AllPublishKeysV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "publishkey", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleRmPublishKey 함수는 PublishKey를 삭제하는 페이지이다.
func handleRmPublishKey(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 5 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	
	type recipe struct {
		User User
		SearchOption
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookieV2(client, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "rmpublishkey", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleRmPublishKeySubmit 함수는 PublishKey를 삭제합니다.
func handleRmPublishKeySubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 5 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}

	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	
	id := r.FormValue("id")
	err = RmPublishKeyV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/publishkey", http.StatusSeeOther)
}

// handleEditPublishKey 함수는 PublishKey를 편집하는 페이지이다.
func handleEditPublishKey(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	
	type recipe struct {
		User User
		SearchOption
		PublishKey
		Setting Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookieV2(client, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	q := r.URL.Query()
	id := q.Get("id")
	rcp.PublishKey, err = GetPublishKeyV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "editpublishkey", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEditPublishKeySubmit 함수는 PublishKey를 수정합니다.
func handleEditPublishKeySubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 4 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	
	key := PublishKey{
		ID:          r.FormValue("id"),
		Description: r.FormValue("description"),
	}
	err = key.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = SetPublishKeyV2(client, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/publishkey", http.StatusSeeOther)
}
