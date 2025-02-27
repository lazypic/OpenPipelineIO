package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"encoding/json"
)

func handleDirectupload(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type recipe struct {
		User
		Setting     Setting
		Projectlist []string
	}
	rcp := recipe{}
	rcp.Projectlist, err = ProjectlistV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User, err = getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Setting = CachedAdminSetting
	err = TEMPLATES.ExecuteTemplate(w, "directupload", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UploadStatus struct {
	FileName string `json:"fileName"`
	SavedPath string `json:"savedPath"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func directUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 30) // 10G 제한, 20 == 10M
	files := r.MultipartForm.File["files"]
	relativePaths := r.MultipartForm.Value["relativePath[]"] // 폴더 구조 유지

	var uploadedFiles []UploadStatus

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "파일 열기 실패", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// 저장할 경로 설정
		if CachedAdminSetting.DirectUploadPath == "" {
			http.Error(w, "need setup direct upload path", http.StatusInternalServerError)
			return
		}
		savePath := filepath.Join(CachedAdminSetting.DirectUploadPath, relativePaths[i])
		os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
		outFile, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "파일 저장 실패", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// 파일 저장
		io.Copy(outFile, file)

		uploadedFiles = append(uploadedFiles, UploadStatus{FileName: fileHeader.Filename, SavedPath: savePath})

	}
	jsonResponse, _ := json.Marshal(uploadedFiles)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

var clients = make(map[*websocket.Conn]bool)

func directUploadProgressHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	clients[conn] = true
}

func directUploadBroadcastProgress(progress int) {
	for client := range clients {
		err := client.WriteJSON(map[string]int{"progress": progress})
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
