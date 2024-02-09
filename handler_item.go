package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"gopkg.in/mgo.v2"
)

// handleSearchSubmit 함수는 검색창의 옵션을 파싱하고 검색 URI로 리다이렉션 한다. // legacy
func handleSearchSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	project := r.FormValue("Project")
	searchword := r.FormValue("Searchword")
	sortkey := r.FormValue("Sortkey")
	task := r.FormValue("Task")
	truestatus := r.FormValue("truestatus")
	// 아래 코드는 임시로 사용한다.

	redirectURL := fmt.Sprintf(`/inputmode?project=%s&searchword=%s&sortkey=%s&task=%s&truestatus=%s`,
		project,
		searchword,
		sortkey,
		task,
		truestatus,
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// handleSearchSubmitV2 함수는 검색창의 옵션을 파싱하고 검색 URI로 리다이렉션 한다.
func handleSearchSubmitV2(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	project := r.FormValue("Project")
	searchword := r.FormValue("Searchword")
	sortkey := r.FormValue("Sortkey")
	task := r.FormValue("Task")
	// status를 체크할 때 마다 truestatus form에 값이 추가되어야 한다.
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	statuslist, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var truestatus []string
	for _, status := range statuslist {
		if "true" == r.FormValue(status.ID) {
			truestatus = append(truestatus, status.ID)
		}
	}
	redirectURL := fmt.Sprintf(`/inputmode?project=%s&searchword=%s&sortkey=%s&task=%s&truestatus=%s`,
		project,
		searchword,
		sortkey,
		task,
		strings.Join(truestatus, ","),
	)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// handleItemDetail 함수는 아이템 디테일 페이지를 출력한다.
func handleItemDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Get Only", http.StatusMethodNotAllowed)
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
	q := r.URL.Query()
	id := q.Get("id")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		User        User
		Projectlist []string
		Projectinfo Project
		SearchOption
		Item                Item
		TasksettingNames    []string
		TasksettingOrderMap map[string]float64
		Status              []Status
		AllStatusIDs        []string
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rcp.SearchOption.Project != "" {
		rcp.Projectinfo, err = getProject(session, rcp.SearchOption.Project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	tasks, err := AllTaskSettings(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.TasksettingOrderMap = make(map[string]float64)
	for _, t := range tasks {
		rcp.TasksettingOrderMap[t.Name] = t.Order
	}
	rcp.Status, err = AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, status := range rcp.Status {
		rcp.AllStatusIDs = append(rcp.AllStatusIDs, status.ID)
	}
	rcp.TasksettingNames, err = TasksettingNames(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Item, err = getItem(session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "detail", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleIndex 함수는 index 페이지이다.
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // 등록되지 않은 URL은 NotFound 처리한다.
		errorHandler(w, r, http.StatusNotFound)
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
	type recipe struct {
		SearchOption
	}
	rcp := recipe{}
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	err = rcp.SearchOption.LoadCookieV2(client, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 아무 상태도 선택되어있지 않다면 기본 상태설정으로 변경한다. // V2
	if len(rcp.SearchOption.TrueStatus) == 0 {
		err := rcp.SearchOption.setStatusDefaultV2()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 혹시나 프로젝트가 삭제되면 프로젝트가 존재하지 않을 수 있다. DB에 프로젝트가 존재하는지 체크한다.
	plist, err := ProjectlistV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 프로젝트가 하나도 없다면, 프로젝트 생성페이지로 이동한다.
	if len(plist) == 0 {
		http.Redirect(w, r, "/addproject", http.StatusSeeOther)
		return
	}
	// 기존 검색바의 프로젝트 문자열이 DB에 존재하는지 체크한다.
	hasProject := false
	for _, p := range plist {
		if p == rcp.SearchOption.Project {
			hasProject = true
			break
		}
		continue
	}
	// Project가 존재하지 않으면 전체 ProjectList의 첫번째 프로젝트를 가지고 온다.
	if !hasProject {
		rcp.SearchOption.Project = plist[0]
	}

	// 리다이렉션 한다.
	url := fmt.Sprintf("/inputmode?project=%s&sortkey=%s&template=index&task=%s&searchword=%s&truestatus=%s&page=%d",
		rcp.SearchOption.Project,
		rcp.SearchOption.Sortkey,
		rcp.SearchOption.Task,
		rcp.SearchOption.Searchword,
		strings.Join(rcp.SearchOption.TrueStatus, ","),
		rcp.SearchOption.Page,
	)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// handleAddShot 함수는 shot을 추가하는 페이지이다.
func handleAddShot(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		User        User
		Projectlist []string
		SearchOption
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "addShot", rcp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddShotSubmit 함수는 shot을 생성한다.
func handleAddShotSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	project := r.FormValue("Project")
	name := r.FormValue("Name")
	typ := r.FormValue("Type")
	season := r.FormValue("Season")
	episode := r.FormValue("Episode")
	mkdir := str2bool(r.FormValue("Mkdir"))
	genTask := str2bool(r.FormValue("GenTask"))
	setRendersize := str2bool(r.FormValue("SetRendersize"))
	netflixID := r.FormValue("NetflixID")
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}

	names := strings.FieldsFunc(name, f)
	type Shot struct {
		Project string
		Name    string
		Error   string
	}
	var success []Shot
	var fails []Shot
	tasks, err := AllTaskSettings(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	admin, err := GetAdminSetting(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pinfo, err := getProject(session, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	initStatus, err := GetInitStatusID(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, name := range names {
		if name == " " || name == "" { // 사용자가 실수로 여러개의 스페이스를 추가할 수 있다.
			continue
		}
		s := Shot{}
		s.Project = project
		s.Name = name
		if !regexpShotname.MatchString(name) {
			s.Error = "지원하는 샷이름 형식이 아닙니다"
			fails = append(fails, s)
			continue
		}
		now := time.Now().Format(time.RFC3339)
		i := Item{}
		i.Name = name
		i.NetflixID = netflixID
		i.SetSeq() // Name을 이용해서 Seq를 설정한다.
		i.SetCut() // Name을 이용해서 Cut을 설정한다.
		i.Type = typ
		i.UseType = typ
		i.Project = project
		i.ID = i.Project + "_" + i.Name + "_" + i.Type
		i.Shottype = "2d"
		i.Season = season
		i.Episode = episode

		// adminsetting에서 값을 가지고 와서 경로를 설정한다.
		var thumbnailImagePath bytes.Buffer
		thumbnailImagePathTmpl, err := template.New("thumbnailImagePath").Parse(admin.ThumbnailImagePath)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		err = thumbnailImagePathTmpl.Execute(&thumbnailImagePath, i)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		var thumbnailMovPath bytes.Buffer
		thumbnailMovPathTmpl, err := template.New("thumbnailMovPath").Parse(admin.ThumbnailMovPath)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		err = thumbnailMovPathTmpl.Execute(&thumbnailMovPath, i)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		var platePath bytes.Buffer
		thumbnailPlatePathTmpl, err := template.New("pathPath").Parse(admin.PlatePath)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		err = thumbnailPlatePathTmpl.Execute(&platePath, i)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}

		i.Thumpath = thumbnailImagePath.String()
		i.Thummov = thumbnailMovPath.String()
		i.Platepath = platePath.String()
		i.Scantime = now
		i.Updatetime = now
		if i.Type == "org" || i.Type == "left" {
			i.StatusV2 = initStatus
			if setRendersize {
				width := int(float64(pinfo.PlateWidth) * admin.DefaultScaleRatioOfUndistortionPlate)
				height := int(float64(pinfo.PlateHeight) * admin.DefaultScaleRatioOfUndistortionPlate)
				i.Platesize = fmt.Sprintf("%dx%d", pinfo.PlateWidth, pinfo.PlateHeight)
				i.Undistortionsize = fmt.Sprintf("%dx%d", width, height)
				i.Rendersize = fmt.Sprintf("%dx%d", width, height)
			}
		} else {
			i.StatusV2 = "none"
		}
		if genTask {
			// 기본적으로 생성해야할 Task를 추가한다.
			if i.Type == "org" || i.Type == "left" {
				i.Tasks = make(map[string]Task)
				for _, task := range tasks {
					if !task.InitGenerate {
						continue
					}
					if task.Type != "shot" {
						continue
					}
					t := Task{
						Title:        task.Name,
						StatusV2:     initStatus,
						Pipelinestep: task.Pipelinestep, // 파이프라인 스텝을 설정한다.
					}
					i.Tasks[task.Name] = t
				}
			}
		}
		err = i.CheckError()
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		err = addItem(session, i)
		if err != nil {
			s.Error = err.Error()
			fails = append(fails, s)
			continue
		}
		// 폴더 생성 옵션을 체크하면 폴더를 생성한다.
		if mkdir {
			shotRootPath, err := ShotRootPath(i)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			err = GenShotRootPath(shotRootPath)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			seqPath, err := SeqPath(i)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			err = GenSeqPath(seqPath)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			shotPath, err := ShotPath(i)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
			err = GenShotPath(shotPath)
			if err != nil {
				s.Error = err.Error()
				fails = append(fails, s)
				continue
			}
		}
		success = append(success, s)
		// slack log
		err = slacklog(session, project, fmt.Sprintf("Add Shot: %s\nProject: %s, Author: %s", name, project, ssid.ID))
		if err != nil {
			log.Println(err)
			continue
		}
	}

	type recipe struct {
		Success []Shot
		Fails   []Shot
		User    User
		SearchOption
		TrueStatus []string
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Success = success
	rcp.Fails = fails
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	status, err := AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, s := range status {
		if s.DefaultOn {
			rcp.TrueStatus = append(rcp.TrueStatus, s.ID)
		}
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addShot_success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddAsset 함수는 asset을 추가하는 페이지이다.
func handleAddAsset(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	type recipe struct {
		User        User
		Projectlist []string
		SearchOption
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	rcp.Projectlist, err = Projectlist(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "addAsset", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddAssetSubmit 함수는 asset을 생성한다.
func handleAddAssetSubmit(w http.ResponseWriter, r *http.Request) {
	ssid, err := GetSessionID(r)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if ssid.AccessLevel == 0 {
		http.Redirect(w, r, "/invalidaccess", http.StatusSeeOther)
		return
	}
	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()
	initStatusID, err := GetInitStatusID(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project := r.FormValue("Project")
	name := r.FormValue("Name")
	assettype := r.FormValue("Assettype")
	construction := r.FormValue("Construction")
	crowdAsset := str2bool(r.FormValue("CrowdAsset"))
	mkdir := str2bool(r.FormValue("Mkdir"))
	genTask := str2bool(r.FormValue("GenTask"))
	netflixID := r.FormValue("NetflixID")
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}
	names := strings.FieldsFunc(name, f)
	type Asset struct {
		Project string
		Name    string
		Error   string
	}
	var success []Asset
	var fails []Asset
	tasks, err := AllTaskSettings(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, n := range names {
		if n == " " || n == "" { // 사용자가 실수로 여러개의 스페이스를 추가할 수 있다.
			continue
		}
		a := Asset{}
		a.Project = project
		a.Name = n
		if !regexpAssetname.MatchString(n) {
			a.Error = "stone 또는 stone01 형식의 이름이 아닙니다"
			fails = append(fails, a)
			continue
		}
		i := Item{}
		i.Name = n
		i.Type = "asset"
		i.Project = project
		i.ID = i.Project + "_" + i.Name + "_" + i.Type
		i.StatusV2 = initStatusID
		i.Updatetime = time.Now().Format(time.RFC3339)
		i.Assettype = assettype
		i.Assettags = []string{assettype, construction}
		i.CrowdAsset = crowdAsset
		i.NetflixID = netflixID
		if genTask {
			// 기본적으로 생성해야할 Task를 추가한다.
			i.Tasks = make(map[string]Task)
			for _, task := range tasks {
				if !task.InitGenerate {
					continue
				}
				if task.Type != "asset" {
					continue
				}
				t := Task{
					Title:        task.Name,
					StatusV2:     initStatusID,
					Pipelinestep: task.Pipelinestep, // 파이프라인 스텝을 설정한다.
				}
				i.Tasks[task.Name] = t
			}
		}
		err = addItem(session, i)
		if err != nil {
			a.Error = err.Error()
			fails = append(fails, a)
			continue
		}
		// 폴더 생성 옵션을 체크하면 폴더를 생성한다.
		if mkdir {
			assetRootPath, err := AssetRootPath(i)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			err = GenAssetRootPath(assetRootPath)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}

			assetTypePath, err := AssetTypePath(i)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			err = GenAssetTypePath(assetTypePath)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}

			assetPath, err := AssetPath(i)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
			err = GenAssetPath(assetPath)
			if err != nil {
				a.Error = err.Error()
				fails = append(fails, a)
				continue
			}
		}
		success = append(success, a)
	}

	type recipe struct {
		Success []Asset
		Fails   []Asset
		User    User
		SearchOption
		TrueStatus []string
		Setting
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	err = rcp.SearchOption.LoadCookie(session, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Success = success
	rcp.Fails = fails
	u, err := getUser(session, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.User = u
	status, err := AllStatus(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, s := range status {
		if s.DefaultOn {
			rcp.TrueStatus = append(rcp.TrueStatus, s.ID)
		}
	}
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "addAsset_success", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
