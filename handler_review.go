package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func handleDailyReviewStatus(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	today := time.Now().Format("2006-01-02")
	// 오늘날짜를 구하고 리다이렉트한다.
	http.Redirect(w, r, "/reviewstatus?searchword=createtime:"+today, http.StatusSeeOther)
}

// handleReviewData 함수는 리뷰 영상데이터를 전송한다.
func handleReviewData(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 2 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	ext := q.Get("ext") // 확장자를 자동으로 가지고 오지 않는 이유는 DB에 접근하는것을 줄이기 위해서이다.
	if ext == "" {
		ext = ".mp4" // 확장자가 없다면 기본적으로 mp4를 불러온다.
	}
	http.ServeFile(w, r, fmt.Sprintf("%s/%s%s", CachedAdminSetting.ReviewDataPath, id, ext))
}

// handleReviewDrawingData 함수는 리뷰 드로잉 데이터를 전송한다.
func handleReviewDrawingData(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel < 2 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	q := r.URL.Query()
	id := q.Get("id")
	frameStr := q.Get("frame")
	_ = q.Get("time") // 브라우저에서 캐쉬되지 않은 이미지를 가지고 오기 위해 time 옵션을 사용한다. 리소스 URL로 이미지를 캐쉬하기 때문이다.
	frame, err := strconv.Atoi(frameStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.ServeFile(w, r, fmt.Sprintf("%s/%s.%06d.png", CachedAdminSetting.ReviewDataPath, id, frame))
}

// handleReviewStatusSubmit 함수는 리뷰 검색창의 검색어를 입력받아 새로운 URI로 리다이렉션 한다.
func handleReviewStatusSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Post Only", http.StatusMethodNotAllowed)
		return
	}
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	searchword := r.FormValue("searchword")
	reviewProject := r.FormValue("reviewproject")
	itemStatus := r.FormValue("itemstatus")
	createtime := r.FormValue("review-createtime")
	task := r.FormValue("review-task")
	redirectURL := fmt.Sprintf("/reviewstatus?searchword=%s&project=%s&itemstatus=%s&createtime=%s&task=%s", searchword, reviewProject, itemStatus, createtime, task)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// handleUploadReviewFile 함수는 리뷰 파일을 업로드 한다.
