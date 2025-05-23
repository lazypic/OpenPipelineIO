package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"go.mongodb.org/mongo-driver/bson"
)

// MaxFileSize 사이즈는 웹에서 전송할 수 있는 최대 사이즈를 2기가로 제한한다.(인트라넷)
const MaxFileSize = 2000 * 1000 * 1024

// LoadTemplates 함수는 템플릿을 로딩합니다.
func LoadTemplates() (*template.Template, error) {
	t := template.New("").Funcs(funcMap)
	t, err := vfstemplate.ParseGlob(assets, t, "/template/*.html")
	return t, err
}

// 템플릿 함수를 로딩합니다.
var funcMap = template.FuncMap{
	"safeURL": func(rawURL string) template.URL {
		// URL을 안전하게 처리합니다. 에러 처리는 예제를 단순화하기 위해 생략됩니다.
		// 실제 코드에서는 url.Parse 후에 에러를 확인하고 적절히 처리해야 합니다.
		parsedURL, _ := url.Parse(rawURL)
		return template.URL(parsedURL.String())
	},
	"AddProductionStartFrame":      AddProductionStartFrame,
	"title":                        strings.Title,
	"Split":                        strings.Split,
	"Join":                         strings.Join,
	"Parentpath":                   filepath.Dir,
	"projectStatus2color":          projectStatus2color,
	"name2seq":                     name2seq,
	"note2body":                    note2body,
	"pmnote2body":                  pmnote2body,
	"GetPath":                      GetPath,
	"ReverseStringSlice":           ReverseStringSlice,
	"ReverseCommentSlice":          ReverseCommentSlice,
	"SortByCreatetimeForPublishes": SortByCreatetimeForPublishes,
	"CutStringSlice":               CutStringSlice,
	"CutCommentSlice":              CutCommentSlice,
	"ToShortTime":                  ToShortTime,
	"ToNormalTime":                 ToNormalTime,
	"List2str":                     List2str,
	"CheckDate":                    CheckDate,
	"CheckUpdate":                  CheckUpdate,
	"CheckDdline":                  CheckDdline,
	"CheckDdlinev2":                CheckDdlinev2,
	"ToHumantime":                  ToHumantime,
	"Framecal":                     Framecal,
	"Add":                          Add,
	"Minus":                        Minus,
	"Scanname2RollMedia":           Scanname2RollMedia,
	"AddTagColon":                  AddTagColon, //Hashtag2tag,
	"Username2Elements":            Username2Elements,
	"RemovePath":                   RemovePath,
	"ShortPhoneNum":                ShortPhoneNum,
	"TaskStatus":                   TaskStatus,
	"TaskUser":                     TaskUser,
	"TaskDate":                     TaskDate,
	"TaskPredate":                  TaskPredate,
	"ProductionVersionFormat":      ProductionVersionFormat,
	"Protocol":                     Protocol,
	"RmProtocol":                   RmProtocol,
	"ProtocolTarget":               ProtocolTarget,
	"userInfo":                     userInfo,
	"onlyID":                       onlyID,
	"mapToSlice":                   mapToSlice,
	"hasStatus":                    hasStatus,
	"GenPageNums":                  GenPageNums,
	"floatToString":                floatToString,
}

func errorHandler(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "NotFound 404")
	}
}

// 도움말 페이지 입니다.
func handleHelp(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	u, err := getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := LoadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type recipe struct {
		User User
		SearchOption
		Sha1ver   string
		BuildTime string
		Status    []Status
		DBIP      string
		DBVer     string
		ServerIP  string
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	rcp.Sha1ver = SHA1VER
	rcp.BuildTime = BUILDTIME
	rcp.DBIP = *flagDBIP

	var result bson.M
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "buildInfo", Value: 1}}).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 버전 정보 출력
	rcp.DBVer = fmt.Sprintf("%s", result["version"])

	err = rcp.SearchOption.LoadCookieV2(client, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.Status, err = AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ip, err := serviceIP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.ServerIP = ip
	err = t.ExecuteTemplate(w, "help", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	type recipe struct {
		IP         string
		Web        bool
		DB         bool
		MountPoint bool
		All        bool
	}
	rcp := recipe{}
	// IP구하기
	ip, err := serviceIP()
	if err != nil {
		rcp.IP = ""
	}
	rcp.IP = ip
	// DB 체크

	client, err := initMongoClient()
	if err != nil {
		rcp.DB = false
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	rcp.DB = true

	// 웹서버 체크
	rcp.Web = true

	// Mount Point 경로 존재하는지 체크
	_, err = os.Stat(CachedAdminSetting.RootPath)
	if os.IsNotExist(err) {
		rcp.MountPoint = false
	}
	rcp.MountPoint = true

	// 모든 요소가 true인지 체크
	isIPexist := false
	if rcp.IP != "" {
		isIPexist = true
	}
	rcp.All = isIPexist && rcp.Web && rcp.DB && rcp.MountPoint

	// json 으로 결과 전송
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 전송되는 컨텐츠의 캐쉬 수명을 설정하는 핸들러입니다.
func maxAgeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", *flagThumbnailAge))
		h.ServeHTTP(w, r)
	})
}

func helpMethodOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
}

// webserver함수는 웹서버의 URL을 선언하는 함수입니다.
func webserver(port string) {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(assets)))
	r.PathPrefix("/thumbnail/").Handler(maxAgeHandler(http.StripPrefix("/thumbnail/", http.FileServer(http.Dir(CachedAdminSetting.ThumbnailRootPath)))))
	r.PathPrefix("/captcha/").Handler(captcha.Server(captcha.StdWidth, captcha.StdHeight)) // Captcha

	// DirectUpload
	r.HandleFunc("/directupload", handleDirectupload).Methods("GET")
	r.HandleFunc("/ws/directupload", directUploadHandler)
	r.HandleFunc("/ws/directuploadprogress", directUploadProgressHandler)

	// ScanPlate
	r.HandleFunc("/scanplate", handleScanPlate)
	r.HandleFunc("/process", handleProcess)
	r.HandleFunc("/upload-scanplate", handleUploadScanPlate)
	r.HandleFunc("/api/scanplate", postHandleAPIScanPlate).Methods("POST")
	r.HandleFunc("/api/scanplate/{id}", deleteHandleAPIScanPlate).Methods("DELETE")
	r.HandleFunc("/api/scanplate/{id}", putHandleAPIScanPlate).Methods("PATCH")
	r.HandleFunc("/api/scanplates", handleAPIScanPlates).Methods("POST")
	r.HandleFunc("/api/scanplatetemp", deleteScanPlateTemp).Methods("DELETE")
	r.HandleFunc("/api/searchfootages", handleAPISearchFootages).Methods("POST")
	r.HandleFunc("/api/ociocolorspace", handleAPIOcioColorspace).Methods("GET")
	r.HandleFunc("/api/protocol", handleAPIProtocol).Methods("GET")

	// Item
	r.HandleFunc("/", handleIndex)
	r.HandleFunc("/inputmode", handleInputMode)
	r.HandleFunc("/searchsubmit", handleSearchSubmit).Methods("POST")
	r.HandleFunc("/searchsubmitv2", handleSearchSubmitV2).Methods("POST")
	r.HandleFunc("/help", handleHelp)
	r.HandleFunc("/setellite", handleSetellite)
	r.HandleFunc("/uploadsetellite", handleUploadSetellite)
	r.HandleFunc("/addshot", handleAddShot)
	r.HandleFunc("/addshot_submit", handleAddShotSubmit)
	r.HandleFunc("/addasset", handleAddAsset)
	r.HandleFunc("/addasset_submit", handleAddAssetSubmit)
	r.HandleFunc("/detail", handleItemDetail).Methods("GET")

	// Review
	r.HandleFunc("/daily-review-status", handleDailyReviewStatus)
	r.HandleFunc("/reviewstatus", handleReviewStatus)
	r.HandleFunc("/reviewdata", handleReviewData)
	r.HandleFunc("/reviewdrawingdata", handleReviewDrawingData)
	r.HandleFunc("/review-status-submit", handleReviewStatusSubmit)
	r.HandleFunc("/upload-reviewfile", handleUploadReviewFile)

	// Project
	r.HandleFunc("/projectinfo", handleProjectinfo)
	r.HandleFunc("/addproject", handleAddProject)
	r.HandleFunc("/addproject_submit", handleAddProjectSubmit)
	r.HandleFunc("/editproject", handleEditProject)
	r.HandleFunc("/editproject_submit", handleEditProjectSubmit)
	r.HandleFunc("/rmproject", handleRmProject)
	r.HandleFunc("/rmproject_submit", handleRmProjectSubmit)
	r.HandleFunc("/noonproject", handleNoOnProject)

	// User
	r.HandleFunc("/signup", handleSignup)
	r.HandleFunc("/signup_submit", handleSignupSubmit)
	r.HandleFunc("/signin", handleSignin)
	r.HandleFunc("/signin_submit", handleSigninSubmit)
	r.HandleFunc("/signin_success", handleSigninSuccess)
	r.HandleFunc("/signout", handleSignout)
	r.HandleFunc("/user", handleUser)
	r.HandleFunc("/users", handleUsers)
	r.HandleFunc("/updatepassword", handleUpdatePassword)
	r.HandleFunc("/updatepassword_submit", handleUpdatePasswordSubmit)
	r.HandleFunc("/edituser", handleEditUser)
	r.HandleFunc("/edituser-submit", handleEditUserSubmit)
	r.HandleFunc("/replacetag", handleReplaceTag)
	r.HandleFunc("/replacetag_submit", handleReplaceTagSubmit).Methods("POST")
	r.HandleFunc("/invalidaccess", handleInvalidAccess)
	r.HandleFunc("/invalidpass", handleInvalidPass)
	r.HandleFunc("/nouser", handleNoUser)

	// Partner
	r.HandleFunc("/partners", handlePartners)
	r.HandleFunc("/projectsforpartner", handleProjectsForPartner)

	// Endpoint
	r.HandleFunc("/endpoints", handleEndpoints)

	// Scenario
	r.HandleFunc("/scenarios", handleScenarios)
	r.HandleFunc("/import-scenario-pdf", handleImportScenarioPdf)

	// Admin Setting
	r.HandleFunc("/adminsetting", handleAdminSetting)
	r.HandleFunc("/adminsetting_submit", handleAdminSettingSubmit)
	r.HandleFunc("/setadminsetting", handleSetAdminSetting)

	// Organization
	r.HandleFunc("/divisions", handleDivisions)
	r.HandleFunc("/departments", handleDepartments)
	r.HandleFunc("/teams", handleTeams)
	r.HandleFunc("/roles", handleRoles)
	r.HandleFunc("/positions", handlePositions)
	r.HandleFunc("/adddivision", handleAddOrganization)
	r.HandleFunc("/editdivision", handleEditDivision)
	r.HandleFunc("/editdivisionsubmit", handleEditDivisionSubmit)
	r.HandleFunc("/adddepartment", handleAddOrganization)
	r.HandleFunc("/editdepartment", handleEditDepartment)
	r.HandleFunc("/editdepartmentsubmit", handleEditDepartmentSubmit)
	r.HandleFunc("/addteam", handleAddOrganization)
	r.HandleFunc("/editteam", handleEditTeam)
	r.HandleFunc("/editteamsubmit", handleEditTeamSubmit)
	r.HandleFunc("/addrole", handleAddOrganization)
	r.HandleFunc("/editrole", handleEditRole)
	r.HandleFunc("/editrolesubmit", handleEditRoleSubmit)
	r.HandleFunc("/addposition", handleAddOrganization)
	r.HandleFunc("/editposition", handleEditPosition)
	r.HandleFunc("/editpositionsubmit", handleEditPositionSubmit)
	r.HandleFunc("/adddivisionsubmit", handleAddDivisionSubmit)
	r.HandleFunc("/adddepartmentsubmit", handleAddDepartmentSubmit)
	r.HandleFunc("/addteamsubmit", handleAddTeamSubmit)
	r.HandleFunc("/addrolesubmit", handleAddRoleSubmit)
	r.HandleFunc("/addpositionsubmit", handleAddPositionSubmit)
	r.HandleFunc("/rmorganization", handleRmOrganization)
	r.HandleFunc("/rmorganization-submit", handleRmOrganizationSubmit)

	// Import, Export: Excel, Json, Csv, Dump
	r.HandleFunc("/importexcel", handleImportExcel)
	r.HandleFunc("/importjson", handleImportJSON)
	r.HandleFunc("/exportexcel", handleExportExcel)
	r.HandleFunc("/exportjson", handleExportJSON)
	r.HandleFunc("/reportexcel", handleReportExcel).Methods("GET")
	r.HandleFunc("/reportjson", handleReportJSON).Methods("GET")
	r.HandleFunc("/excel-submit", handleExcelSubmit).Methods("POST")
	r.HandleFunc("/json-submit", handleJSONSubmit).Methods("POST")
	r.HandleFunc("/exportexcel-submit", handleExportExcelSubmit).Methods("POST")
	r.HandleFunc("/exportjson-submit", handleExportJSONSubmit).Methods("POST")
	r.HandleFunc("/upload-excel", handleUploadExcel)
	r.HandleFunc("/upload-json", handleUploadJSON)
	r.HandleFunc("/download-excel-template", handleDownloadExcelTemplate).Methods("GET")
	r.HandleFunc("/download-excel-file", handleDownloadExcelFile).Methods("GET")
	r.HandleFunc("/download-json-file", handleDownloadJSONFile).Methods("GET")
	r.HandleFunc("/download-csv-file", handleDownloadCsvFile).Methods("GET")

	// Task
	r.HandleFunc("/tasksettings", handleTasksettings)
	r.HandleFunc("/addtasksetting", handleAddTasksetting)
	r.HandleFunc("/rmtasksetting", handleRmTasksetting)
	r.HandleFunc("/edittasksetting", handleEditTasksetting)
	r.HandleFunc("/addtasksetting-submit", handleAddTasksettingSubmit)
	r.HandleFunc("/rmtasksetting-submit", handleRmTasksettingSubmit)
	r.HandleFunc("/edittasksetting-submit", handleEditTasksettingSubmit)

	// Status
	r.HandleFunc("/status", handleStatus)
	r.HandleFunc("/addstatus", handleAddStatus)
	r.HandleFunc("/addstatus-submit", handleAddStatusSubmit)
	r.HandleFunc("/editstatus", handleEditStatus)
	r.HandleFunc("/editstatus-submit", handleEditStatusSubmit)
	r.HandleFunc("/rmstatus", handleRmStatus)
	r.HandleFunc("/rmstatus-submit", handleRmStatusSubmit)

	// Publish Key
	r.HandleFunc("/publishkey", handlePublishKey)
	r.HandleFunc("/addpublishkey", handleAddPublishKey)
	r.HandleFunc("/addpublishkey-submit", handleAddPublishKeySubmit)
	r.HandleFunc("/editpublishkey", handleEditPublishKey)
	r.HandleFunc("/editpublishkey-submit", handleEditPublishKeySubmit)
	r.HandleFunc("/rmpublishkey", handleRmPublishKey)
	r.HandleFunc("/rmpublishkey-submit", handleRmPublishKeySubmit)

	// Error
	r.HandleFunc("/error-captcha", handleErrorCaptcha)

	//Health & Help
	r.HandleFunc("/health", handleHealth)
	r.HandleFunc("/api/statusinfo", handleAPIStatusInfo).Methods(http.MethodGet, http.MethodOptions)

	// Statistics
	r.HandleFunc("/statistics", handleStatistics).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/statistics/statusperproject", handleStatisticsStatusPerProject).Methods(http.MethodGet, http.MethodOptions)

	// Statistics API
	r.HandleFunc("/api/statistics/projectnum", handleAPIStatisticsProjectnum).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/statistics/deadlinenum", handleAPIStatisticsDeadlineNum).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/statistics/needdeadlinenum", handleAPIStatisticsNeedDeadlineNum).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/statistics/shottype", handleAPIStatisticsShottype).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/statistics/itemtype", handleAPIStatisticsItemtype).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/shot", handleAPI1StatisticsShot).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api2/statistics/shot", handleAPI2StatisticsShot).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/asset", handleAPI1StatisticsAsset).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api2/statistics/asset", handleAPI2StatisticsAsset).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/task", handleAPI1StatisticsTask).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api2/statistics/task", handleAPI2StatisticsTask).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/pipelinestep", handleAPI1StatisticsPipelinestep).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api2/statistics/pipelinestep", handleAPI2StatisticsPipelinestep).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/tag", handleAPI1StatisticsTag).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api2/statistics/tag", handleAPI2StatisticsTag).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api1/statistics/user", handleAPI1StatisticsUser).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api2/statistics/user", handleAPI2StatisticsUser).Methods(http.MethodGet, http.MethodOptions)

	// restAPI Project
	r.HandleFunc("/api/project", handleAPIProject).Methods("GET")
	r.HandleFunc("/api2/projects", handleAPI2Projects).Methods("GET")
	r.HandleFunc("/api/addproject", handleAPIAddproject).Methods("POST")
	r.HandleFunc("/api/projecttags", handleAPIProjectTags).Methods("GET")
	r.HandleFunc("/api/projectassettags", handleAPIProjectAssetTags).Methods("GET")

	// restAPI Onset(Setellite)
	r.HandleFunc("/api/setellite", handleAPISetelliteItems)
	r.HandleFunc("/api/setellitesearch", handleAPISetelliteSearch)

	// restAPI Item
	r.HandleFunc("/api/timeinfo", handleAPITimeinfo).Methods("POST")
	r.HandleFunc("/api2/item", handleAPI2GetItem).Methods("GET")
	r.HandleFunc("/api/rmitemid", handleAPIRmItemID).Methods("POST")
	r.HandleFunc("/api3/items", handleAPI3Items)
	r.HandleFunc("/api/seqs", handleAPISeqs)
	r.HandleFunc("/api/allshots", handleAPIAllShots)
	r.HandleFunc("/api2/shots", handleAPI2Shots)
	r.HandleFunc("/api/shot", handleAPIGetShot).Methods("GET")
	r.HandleFunc("/api/item", postHandleAPIItem).Methods("POST")
	r.HandleFunc("/api/asset", handleAPIGetAsset).Methods("GET")
	r.HandleFunc("/api/assets", handleAPIAssets)
	r.HandleFunc("/api/setplatesize", handleAPISetPlateSize).Methods("POST")
	r.HandleFunc("/api/setundistortionsize", handleAPISetUnDistortionSize).Methods("POST")
	r.HandleFunc("/api2/setrendersize", handleAPI2SetRenderSize).Methods("POST")
	r.HandleFunc("/api/setoverscanratio", handleAPISetOverscanRatio).Methods("POST")
	r.HandleFunc("/api/setcamerapubpath", handleAPISetCameraPubPath).Methods("POST")
	r.HandleFunc("/api/setcamerapubtask", handleAPISetCameraPubTask).Methods("POST")
	r.HandleFunc("/api/setcameralensmm", handleAPISetCameraLensmm).Methods("POST")
	r.HandleFunc("/api/setcameraprojection", handleAPISetCameraProjection).Methods("POST")
	r.HandleFunc("/api/setseq", handleAPISetSeq).Methods("POST")
	r.HandleFunc("/api/setscene", handleAPISetScene).Methods("POST")
	r.HandleFunc("/api/setcut", handleAPISetCut).Methods("POST")
	r.HandleFunc("/api/setepisode", handleAPISetEpisode).Methods("POST")
	r.HandleFunc("/api/setplatepath", handleAPISetPlatePath).Methods("POST")
	r.HandleFunc("/api/setthummov", handleAPI2SetThummov).Methods("POST")
	r.HandleFunc("/api/setbeforemov", handleAPISetBeforemov).Methods("POST")
	r.HandleFunc("/api/setaftermov", handleAPISetAftermov).Methods("POST")
	r.HandleFunc("/api/seteditmov", handleAPISetEditmov).Methods("POST")
	r.HandleFunc("/api2/settaskstatus", handleAPI2SetTaskStatus).Methods("POST")
	r.HandleFunc("/api/taskstatusnum", handleAPITaskStatusNum)
	r.HandleFunc("/api/taskanduserstatusnum", handleAPITaskAndUserStatusNum)
	r.HandleFunc("/api/userstatusnum", handleAPIUserStatusNum)
	r.HandleFunc("/api/statusnum", handleAPIStatusNum)
	r.HandleFunc("/api/addtask", handleAPIAddTask).Methods("POST")
	r.HandleFunc("/api/rmtask", handleAPIRmTask).Methods("POST")
	r.HandleFunc("/api/settaskuser", handleAPISetTaskUser).Methods("POST")
	r.HandleFunc("/api/settaskusercomment", handleAPISetTaskUserComment).Methods("POST")
	r.HandleFunc("/api/setplatein", handleAPISetPlateIn).Methods("POST")
	r.HandleFunc("/api/setplateout", handleAPISetPlateOut).Methods("POST")
	r.HandleFunc("/api/setjustin", handleAPISetJustIn).Methods("POST")
	r.HandleFunc("/api/setjustout", handleAPISetJustOut).Methods("POST")
	r.HandleFunc("/api/setscanin", handleAPISetScanIn).Methods("POST")
	r.HandleFunc("/api/setscanout", handleAPISetScanOut).Methods("POST")
	r.HandleFunc("/api/setscanframe", handleAPISetScanFrame).Methods("POST")
	r.HandleFunc("/api/sethandlein", handleAPISetHandleIn).Methods("POST")
	r.HandleFunc("/api/sethandleout", handleAPISetHandleOut).Methods("POST")
	r.HandleFunc("/api/setshottype", handleAPISetShotType).Methods("POST")
	r.HandleFunc("/api/setusetype", handleAPISetUseType).Methods("POST")
	r.HandleFunc("/api/setassettype", handleAPISetAssetType).Methods("POST")
	r.HandleFunc("/api/setoutputname", handleAPISetOutputName)
	r.HandleFunc("/api2/setrnum", handleAPI2SetRnum).Methods("POST")
	r.HandleFunc("/api/setdeadline2d", handleAPISetDeadline2D).Methods("POST")
	r.HandleFunc("/api/setdeadline3d", handleAPISetDeadline3D).Methods("POST")
	r.HandleFunc("/api/setscantimecodein", handleAPISetScanTimecodeIn).Methods("POST")
	r.HandleFunc("/api/setscantimecodeout", handleAPISetScanTimecodeOut).Methods("POST")
	r.HandleFunc("/api/setjusttimecodein", handleAPISetJustTimecodeIn).Methods("POST")
	r.HandleFunc("/api/setjusttimecodeout", handleAPISetJustTimecodeOut).Methods("POST")
	r.HandleFunc("/api/setfinver", handleAPISetFinver).Methods("POST")
	r.HandleFunc("/api/setfindate", handleAPISetFindate)
	r.HandleFunc("/api/addtag", handleAPIAddTag).Methods("POST")
	r.HandleFunc("/api/addassettag", handleAPIAddAssetTag).Methods("POST")
	r.HandleFunc("/api/renametag", handleAPIRenameTag).Methods("POST")
	r.HandleFunc("/api/rmtag", handleAPIRmTag).Methods("POST")
	r.HandleFunc("/api/rmassettag", handleAPIRmAssetTag).Methods("POST")
	r.HandleFunc("/api/setnote", handleAPISetNote).Methods("POST")
	r.HandleFunc("/api/addcomment", handleAPIAddComment).Methods("POST")
	r.HandleFunc("/api/editcomment", handleAPIEditComment).Methods("POST")
	r.HandleFunc("/api/rmcomment", handleAPIRmComment).Methods("POST")
	r.HandleFunc("/api/addsource", handleAPIAddSource).Methods("POST")
	r.HandleFunc("/api/rmsource", handleAPIRmSource).Methods("POST")
	r.HandleFunc("/api/addreference", handleAPIAddReference).Methods("POST")
	r.HandleFunc("/api/rmreference", handleAPIRmReference).Methods("POST")
	r.HandleFunc("/api/search", handleAPISearch)
	r.HandleFunc("/api/deadline2d", handleAPIDeadline2D)
	r.HandleFunc("/api/deadline3d", handleAPIDeadline3D)
	r.HandleFunc("/api2/settaskmov", handleAPI2SetTaskMov).Methods("POST")
	r.HandleFunc("/api/settaskusernote", handleAPISetTaskUserNote).Methods("POST")
	r.HandleFunc("/api/setretimeplate", handleAPISetRetimePlate).Methods("POST")
	r.HandleFunc("/api/setobjectid", handleAPISetObjectID)
	r.HandleFunc("/api/setscanname", handleAPISetScanname).Methods("POST")
	r.HandleFunc("/api/settaskdate", handleAPISetTaskDate).Methods("POST")
	r.HandleFunc("/api/taskduration", handleAPISetTaskDuration).Methods("POST")
	r.HandleFunc("/api/settaskexpectday", handleAPISetTaskExpectDay)
	r.HandleFunc("/api/settaskresultday", handleAPISetTaskResultDay)
	r.HandleFunc("/api/task", handleAPITask).Methods("POST")
	r.HandleFunc("/api/settaskstart", handleAPISetTaskStartdate).Methods("POST")
	r.HandleFunc("/api/settaskend", handleAPISetTaskEnd).Methods("POST")
	r.HandleFunc("/api/shottype", handleAPIShottype).Methods("POST")
	r.HandleFunc("/api/setcrowdasset", handleAPISetCrowdAsset)
	r.HandleFunc("/api/mailinfo", handleAPIMailInfo).Methods("POST")
	r.HandleFunc("/api/usetypes", handleAPIUseTypes).Methods("GET")
	r.HandleFunc("/api/addpublish", handleAPIAddTaskPublish).Methods("POST")
	r.HandleFunc("/api/setpublishstatus", handleAPISetTaskPublishStatus).Methods("POST")
	r.HandleFunc("/api/rmpublish", handleAPIRmTaskPublish).Methods("POST")
	r.HandleFunc("/api/rmpublishkey", handleAPIRmTaskPublishKey).Methods("POST")
	r.HandleFunc("/api/uploadthumbnail", handleAPIUploadThumbnail).Methods("POST")

	// restAPI USER
	r.HandleFunc("/api2/user", handleAPI2User)
	r.HandleFunc("/api/users", handleAPISearchUser).Methods("GET")
	r.HandleFunc("/api/setleaveuser", handleAPISetLeaveUser).Methods("POST")
	r.HandleFunc("/api/autocompliteusers", handleAPIAutoCompliteUsers).Methods("GET")
	r.HandleFunc("/api/initpassword", handleAPIInitPassword).Methods("POST")
	r.HandleFunc("/api/ansiblehosts", handleAPIAnsibleHosts).Methods("GET")

	// restAPI Organization
	r.HandleFunc("/api/teams", handleAPIAllTeams).Methods("GET")

	// restAPI Tasksetting
	r.HandleFunc("/api/tasksetting", handleAPITasksetting)
	r.HandleFunc("/api/shottasksetting", handleAPIShotTasksetting).Methods("GET")
	r.HandleFunc("/api/assettasksetting", handleAPIAssetTasksetting)
	r.HandleFunc("/api/categorytasksettings", handleAPICategoryTasksettings)

	// restAPI Status
	r.HandleFunc("/api/status", handleAPIStatus)
	r.HandleFunc("/api/addstatus", handleAPIAddStatus)

	// restAPI PublishKey
	r.HandleFunc("/api/publishkeys", handleAPIPublishKeys).Methods("GET")
	r.HandleFunc("/api/getpublish", handleAPIGetPublish)

	// restAPI Review
	r.HandleFunc("/api/addreview", handleAPIAddReview).Methods("POST")
	r.HandleFunc("/api/review", handleAPIReview).Methods("POST")
	r.HandleFunc("/api/searchreview", handleAPISearchReview).Methods("POST")
	r.HandleFunc("/api/setreviewitemstatus", handleAPISetReviewItemStatus).Methods("POST")
	r.HandleFunc("/api/addreviewstatusmodecomment", handleAPIAddReviewStatusModeComment).Methods("POST")
	r.HandleFunc("/api/editreviewcomment", handleAPIEditReviewComment).Methods("POST")
	r.HandleFunc("/api/rmreviewcomment", handleAPIRmReviewComment).Methods("POST")
	r.HandleFunc("/api/rmreview", handleAPIRmReview)
	r.HandleFunc("/api/setreviewproject", handleAPISetReviewProject)
	r.HandleFunc("/api/setreviewtask", handleAPISetReviewTask)
	r.HandleFunc("/api/setreviewname", handleAPISetReviewName)
	r.HandleFunc("/api/setreviewpath", handleAPISetReviewPath)
	r.HandleFunc("/api/setreviewcreatetime", handleAPISetReviewCreatetime)
	r.HandleFunc("/api/setreviewmainversion", handleAPISetReviewMainVersion)
	r.HandleFunc("/api/setreviewsubversion", handleAPISetReviewSubVersion)
	r.HandleFunc("/api/setreviewfps", handleAPISetReviewFps)
	r.HandleFunc("/api/setreviewdescription", handleAPISetReviewDescription)
	r.HandleFunc("/api/setreviewcamerainfo", handleAPISetReviewCameraInfo)
	r.HandleFunc("/api/setreviewprocessstatus", handleAPISetReviewProcessStatus)
	r.HandleFunc("/api/uploadreviewdrawing", handleAPIUploadReviewDrawing)
	r.HandleFunc("/api/rmreviewdrawing", handleAPIRmReviewDrawing)
	r.HandleFunc("/api/reviewdrawingframe", handleAPIReviewDrawingFrame)
	r.HandleFunc("/api/setreviewagainforwaitstatustoday", handleAPISetReviewAgainForWaitStatusToday)
	r.HandleFunc("/api/reviewoutputdatapath", handleAPIReviewOutputDataPath)

	// REST API PDF
	r.HandleFunc("/api/pdf-to-json", handlerAPIPdfToJson).Methods("POST")

	// REST API Partner
	r.HandleFunc("/api/partner", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/partner", postPartnerHandler).Methods("POST")
	r.HandleFunc("/api/partner/{id}", getPartnerHandler).Methods("GET")
	r.HandleFunc("/api/partner/{id}", putPartnerHandler).Methods("PUT")
	r.HandleFunc("/api/partner/{id}", deletePartnerHandler).Methods("DELETE")
	r.HandleFunc("/api/partners", getPartnersHandler).Methods("GET")
	r.HandleFunc("/api/partnerscodename", getPartnersCodenameHandler).Methods("GET")

	// REST API Partner
	r.HandleFunc("/api/projectforpartner", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/projectforpartner", postProjectForPartnerHandler).Methods("POST")
	r.HandleFunc("/api/projectforpartner/{id}", getProjectForPartnerHandler).Methods("GET")
	r.HandleFunc("/api/projectforpartner/{id}", putProjectForPartnerHandler).Methods("PUT")
	r.HandleFunc("/api/projectforpartner/{id}", deleteProjectForPartnerHandler).Methods("DELETE")

	// REST API Endpoint
	r.HandleFunc("/api/endpoint", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/endpoint", postEndpointHandler).Methods("POST")
	r.HandleFunc("/api/endpoint/{id}", getEndpointHandler).Methods("GET")
	r.HandleFunc("/api/endpoints", getEndpointsHandler).Methods("GET")
	r.HandleFunc("/api/endpoint/{id}", putEndpointHandler).Methods("PUT")
	r.HandleFunc("/api/endpoint/{id}", deleteEndpointHandler).Methods("DELETE")

	// REST API Scenario
	r.HandleFunc("/api/scenario", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/scenario", postScenarioHandler).Methods("POST")
	r.HandleFunc("/api/scenario/{id}", getScenarioHandler).Methods("GET")
	r.HandleFunc("/api/scenario/{id}", putScenarioHandler).Methods("PUT")
	r.HandleFunc("/api/scenario/{id}", deleteScenarioHandler).Methods("DELETE")
	r.HandleFunc("/api/scenarios", postScenariosHandler).Methods("POST")

	r.HandleFunc("/api/ganimage/{id}", postGANImageHandler).Methods("POST")
	r.HandleFunc("/api/ganimages/{id}", getGANImagesHandler).Methods("GET")
	r.HandleFunc("/api/ganimage/{id}", putGANImageHandler).Methods("PUT")
	r.HandleFunc("/api/ganimage/{id}", deleteGANImageHandler).Methods("DELETE")

	// REST API Money
	r.HandleFunc("/api/money", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/money", postMoneyHandler).Methods("POST")
	r.HandleFunc("/api/money/{id}", getMoneyHandler).Methods("GET")
	r.HandleFunc("/api/money/{id}", putMoneyHandler).Methods("PUT")
	r.HandleFunc("/api/money/{id}", deleteMoneyHandler).Methods("DELETE")

	// REST API Moneytype
	r.HandleFunc("/api/moneytype", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/moneytype", postMoneytypeHandler).Methods("POST")
	r.HandleFunc("/api/moneytype/{id}", getMoneytypeHandler).Methods("GET")
	r.HandleFunc("/api/moneytype/{id}", putMoneytypeHandler).Methods("PUT")
	r.HandleFunc("/api/moneytype/{id}", deleteMoneytypeHandler).Methods("DELETE")

	// REST API FullCalendar Event
	r.HandleFunc("/api/fcevent", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/fcevent", postFCEventHandler).Methods("POST")
	r.HandleFunc("/api/fcevent/{id}", getFCEventHandler).Methods("GET")
	r.HandleFunc("/api/fcevent/{id}", putFCEventHandler).Methods("PUT")
	r.HandleFunc("/api/fcevent/{id}", deleteFCEventHandler).Methods("DELETE")

	// REST API FullCalendar Resource
	r.HandleFunc("/api/fcresource", helpMethodOptionsHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions)
	r.HandleFunc("/api/fcresource", postFCResourceHandler).Methods("POST")
	r.HandleFunc("/api/fcresource/{id}", getFCResourceHandler).Methods("GET")
	r.HandleFunc("/api/fcresource/{id}", putFCResourceHandler).Methods("PUT")
	r.HandleFunc("/api/fcresource/{id}", deleteFCResourceHandler).Methods("DELETE")

	r.Use(mux.CORSMethodMiddleware(r))
	http.Handle("/", r)

	if port == ":443" { // https ports
		err := http.ListenAndServeTLS(port, *flagCertFullchanin, *flagCertPrivkey, r)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := http.ListenAndServe(port, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}
