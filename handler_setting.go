package main

import (
	"context"
	"net/http"
	"strconv"
)

// handleAdminSetting 함수는 admin setting을 수정하는 페이지이다.
func handleAdminSetting(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel != AdminAccessLevel { // Admin
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
		User
		Projectlist []string
		SearchOption
		Setting Setting
	}
	rcp := recipe{}
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
	rcp.Setting, err = GetAdminSettingV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = TEMPLATES.ExecuteTemplate(w, "adminsetting", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAdminSettingSubmint 함수는 adminsetting을 업데이트 한다.
func handleAdminSettingSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel != AdminAccessLevel {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	s := Setting{}
	s.ID = "admin"
	s.AppName = r.FormValue("AppName")
	s.EmailDNS = r.FormValue("EmailDNS")
	s.RootPath = r.FormValue("RootPath")
	s.ScanPlateUploadPath = r.FormValue("ScanPlateUploadPath")
	s.ProjectPath = r.FormValue("ProjectPath")
	s.ProjectPathPermission = r.FormValue("ProjectPathPermission")
	s.ProjectPathUID = r.FormValue("ProjectPathUID")
	s.ProjectPathGID = r.FormValue("ProjectPathGID")
	s.ShotRootPath = r.FormValue("ShotRootPath")
	s.ShotRootPathPermission = r.FormValue("ShotRootPathPermission")
	s.ShotRootPathUID = r.FormValue("ShotRootPathUID")
	s.ShotRootPathGID = r.FormValue("ShotRootPathGID")
	s.SeqPath = r.FormValue("SeqPath")
	s.SeqPathPermission = r.FormValue("SeqPathPermission")
	s.SeqPathUID = r.FormValue("SeqPathUID")
	s.SeqPathGID = r.FormValue("SeqPathGID")
	s.ShotPath = r.FormValue("ShotPath")
	s.ShotPathPermission = r.FormValue("ShotPathPermission")
	s.ShotPathUID = r.FormValue("ShotPathUID")
	s.ShotPathGID = r.FormValue("ShotPathGID")
	s.AssetRootPath = r.FormValue("AssetRootPath")
	s.AssetRootPathPermission = r.FormValue("AssetRootPathPermission")
	s.AssetRootPathUID = r.FormValue("AssetRootPathUID")
	s.AssetRootPathGID = r.FormValue("AssetRootPathGID")
	s.AssetTypePath = r.FormValue("AssetTypePath")
	s.AssetTypePathPermission = r.FormValue("AssetTypePathPermission")
	s.AssetTypePathUID = r.FormValue("AssetTypePathUID")
	s.AssetTypePathGID = r.FormValue("AssetTypePathGID")
	s.AssetPath = r.FormValue("AssetPath")
	s.AssetPathPermission = r.FormValue("AssetPathPermission")
	s.AssetPathUID = r.FormValue("AssetPathUID")
	s.AssetPathGID = r.FormValue("AssetPathGID")
	s.ThumbnailRootPath = r.FormValue("ThumbnailRootPath")
	s.ThumbnailRootPathPermission = r.FormValue("ThumbnailRootPathPermission")
	s.ThumbnailRootPathUID = r.FormValue("ThumbnailRootPathUID")
	s.ThumbnailRootPathGID = r.FormValue("ThumbnailRootPathGID")
	s.ThumbnailImagePath = r.FormValue("ThumbnailImagePath")
	s.ThumbnailImagePathPermission = r.FormValue("ThumbnailImagePathPermission")
	s.ThumbnailImagePathUID = r.FormValue("ThumbnailImagePathUID")
	s.ThumbnailImagePathGID = r.FormValue("ThumbnailImagePathGID")
	s.ThumbnailMovPath = r.FormValue("ThumbnailMovPath")
	s.ThumbnailMovPathPermission = r.FormValue("ThumbnailMovPathPermission")
	s.ThumbnailMovPathUID = r.FormValue("ThumbnailMovPathUID")
	s.ThumbnailMovPathGID = r.FormValue("ThumbnailMovPathGID")
	s.PlatePath = r.FormValue("PlatePath")
	s.PlatePathPermission = r.FormValue("PlatePathPermission")
	s.PlatePathUID = r.FormValue("PlatePathUID")
	s.PlatePathGID = r.FormValue("PlatePathGID")
	s.ReviewDataPath = r.FormValue("ReviewDataPath")
	s.ReviewDataPathPermission = r.FormValue("ReviewDataPathPermission")
	s.ReviewDataPathUID = r.FormValue("ReviewDataPathUID")
	s.ReviewDataPathGID = r.FormValue("ReviewDataPathGID")
	s.ReviewUploadPath = r.FormValue("ReviewUploadPath")
	s.ReviewUploadPathPermission = r.FormValue("ReviewUploadPathPermission")
	s.ReviewUploadPathUID = r.FormValue("ReviewUploadPathUID")
	s.ReviewUploadPathGID = r.FormValue("ReviewUploadPathGID")
	s.DirectUploadPath = r.FormValue("DirectUploadPath")
	s.DirectUploadPathPermission = r.FormValue("DirectUploadPathPermission")
	s.DirectUploadPathUID = r.FormValue("DirectUploadPathUID")
	s.DirectUploadPathGID = r.FormValue("DirectUploadPathGID")
	s.InitPassword = r.FormValue("InitPassword")
	s.OCIOConfig = r.FormValue("OCIOConfig")
	s.FFmpeg = r.FormValue("FFmpeg")
	s.FFprobe = r.FormValue("FFprobe")
	threads, err := strconv.Atoi(r.FormValue("FFmpegThreads"))
	if err != nil {
		threads = 1
	}
	if threads == 0 {
		threads = 1 // 최소한 1개의 CPU는 셋팅되어야 한다.
	}
	s.FFmpegThreads = threads
	s.OpenImageIO = r.FormValue("OpenImageIO")
	s.Iinfo = r.FormValue("Iinfo")
	s.Curl = r.FormValue("Curl")
	s.SlateFontPath = r.FormValue("SlateFontPath")
	s.RVPath = r.FormValue("RVPath")
	s.Protocol = r.FormValue("Protocol")
	s.WFS = r.FormValue("WFS")
	ratio, err := strconv.ParseFloat(r.FormValue("DefaultScaleRatioOfUndistortionPlate"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.DefaultScaleRatioOfUndistortionPlate = ratio
	itemNumberOfPage, err := strconv.Atoi(r.FormValue("ItemNumberOfPage"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.ItemNumberOfPage = itemNumberOfPage
	multipartFormBufferSize, err := strconv.Atoi(r.FormValue("MultipartFormBufferSize"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.MultipartFormBufferSize = multipartFormBufferSize
	thumbnailImageWidth, err := strconv.Atoi(r.FormValue("ThumbnailImageWidth"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if thumbnailImageWidth == 0 {
		thumbnailImageWidth = 410
	}
	s.ThumbnailImageWidth = thumbnailImageWidth
	thumbnailImageHeight, err := strconv.Atoi(r.FormValue("ThumbnailImageHeight"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if thumbnailImageHeight == 0 {
		thumbnailImageHeight = 222
	}
	s.ThumbnailImageHeight = thumbnailImageHeight
	productionStartFrame, err := strconv.Atoi(r.FormValue("ProductionStartFrame"))
	if err != nil {
		productionStartFrame = 1
	}
	s.ProductionStartFrame = productionStartFrame
	productionPaddingVersionNumber, err := strconv.Atoi(r.FormValue("ProductionPaddingVersionNumber"))
	if err != nil {
		productionPaddingVersionNumber = 2
	}
	s.ProductionPaddingVersionNumber = productionPaddingVersionNumber
	s.EnableEndpoint = str2bool(r.FormValue("EnableEndpoint"))
	s.EnableDirectupload = str2bool(r.FormValue("EnableDirectupload"))
	s.EnableDirectuploadWithDate = str2bool(r.FormValue("EnableDirectuploadWithDate"))
	s.EnableDirectuploadWithProject = str2bool(r.FormValue("EnableDirectuploadWithProject"))
	s.EnableDirectuploadWithCompanyID = str2bool(r.FormValue("EnableDirectuploadWithCompanyID"))
	s.FullcalendarSchedulerLicenseKey = r.FormValue("FullcalendarSchedulerLicenseKey")
	s.AudioCodec = r.FormValue("audiocodec")
	// DB에 값을 저장합니다.
	err = SetAdminSettingV2(client, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	err = TEMPLATES.ExecuteTemplate(w, "setadminsetting", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSetAdminSetting 함수는 admin setting을 수정하는 페이지이다.
func handleSetAdminSetting(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel != 10 { // Admin
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	type recipe struct {
		User        User
		Projectlist []string
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
	rcp.Projectlist, err = ProjectlistV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "update_adminsetting", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
