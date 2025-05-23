package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleScanPlate(w http.ResponseWriter, r *http.Request) {
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
		SearchOption
		Setting          Setting
		Projectlist      []string
		TasksettingNames []string
		Status           []Status
		Colorspaces      []string `json:"colorspaces"`
	}
	rcp := recipe{}
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
	rcp.User, err = getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Setting = CachedAdminSetting
	rcp.Colorspaces, err = loadOCIOConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "scanplate", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
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
		SearchOption
		Setting              Setting
		Projectlist          []string
		TasksettingNames     []string
		Status               []Status
		ProcessingScanPlates []ScanPlate
	}
	rcp := recipe{}
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
	rcp.User, err = getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Setting = CachedAdminSetting
	rcp.ProcessingScanPlates, err = GetUnDoneScanPlate(client)
	//rcp.ProcessingScanPlates, err = allScanPlate(client)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = TEMPLATES.ExecuteTemplate(w, "process", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUploadScanPlate(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}

	buffer := CachedAdminSetting.MultipartFormBufferSize
	err = r.ParseMultipartForm(int64(buffer)) // grab the multipart form
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path := CachedAdminSetting.ScanPlateUploadPath + "/temp"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0775)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	type recipe struct {
		Path     string `json:"path"`
		Ext      string `json:"ext"`
		Filename string `json:"filename"`
	}
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
			switch mimeType {
			case "image/jpeg", "image/png", "video/quicktime", "image/x-exr", "application/octet-stream":
				ext := strings.ToLower(filepath.Ext(f.Filename))
				if ext != ".jpg" && ext != ".png" && ext != ".mov" && ext != ".dpx" && ext != ".exr" { // .jpg .dpx .exr 외에는 허용하지 않는다.
					http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
					return
				}
				data, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = ioutil.WriteFile(path+"/"+f.Filename, data, 0775)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				rcp := recipe{}
				rcp.Path = path
				rcp.Filename = f.Filename
				rcp.Ext = ext
				json, err := json.Marshal(rcp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Write(json)
			default:
				http.Error(w, "허용하지 않는 파일 포맷입니다", http.StatusBadRequest)
				return
			}
		}
	}
}

func deleteScanPlateTemp(w http.ResponseWriter, r *http.Request) {
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
	// token 체크
	_, accessLevel, err := TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accessLevel != AdminAccessLevel {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	path := CachedAdminSetting.ScanPlateUploadPath + "/temp"

	os.RemoveAll(path)
	os.MkdirAll(path, 0775)

	type recipe struct {
		Path string `json:"path"`
	}
	rcp := recipe{}
	rcp.Path = path
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPISearchFootages(w http.ResponseWriter, r *http.Request) {
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
	// token 체크
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "need path value", http.StatusBadRequest)
		return
	}
	typ := r.FormValue("type")
	episode := r.FormValue("episode")

	method := r.FormValue("method")
	incolorspace := r.FormValue("incolorspace")
	outcolorspace := r.FormValue("outcolorspace")
	regex := r.FormValue("regex")

	regexp, regerr := regexp.Compile(regex)
	if regerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var items []ScanPlate
	if method == "" || method == "image" {
		items, err = searchImage(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if method == "mov" {
		items, err = searchMovs(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		items, err = searchSeq(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// set value in scanplate items.
	for n := range items {
		items[n].Episode = episode
		items[n].Type = typ
		items[n].InColorspace = incolorspace
		items[n].OutColorspace = outcolorspace
		regGroup := regexp.SubexpNames()
		regValue := regexp.FindStringSubmatch(items[n].Base)
		// 사용자에게 받은 Regex값을 이용해서 Name값을 추출한다.
		if regerr == nil {
			for order, groupName := range regGroup {
				switch groupName {
				case "name", "Name":
					items[n].Name = regValue[order]
				case "type", "Type":
					if typ != "" {
						items[n].Type = typ
					} else {
						items[n].Type = regValue[order]
					}
				}
			}
		}
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

func handleAPIOcioColorspace(w http.ResponseWriter, r *http.Request) {
	type recipe struct {
		Colorspaces []string `json:"colorspaces"`
	}
	rcp := recipe{}
	colorspace, err := loadOCIOConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Colorspaces = colorspace
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIProtocol(w http.ResponseWriter, r *http.Request) {
	type recipe struct {
		Protocol string `json:"protocol"`
	}
	rcp := recipe{}
	rcp.Protocol = CachedAdminSetting.Protocol
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func postHandleAPIScanPlate(w http.ResponseWriter, r *http.Request) {
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
	// token 체크
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scan := ScanPlate{}
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&scan)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	scan.ID = primitive.NewObjectID()
	// 스캔플레이트를 등록합니다.
	err = addScanPlate(client, scan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Json 값 출력
	data, err := json.Marshal(scan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func deleteHandleAPIScanPlate(w http.ResponseWriter, r *http.Request) {
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
	// token 체크
	_, accessLevel, err := TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accessLevel != AdminAccessLevel {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	// 어떠한 데이터가 삭제되었는지 확인하기 위해서 구한다.
	s, err := getScanPlate(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = rmScanPlate(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func putHandleAPIScanPlate(w http.ResponseWriter, r *http.Request) {
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
	// token 체크
	_, accessLevel, err := TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accessLevel != AdminAccessLevel {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "need id", http.StatusBadRequest)
		return
	}
	// 어떠한 데이터가 삭제되었는지 확인하기 위해서 구한다.
	s, err := getScanPlate(client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.ParseForm()
	processStatus := r.FormValue("processstatus")
	if processStatus == "" {
		http.Error(w, "need processstatus", http.StatusBadRequest)
		return
	}
	if !(processStatus == "wait" || processStatus == "processing" || processStatus == "done" || processStatus == "error") {
		http.Error(w, "support status: wait,processing,done,error", http.StatusBadRequest)
		return
	}
	s.ProcessStatus = processStatus
	err = setScanPlate(client, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIScanPlates(w http.ResponseWriter, r *http.Request) {
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
	// token 체크
	_, _, err = TokenHandlerV2(r, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scans := []ScanPlate{}
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&scans)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	for _, scan := range scans {
		scan.ID = primitive.NewObjectID()
		err = addScanPlate(client, scan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Json 값 출력
	data, err := json.Marshal(scans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
