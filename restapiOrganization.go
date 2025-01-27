package main

import (
	"context"
	"encoding/json"
	"net/http"
)

// handleAPIAllTeams 함수는 모든 팀 조직정보를 반환한다.
func handleAPIAllTeams(w http.ResponseWriter, r *http.Request) {
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
	teams, err := allTeamsV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipe struct {
		Data []Team `json:"data"`
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rcp := recipe{}
	rcp.Data = teams
	err = json.NewEncoder(w).Encode(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
