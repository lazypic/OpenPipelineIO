package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"encoding/json"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	// 만약 사용자에게 AccessProjects가 설정되어있다면 해당리스트를 사용한다.
	if len(rcp.User.AccessProjects) != 0 {
		var accessProjects []string
		for _, i := range rcp.Projectlist {
			for _, j := range rcp.User.AccessProjects {
				if i != j {
					continue
				}
				accessProjects = append(accessProjects, j)
			}
		}
		rcp.Projectlist = accessProjects
	}

	rcp.Setting = CachedAdminSetting
	err = TEMPLATES.ExecuteTemplate(w, "directupload", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UploadStatus struct {
	FileName  string `json:"fileName"`
	SavedPath string `json:"savedPath"`
	Progress  int    `json:"progress"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func directUploadHandler(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	r.ParseMultipartForm(10 << 30) // 10G 제한, 20 == 10M
	files := r.MultipartForm.File["files"]
	relativePaths := r.MultipartForm.Value["relativePath[]"] // 폴더 구조 유지
	projects := r.MultipartForm.Value["project"] // 폴더 구조 유지

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
		targetPath := CachedAdminSetting.DirectUploadPath

		// add CompanyID path
		if CachedAdminSetting.EnableDirectuploadWithCompanyID {
			u, err := getUserV2(client, ssid.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if u.CompanyID != "" {
				targetPath = filepath.Join(targetPath, u.CompanyID)
			} else {
				targetPath = filepath.Join(targetPath, u.ID)
			}
		}

		// add project path
		if len(projects) > 0 {
			selectProject := projects[0]
			if selectProject != "" {
				targetPath = filepath.Join(targetPath, selectProject)
			}
		}

		savePath := filepath.Join(targetPath, relativePaths[i])
		os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
		outFile, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "파일 저장 실패", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// 파일을 청크 단위로 저장하며 진행률 전송
		buffer := make([]byte, 1024*64) // 64KB 버퍼
		totalSize := fileHeader.Size
		written := int64(0)

		for {
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				http.Error(w, "파일 읽기 실패", http.StatusInternalServerError)
				return
			}
			if n == 0 {
				break
			}

			// 파일 저장
			if _, err := outFile.Write(buffer[:n]); err != nil {
				http.Error(w, "파일 쓰기 실패", http.StatusInternalServerError)
				return
			}
			written += int64(n)

			// 진행률 계산 및 WebSocket으로 전송
			progress := int((written * 100) / totalSize)
			directUploadBroadcastProgress(UploadStatus{
				FileName:  fileHeader.Filename,
				SavedPath: savePath,
				Progress:  progress,
			})

		}
		// 완료 시 100% 전송
		directUploadBroadcastProgress(UploadStatus{
			FileName:  fileHeader.Filename,
			SavedPath: savePath,
			Progress:  100,
		})
	}

	jsonResponse, _ := json.Marshal(uploadedFiles)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

var clients = make(map[*websocket.Conn]bool)
var clientsLock sync.Mutex

func directUploadProgressHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	clientsLock.Lock()
	clients[conn] = true
	clientsLock.Unlock()

	fmt.Println("WebSocket client connected for upload progress")

	// WebSocket 연결이 끊어지면 클라이언트 목록에서 삭제
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			clientsLock.Lock()
			delete(clients, conn)
			clientsLock.Unlock()
			fmt.Println("WebSocket client disconnected")
			break
		}
	}
}

func directUploadBroadcastProgress(status UploadStatus) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	for client := range clients {
		err := client.WriteJSON(status)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
