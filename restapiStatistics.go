package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func handleAPIStatisticsDeadlineNum(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()

	date := q.Get("date") // 사용자로부터 "2022-06" 형태로 받는다.

	// 전체 프로젝트 리스트를 구한다.
	var projects []string
	projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	type Recipe struct {
		Projects map[string]int64 `json:"projects"`
		Total    int64            `json:"total"`
	}
	rcp := Recipe{}
	rcp.Projects = make(map[string]int64)
	shotFilter := bson.A{bson.D{{Key: "type", Value: "org"}}, bson.D{{Key: "type", Value: "left"}}}                               // 타입이 샷이면서(일반샷,입체샷)
	dateFilter := bson.D{{Key: "ddline2d", Value: primitive.Regex{Pattern: date, Options: "i"}}, {Key: "$or", Value: shotFilter}} // 날짜가 포함된 아이템 검색
	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		count, err := collection.CountDocuments(ctx, dateFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		rcp.Projects[project] = count
		rcp.Total += count
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIStatisticsNeedDeadlineNum(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 전체 프로젝트 리스트를 구한다.
	var projects []string
	projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	type Recipe struct {
		Projects map[string]int64 `json:"projects"`
		Total    int64            `json:"total"`
	}
	rcp := Recipe{}
	rcp.Projects = make(map[string]int64)
	shotFilter := bson.A{bson.D{{Key: "type", Value: "org"}}, bson.D{{Key: "type", Value: "left"}}} // 타입이 샷이면서(일반샷,입체샷)
	dateFilter := bson.D{{Key: "ddline2d", Value: ""}, {Key: "$or", Value: shotFilter}}             // 날짜가 포함된 아이템 검색
	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		count, err := collection.CountDocuments(ctx, dateFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		rcp.Projects[project] = count
		rcp.Total += count
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIStatisticsShottype(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 전체 프로젝트 리스트를 구한다.
	var projects []string
	projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type TypeNum struct {
		Type2D   int64 `json:"type2d"`
		Type3D   int64 `json:"type3d"`
		TypeNone int64 `json:"typenone"`
	}
	type Recipe struct {
		Projects map[string]TypeNum `json:"projects"`
		TotalNum TypeNum            `json:"totalnum"`
	}
	rcp := Recipe{}
	rcp.Projects = make(map[string]TypeNum)
	shotFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}} // 타입이 샷이면서(일반샷,입체샷)
	type2DFilter := bson.D{{"shottype", "2d"}, {"$or", shotFilter}}         // shottype이 2d 샷 검색
	type3DFilter := bson.D{{"shottype", "3d"}, {"$or", shotFilter}}         // shottype이 3d 샷 검색
	typeNoneFilter := bson.D{{"shottype", ""}, {"$or", shotFilter}}         // shottype이 "" 샷 검색
	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		type2dNum, err := collection.CountDocuments(ctx, type2DFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		type3dNum, err := collection.CountDocuments(ctx, type3DFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		typeNoneNum, err := collection.CountDocuments(ctx, typeNoneFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		typeNum := TypeNum{}
		typeNum.Type2D = type2dNum
		typeNum.Type3D = type3dNum
		typeNum.TypeNone = typeNoneNum
		rcp.Projects[project] = typeNum
		rcp.TotalNum.Type2D += type2dNum
		rcp.TotalNum.Type3D += type3dNum
		rcp.TotalNum.TypeNone += typeNoneNum
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIStatisticsItemtype(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 전체 프로젝트 리스트를 구한다.
	var projects []string
	projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	type TypeNum struct {
		Shot  int64 `json:"shot"`
		Asset int64 `json:"asset"`
	}
	type Recipe struct {
		Projects map[string]TypeNum `json:"projects"`
		TotalNum TypeNum            `json:"totalnum"`
	}
	rcp := Recipe{}
	rcp.Projects = make(map[string]TypeNum)
	shottypeFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}} // 타입이 샷인 조건
	assetFilter := bson.D{{"type", "asset"}}                                    // 타입이 에셋인 것
	shotFilter := bson.D{{"$or", shottypeFilter}}                               // 타입이 샷인 것
	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		shotNum, err := collection.CountDocuments(ctx, shotFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		assetNum, err := collection.CountDocuments(ctx, assetFilter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		typeNum := TypeNum{}
		typeNum.Shot = shotNum
		typeNum.Asset = assetNum
		rcp.Projects[project] = typeNum
		rcp.TotalNum.Shot += shotNum
		rcp.TotalNum.Asset += assetNum
	}
	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI1StatisticsShot(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	type Statusnum struct {
		None    int64 `json:"none"`
		Hold    int64 `json:"hold"`
		Done    int64 `json:"done"`
		Out     int64 `json:"out"`
		Assign  int64 `json:"assign"`
		Ready   int64 `json:"ready"`
		Wip     int64 `json:"wip"`
		Confirm int64 `json:"confirm"`
		Omit    int64 `json:"omit"`
		Client  int64 `json:"client"`
	}
	type Recipe struct {
		Projects map[string]Statusnum `json:"projects"`
		Total    Statusnum            `json:"total"`
	}
	rcp := Recipe{}
	rcp.Projects = make(map[string]Statusnum)

	for _, project := range projects {
		// 프로젝트별로 설정하기.
		currentProjectNum := Statusnum{}
		rcp.Projects[project] = currentProjectNum
	}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI1StatisticsTag(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	type Recipe struct {
		None    int64 `json:"none"`
		Hold    int64 `json:"hold"`
		Done    int64 `json:"done"`
		Out     int64 `json:"out"`
		Assign  int64 `json:"assign"`
		Ready   int64 `json:"ready"`
		Wip     int64 `json:"wip"`
		Confirm int64 `json:"confirm"`
		Omit    int64 `json:"omit"`
		Client  int64 `json:"client"`
	}
	rcp := Recipe{}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI1StatisticsUser(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	task := q.Get("task")
	if task == "" {
		http.Error(w, "Need task name", http.StatusBadRequest)
		return
	}
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	type Recipe struct {
	}
	rcp := Recipe{}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI1StatisticsTask(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	type Recipe struct {
	}
	rcp := Recipe{}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI1StatisticsPipelinestep(w http.ResponseWriter, r *http.Request) {
	// 나중에 Task 입력을 받지 않고 처리할 수 있도록 리펙토링 해야한다.
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	type Recipe struct {
	}
	rcp := Recipe{}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI2StatisticsShot(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 모든 상태를 가지고 옵니다.
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Recipe struct {
		Status  map[string]int64  `json:"status"`
		Filters map[string]bson.D `json:"-"` // 통계를 위한 필터저장에만 사용한다. 반환하지 않는다.
	}
	rcp := Recipe{}
	rcp.Status = make(map[string]int64)
	rcp.Filters = make(map[string]bson.D)
	// filter를 생성합니다.
	shotFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	for _, s := range status {
		rcp.Filters[s.ID] = bson.D{{"statusv2", s.ID}, {"$or", shotFilter}}
	}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		// filter를 for 문 돌면서 나오는 카운트를 검색하고 상태에 넣는다.
		for status, filter := range rcp.Filters {
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rcp.Status[status] += count
		}
	}

	data, err := json.Marshal(rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI2StatisticsAsset(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 모든 상태를 가지고 옵니다.
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Recipe struct {
		Status  map[string]int64  `json:"status"`
		Filters map[string]bson.D `json:"-"` // 통계를 위한 필터저장에만 사용한다. 반환하지 않는다.
	}
	rcp := Recipe{}
	rcp.Status = make(map[string]int64)
	rcp.Filters = make(map[string]bson.D)
	// filter를 생성합니다.
	for _, s := range status {
		rcp.Filters[s.ID] = bson.D{{"statusv2", s.ID}, {"type", "asset"}}
	}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		// filter를 for 문 돌면서 나오는 카운트를 검색하고 상태에 넣는다.
		for status, filter := range rcp.Filters {
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rcp.Status[status] += count
		}
	}

	data, err := json.Marshal(rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI2StatisticsTask(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	task := q.Get("task")
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}

	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 모든 상태를 가지고 옵니다.
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Recipe struct {
		Status  map[string]int64  `json:"status"`
		Filters map[string]bson.D `json:"-"` // 통계를 위한 필터저장에만 사용한다. 반환하지 않는다.
	}
	rcp := Recipe{}
	typeFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	if typ == "asset" {
		typeFilter = bson.A{bson.D{{"type", "asset"}}}
	}
	rcp.Status = make(map[string]int64)
	rcp.Filters = make(map[string]bson.D)
	// filter를 생성합니다.
	for _, s := range status {
		rcp.Filters[s.ID] = bson.D{{"tasks." + task + ".statusv2", s.ID}, {"$or", typeFilter}}
	}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		// filter를 for 문 돌면서 나오는 카운트를 검색하고 상태에 넣는다.
		for status, filter := range rcp.Filters {
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rcp.Status[status] += count
		}
	}

	data, err := json.Marshal(rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI2StatisticsPipelinestep(w http.ResponseWriter, r *http.Request) {
	// 나중에 Task 입력을 받지 않고 처리할 수 있도록 리펙토링 해야한다.
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	task := q.Get("task")
	pipelinestep := q.Get("pipelinestep")
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}

	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = ProjectlistV2(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 모든 상태를 가지고 옵니다.
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Recipe struct {
		Status  map[string]int64  `json:"status"`
		Filters map[string]bson.D `json:"-"` // 통계를 위한 필터저장에만 사용한다. 반환하지 않는다.
	}
	rcp := Recipe{}
	typeFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	if typ == "asset" {
		typeFilter = bson.A{bson.D{{"type", "asset"}}}
	}
	rcp.Status = make(map[string]int64)
	rcp.Filters = make(map[string]bson.D)
	// filter를 생성합니다.
	for _, s := range status {
		rcp.Filters[s.ID] = bson.D{{"tasks." + task + ".statusv2", s.ID}, {"tasks." + task + ".pipelinestep", pipelinestep}, {"$or", typeFilter}}
	}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		// filter를 for 문 돌면서 나오는 카운트를 검색하고 상태에 넣는다.
		for status, filter := range rcp.Filters {
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rcp.Status[status] += count
		}
	}

	data, err := json.Marshal(rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI2StatisticsTag(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	tagName := q.Get("name")
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}

	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 모든 상태를 가지고 옵니다.
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Recipe struct {
		Status  map[string]int64  `json:"status"`
		Filters map[string]bson.D `json:"-"` // 통계를 위한 필터저장에만 사용한다. 반환하지 않는다.
	}
	rcp := Recipe{}
	typeFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	if typ == "asset" {
		typeFilter = bson.A{bson.D{{"type", "asset"}}}
	}
	rcp.Status = make(map[string]int64)
	rcp.Filters = make(map[string]bson.D)
	// filter를 생성합니다.
	for _, s := range status {
		rcp.Filters[s.ID] = bson.D{{"statusv2", s.ID}, {"tag", tagName}, {"$or", typeFilter}}
	}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		// filter를 for 문 돌면서 나오는 카운트를 검색하고 상태에 넣는다.
		for status, filter := range rcp.Filters {
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rcp.Status[status] += count
		}
	}

	data, err := json.Marshal(rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI2StatisticsUser(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	name := q.Get("name")
	task := q.Get("task")
	if task == "" {
		http.Error(w, "Need task name", http.StatusBadRequest)
		return
	}
	typ := q.Get("type")
	if typ == "" {
		typ = "shot"
	}
	if !(typ == "shot" || typ == "asset") {
		http.Error(w, "The type value must be either 'shot' or 'asset' value.", http.StatusBadRequest)
		return
	}

	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 모든 상태를 가지고 옵니다.
	status, err := AllStatusV2(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Recipe struct {
		Status  map[string]int64  `json:"status"`
		Filters map[string]bson.D `json:"-"` // 통계를 위한 필터저장에만 사용한다. 반환하지 않는다.
	}
	rcp := Recipe{}
	typeFilter := bson.A{bson.D{{"type", "org"}}, bson.D{{"type", "left"}}}
	if typ == "asset" {
		typeFilter = bson.A{bson.D{{"type", "asset"}}}
	}
	rcp.Status = make(map[string]int64)
	rcp.Filters = make(map[string]bson.D)
	// filter를 생성합니다.
	for _, s := range status {
		rcp.Filters[s.ID] = bson.D{{"statusv2", s.ID}, {"tasks." + task + ".user", primitive.Regex{Pattern: name, Options: "i"}}, {"$or", typeFilter}}
	}

	for _, project := range projects {
		collection := client.Database("project").Collection(project)
		// filter를 for 문 돌면서 나오는 카운트를 검색하고 상태에 넣는다.
		for status, filter := range rcp.Filters {
			count, err := collection.CountDocuments(ctx, filter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rcp.Status[status] += count
		}
	}

	data, err := json.Marshal(rcp.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPI1StatisticsAsset(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := r.URL.Query()
	project := q.Get("project")
	var projects []string
	if project != "" {
		projects = append(projects, project)
	} else {
		projects, err = client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	type Recipe struct {
	}
	rcp := Recipe{}

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleAPIStatisticsProjectnum(w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projects, err := client.Database("projectinfo").ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Recipe struct {
		Count    int      `json:"count"`
		Projects []string `json:"projects"`
	}
	rcp := Recipe{}
	rcp.Projects = projects
	rcp.Count = len(projects)

	data, err := json.Marshal(rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
