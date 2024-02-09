package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

// handleInputMode 함수는 수정을 편하게 입력하는 페이지 이다.
func handleInputMode(w http.ResponseWriter, r *http.Request) {
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
		User
		SessionID   string
		Projectlist []string
		Items       []Item
		Ddline3d    []string
		Ddline2d    []string
		Tags        []string
		Assettags   []string
		AllAssets   []string
		SearchOption
		Searchnum                Infobarnum
		Totalnum                 Infobarnum
		Projectinfo              Project
		MailDNS                  string
		OS                       string
		TasksettingNames         []string
		TasksettingOrderMap      map[string]float64
		Dday                     string
		Status                   []Status
		AllStatusIDs             []string
		TotalPageNum             int
		Setting                  Setting
		FullCalendarEventJson    string
		FullCalendarResourceJson string
	}
	rcp := recipe{}
	rcp.Setting = CachedAdminSetting
	_, rcp.OS, _ = GetInfoFromRequestHeader(r)
	rcp.MailDNS = CachedAdminSetting.EmailDNS

	rcp.SessionID = ssid.ID
	client, err := initMongoClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	tasks, err := AllTaskSettingsV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.TasksettingOrderMap = make(map[string]float64)
	for _, t := range tasks {
		rcp.TasksettingOrderMap[t.Name] = t.Order
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

	for _, status := range rcp.Status {
		rcp.AllStatusIDs = append(rcp.AllStatusIDs, status.ID)
	}
	rcp.SearchOption = handleRequestToSearchOption(r)
	rcp.User, err = getUserV2(client, ssid.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Projectlist, err = OnProjectlistV2(client)
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

	if len(rcp.Projectlist) == 0 {
		http.Redirect(w, r, "/noonproject", http.StatusSeeOther)
		return
	}

	rcp.Ddline3d, err = DistinctDdlineV2(client, rcp.SearchOption.Project, "ddline3d")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Ddline2d, err = DistinctDdlineV2(client, rcp.SearchOption.Project, "ddline2d")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Tags, err = DistinctV2(client, rcp.SearchOption.Project, "tag")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Assettags, err = DistinctV2(client, rcp.SearchOption.Project, "assettags")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.AllAssets, err = AllAssetsV2(client, rcp.SearchOption.Project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Totalnum, err = TotalnumV2(client, rcp.SearchOption.Project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rcp.Totalnum.calculatePercent()
	if rcp.SearchOption.Project != "" {
		rcp.Projectinfo, err = getProjectV2(client, rcp.SearchOption.Project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dday, err := ToDday(rcp.Projectinfo.Deadline)
		if err != nil {
			log.Println(err)
		}
		rcp.Dday = dday
	}

	// 페이지 검색을 진행한다. 페이지수에 맞는 아이템 갯수만 반환해야한다.
	rcp.Items, rcp.TotalPageNum, err = SearchPageV2(client, rcp.SearchOption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 검색바에 출력되는 갯수를 연산한다. 전체에서 갯수를 구해야한다.
	items, err := SearchV2(client, rcp.SearchOption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.Searchnum, err = SearchStatusNum(rcp.SearchOption, items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 달력 이벤트를 렌더링한다.

	event, resource, err := ItemsToFCEventsAndFCResource(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	eventJson, err := json.Marshal(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.FullCalendarEventJson = string(eventJson)
	resourceJson, err := json.Marshal(resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rcp.FullCalendarResourceJson = string(resourceJson)
	// fmt.Println(rcp.FullCalendarResourceJson)
	// fmt.Println(rcp.FullCalendarEventJson)
	// 최종적으로 사용된 프로젝트명을 쿠키에 저장한다.
	cookie := http.Cookie{
		Name:   "Project",
		Value:  rcp.SearchOption.Project,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Task",
		Value:  rcp.SearchOption.Task,
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{
		Name:   "Searchword",
		Value:  base64.StdEncoding.EncodeToString([]byte(rcp.SearchOption.Searchword)), //  쿠키는 UTF-8을 저장할 때 에러가 발생한다.
		MaxAge: 0,
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "text/html")
	err = TEMPLATES.ExecuteTemplate(w, "index", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
