package main

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/disintegration/imaging"
)

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}


// handleAPIUploadThumbnail 함수는 thumbnail 이미지를 업로드 하는 RestAPI 이다.
func handleAPIUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	type Recipe struct {
		Project        string `json:"project"`
		ID             string `json:"id"`
		Path           string `json:"path"`
		UploadFilename string `json:"uploadfilename"`
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
	// 어드민 셋팅을 불러온다.
	uid, err := strconv.Atoi(CachedAdminSetting.ThumbnailImagePathUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gid, err := strconv.Atoi(CachedAdminSetting.ThumbnailImagePathGID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	permission, err := strconv.ParseInt(CachedAdminSetting.ThumbnailImagePathPermission, 8, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 폼을 분석한다.
	err = r.ParseMultipartForm(int64(CachedAdminSetting.MultipartFormBufferSize))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project := r.FormValue("project")
	if project == "" {
		writeJSONError(w, "need project", http.StatusBadRequest)
		return
	}
	rcp.Project = project

	id := r.FormValue("id")
	if id == "" {
		writeJSONError(w, "need id", http.StatusBadRequest)
		return
	}
	rcp.ID = id

	files := r.MultipartForm.File["image"]

	if len(files) != 1 { // 파일이 없다면 에러처리한다.
		writeJSONError(w, "multiple file cannot be set", http.StatusBadRequest)
		return
	}

	f := files[0]
	rcp.UploadFilename = f.Filename // 파일명을 추출한다.
	if f.Size == 0 {
		writeJSONError(w, "file size is 0 bytes", http.StatusBadRequest)
		return
	}
	file, err := f.Open()
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		writeJSONError(w, "invalid image file", http.StatusBadRequest)
		return
	}
	if format != "jpeg" && format != "png" {
		writeJSONError(w, "unsupported image format: "+format, http.StatusBadRequest)
		return
	}


	// adminsetting에 설정된 썸네일 템플릿에 실제 값을 넣는다.
	var thumbImgPath bytes.Buffer
	thumbImgPathTmpl, err := template.New("thumbImgPath").Parse(CachedAdminSetting.ThumbnailImagePath)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = thumbImgPathTmpl.Execute(&thumbImgPath, rcp)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 썸네일 이미지가 이미 존재하는 경우 이미지 파일을 지운다.
	if _, err := os.Stat(thumbImgPath.String()); os.IsExist(err) {
		err = os.Remove(thumbImgPath.String())
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 썸네일 경로를 생성한다.
	path, _ := path.Split(thumbImgPath.String())
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 폴더를 생성한다.
		err = os.MkdirAll(path, os.FileMode(permission))
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 위 폴더가 잘 생성되어 존재한다면 폴더의 권한을 설정한다.
		if _, err := os.Stat(path); os.IsExist(err) {
			err = os.Chown(path, uid, gid)
			if err != nil {
				writeJSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	// 사용자가 업로드한 데이터를 이미지 자료구조로 만들고 리사이즈 한다.
	resizedImage := imaging.Fill(img, CachedAdminSetting.ThumbnailImageWidth, CachedAdminSetting.ThumbnailImageHeight, imaging.Center, imaging.Lanczos)
	err = imaging.Save(resizedImage, thumbImgPath.String())
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Path = thumbImgPath.String()

	json, err := json.Marshal(rcp)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
