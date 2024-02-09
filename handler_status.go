package main

import (
	"context"
	"net/http"
	"strconv"
)

// handleAddStatus 함수는 Status를 추가하는 페이지이다.
func handleAddStatus(w http.ResponseWriter, r *http.Request) {
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
	err = TEMPLATES.ExecuteTemplate(w, "addstatus", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddStatusSubmit 함수는 Status를 추가합니다.
func handleAddStatusSubmit(w http.ResponseWriter, r *http.Request) {
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

	order, err := strconv.ParseFloat(r.FormValue("order"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := Status{
		ID:          r.FormValue("id"),
		Description: r.FormValue("description"),
		TextColor:   "#" + r.FormValue("textcolor"),
		BGColor:     "#" + r.FormValue("bgcolor"),
		BorderColor: "#" + r.FormValue("bordercolor"),
		Order:       order,
		DefaultOn:   str2bool(r.FormValue("defaulton")),
		InitStatus:  str2bool(r.FormValue("initstatus")),
	}
	err = s.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = AddStatusV2(client, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/status", http.StatusSeeOther)
}

// handleStatus 함수는 Status 를 보는 페이지이다.
func handleStatus(w http.ResponseWriter, r *http.Request) {
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
		User   User
		Status []Status
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.Status, err = AllStatusV2(client)
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
	err = TEMPLATES.ExecuteTemplate(w, "status", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleRmStatus 함수는 Status를 삭제하는 페이지이다.
func handleRmStatus(w http.ResponseWriter, r *http.Request) {
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
	err = TEMPLATES.ExecuteTemplate(w, "rmstatus", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleRmStatusSubmit 함수는 Status를 삭제합니다.
func handleRmStatusSubmit(w http.ResponseWriter, r *http.Request) {
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
	err = RmStatusV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/status", http.StatusSeeOther)
}

// handleEditStatus 함수는 Status를 편집하는 페이지이다.
func handleEditStatus(w http.ResponseWriter, r *http.Request) {
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
		Status
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
	rcp.Status, err = GetStatusV2(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "editstatus", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleEditStatusSubmit 함수는 Status를 추가합니다.
func handleEditStatusSubmit(w http.ResponseWriter, r *http.Request) {
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

	order, err := strconv.ParseFloat(r.FormValue("order"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := Status{
		ID:          r.FormValue("id"),
		Description: r.FormValue("description"),
		TextColor:   "#" + r.FormValue("textcolor"),
		BGColor:     "#" + r.FormValue("bgcolor"),
		BorderColor: "#" + r.FormValue("bordercolor"),
		Order:       order,
		DefaultOn:   str2bool(r.FormValue("defaulton")),
		InitStatus:  str2bool(r.FormValue("initstatus")),
	}
	err = s.CheckError()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = SetStatusV2(client, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/status", http.StatusSeeOther)
}