func handleUploadReviewFile(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		Path     string `json:"path"`
		Unixtime string `json:"unixtime"`
		Project  string `json:"project"`
		Vendor   string `json:"vendor"`
		Partner  string `json:"partner"`
		Task     string `json:"task"`
		Year     string `json:"year"`
		Month    string `json:"month"`
		Day      string `json:"day"`
		Type     string `json:"type"`
		Ext      string `json:"ext"`
	}
	rcp := Recipe{}
	rcp.Unixtime = fmt.Sprintf("%d", time.Now().Unix())
	// MultipartForm을 파싱합니다.
	buffer := CachedAdminSetting.MultipartFormBufferSize // 이 절차를 위해 매번 DB에 접근하지 않기 위해서 CachedAdminSetting을 이용합니다.
	err := r.ParseMultipartForm(int64(buffer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Project = r.FormValue("project")
	rcp.Vendor = r.FormValue("vendor")
	rcp.Partner = r.FormValue("partner")
	rcp.Task = r.FormValue("task")
	y, m, d := time.Now().Date()
	rcp.Year = fmt.Sprintf("%04d", y)
	rcp.Month = fmt.Sprintf("%02d", int(m))
	rcp.Day = fmt.Sprintf("%02d", d)
	for _, files := range r.MultipartForm.File {
		for _, f := range files {
			if f.Size == 0 {
				http.Error(w, "파일사이즈가 0 바이트입니다", http.StatusInternalServerError)
				return
			}
			file, err := f.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				continue
			}
			defer file.Close()
			mimeType := f.Header.Get("Content-Type")
			ext := strings.ToLower(filepath.Ext(f.Filename))
			switch mimeType {
			case "image/jpeg":
				rcp.Type = "image"
				rcp.Ext = ".jpg"
				if !(ext == ".jpg" || ext == ".jpeg") {
					http.Error(w, ".jpg, .png 이미지만 허용합니다", http.StatusBadRequest)
					return
				}
			case "image/png":
				rcp.Type = "image"
				rcp.Ext = ".png"
				if !(ext == ".png") {
					http.Error(w, ".jpg, .png 이미지만 허용합니다", http.StatusBadRequest)
					return
				}
			case "video/quicktime", "video/mp4":
				rcp.Type = "clip"
				rcp.Ext = ".mp4"
				if !(ext == ".mov" || ext == ".mp4") { // .mov, .mp4 외에는 허용하지 않는다.
					http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
					return
				}
			default:
				//허용하지 않는 파일 포맷입니다.
				http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
				return
			}
			data, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if CachedAdminSetting.ReviewUploadPathPermission == "" {
				http.Error(w, "Need Upload path and file permission in adminsetting.", http.StatusBadRequest)
				return
			}
			filePerm, err := strconv.ParseInt(CachedAdminSetting.ReviewUploadPathPermission, 8, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var path string
			if CachedAdminSetting.ReviewUploadPath == "" {
				path, err = os.MkdirTemp("", "")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				// Review Upload Path를 파싱한다.
				var reviewUploadPath bytes.Buffer
				reviewUploadPathTmpl, err := template.New("reviewUploadPath").Parse(CachedAdminSetting.ReviewUploadPath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = reviewUploadPathTmpl.Execute(&reviewUploadPath, rcp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				path = reviewUploadPath.String()
				// path가 존재하지 않으면 생성한다.
				if _, err := os.Stat(path); os.IsNotExist(err) {
					err = os.MkdirAll(path, os.FileMode(filePerm))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				uid, err := strconv.Atoi(CachedAdminSetting.ReviewUploadPathUID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				gid, err := strconv.Atoi(CachedAdminSetting.ReviewUploadPathGID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = os.Chown(path, uid, gid)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			err = os.WriteFile(path+"/"+f.Filename, data, os.FileMode(filePerm))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rcp.Path = path + "/" + f.Filename
		}
	}
	// 업로드 경로를 리턴합니다. Dropzone에서 활용하기 위해 json으로 반환합니다.
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// handleReviewStatus 함수는 Status 방식의 리뷰 페이지이다.
func handleReviewStatus(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	type recipe struct {
		User        User
		Projectlist []string
		SearchOption
		Searchword       string
		Status           []Status // css 생성을 위해서 필요함
		CurrentReview    Review   // 현재 리뷰 자료구조
		Reviews          []Review // 옆 Review 항목
		ReviewGroup      []Review // 하단 Review 항목
		TasksettingNames []string
		Project          string
		ItemStatus       string
		Createtime       string
		Task             string
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.Project = q.Get("project")
	rcp.ItemStatus = q.Get("itemstatus")
	rcp.Createtime = q.Get("createtime")
	rcp.Task = q.Get("task")

	rcp.Searchword = q.Get("searchword")
	id := q.Get("id")
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
	rcp.Projectlist, err = ProjectlistV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.TasksettingNames, err = TaskSettingNamesV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Status, err = AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Searchword = setSearchFilter(rcp.Searchword, "project", rcp.Project)
	rcp.Searchword = setSearchFilter(rcp.Searchword, "itemstatus", rcp.ItemStatus)
	rcp.Searchword = setSearchFilter(rcp.Searchword, "task", rcp.Task)
	rcp.Searchword = setSearchFilter(rcp.Searchword, "createtime", rcp.Createtime)

	rcp.Reviews, err = searchReviewV2(client, rcp.Searchword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if id != "" {
		review, err := getReviewV2(client, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rcp.CurrentReview = review
		rcp.ReviewGroup, err = searchReviewV2(client, fmt.Sprintf("project:%s name:%s", review.Project, review.Name))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "reviewstatus", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
