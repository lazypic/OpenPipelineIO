package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/digital-idea/ditime"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func setItem(session *mgo.Session, i Item) error {
	session.SetMode(mgo.Monotonic, true)
	i.Updatetime = time.Now().Format(time.RFC3339)
	status, err := AllStatus(session)
	if err != nil {
		return err
	}
	i.updateStatusV2(status)
	i.setRnumTag() // 롤넘버에 따른 테그 셋팅
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": i.ID}, i)
	if err != nil {
		return err
	}
	return nil
}

func getItem(session *mgo.Session, id string) (Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	var result Item
	err := c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		return Item{}, err
	}
	return result, nil
}

// Shot 함수는 프로젝트명, 샷이름을 이용해서 샷정보를 반환한다.
func Shot(session *mgo.Session, project string, name string) (Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	var result Item
	// org, left 갯수를 구한다.
	orgnum, err := c.Find(bson.M{"id": name + "_org"}).Count()
	if err != nil {
		return Item{}, err
	}
	leftnum, err := c.Find(bson.M{"id": name + "_left"}).Count()
	if err != nil {
		return Item{}, err
	}
	q := bson.M{"id": name + "_org"}
	if leftnum == 1 && orgnum == 0 {
		q = bson.M{"id": name + "_left"}
	}
	err = c.Find(q).One(&result)
	if err != nil {
		return Item{}, err
	}
	return result, nil
}

// Asset 함수는 프로젝트명, 에셋 이름을 입력받아 에셋정보를 반환한다.
func Asset(session *mgo.Session, project string, name string) (Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	var result Item
	err := c.Find(bson.M{"id": name + "_asset"}).One(&result)
	if err != nil {
		return Item{}, err
	}
	return result, nil
}

// AllAsset 함수는 모든 에셋리스트를 가지고 옵니다.
func AllAssets(session *mgo.Session, project string) ([]string, error) {
	var result []string
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Find(bson.M{"type": "asset"}).Distinct("name", &result)
	if err != nil {
		return nil, err
	}
	sort.Strings(result)
	return result, nil
}

// SearchName 함수는 입력된 문자열이 'name'키 값에 포함되어 있다면 해당 아이템을 반환한다.
func SearchName(session *mgo.Session, project string, name string) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	results := []Item{}
	err := c.Find(bson.M{"name": &bson.RegEx{Pattern: name, Options: "i"}}).Sort("name").All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// UseTypes 함수는 project, name을 받아서 사용되는 모든 타입을 반환한다.
func UseTypes(session *mgo.Session, project string, name string) ([]string, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	items := []Item{}
	err := c.Find(bson.M{"name": name}).Sort("type").All(&items)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, i := range items {
		results = append(results, i.Type)
	}
	return results, nil
}

// Seqs 함수는 프로젝트 이름을 받아서 seq 리스트를 반환한다.
func Seqs(session *mgo.Session, project string) ([]string, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	var results []Item
	err := c.Find(bson.M{"$or": []bson.M{{"type": "org"}, {"type": "left"}}}).Select(bson.M{"seq": 1}).All(&results)
	if err != nil {
		return nil, err
	}
	keys := make(map[string]bool)
	for _, result := range results {
		seq := result.Seq
		keys[seq] = true
	}
	seqs := []string{}
	for k := range keys {
		seqs = append(seqs, k)
	}
	sort.Strings(seqs)
	return seqs, nil
}

// Shots 함수는 프로젝트 이름과 입력된 시퀀스가 'name'키 값에 포함되어 있다면 shots 리스트를 반환한다.
func Shots(session *mgo.Session, project string, seq string) ([]string, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	var results []Item
	query := bson.M{"name": &bson.RegEx{Pattern: seq, Options: "i"}, "type": bson.M{"$in": []string{"org", "left"}}}
	err := c.Find(query).Select(bson.M{"name": 1}).All(&results)
	if err != nil {
		return nil, err
	}
	keys := make(map[string]bool)
	for _, result := range results {
		if !strings.Contains(result.Name, "_") {
			continue
		}
		shot := strings.Split(result.Name, "_")[1]
		keys[shot] = true
	}
	shots := []string{}
	for k := range keys {
		shots = append(shots, k)
	}
	sort.Strings(shots)
	return shots, nil
}

func rmItem(session *mgo.Session, project, name, usertyp string) error {
	session.SetMode(mgo.Monotonic, true)
	var typ string
	if usertyp == "" {
		t, err := Type(session, project, name)
		if err != nil {
			return err
		}
		typ = t
	} else {
		typ = usertyp
	}
	c := session.DB(*flagDBName).C("items")
	num, err := c.Find(bson.M{"name": name, "type": typ}).Count()
	if err != nil {
		return err
	}
	if num == 0 {
		return errors.New("삭제할 아이템이 없습니다")
	}
	err = c.Remove(bson.M{"name": name, "type": typ})
	if err != nil {
		return err
	}
	return nil
}

// Distinct 함수는 프로젝트, dict key를 받아서 key에 사용되는 모든 문자열을 반환한다. 예) 태그
func Distinct(session *mgo.Session, project string, key string) ([]string, error) {
	var result []string
	if project == "" || key == "" {
		return result, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Find(bson.M{"project": project}).Distinct(key, &result)
	if err != nil {
		return nil, err
	}
	sort.Strings(result)
	return result, nil
}

// DistinctDdline 함수는 프로젝트, dict key를 받아서 key에 사용되는 모든 마감일을 반환한다. 예) 태그
func DistinctDdline(session *mgo.Session, project string, key string) ([]string, error) {
	var result []string
	if project == "" || key == "" {
		return result, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Find(bson.M{"project": project}).Distinct(key, &result)
	if err != nil {
		return nil, err
	}
	//result로 datelist를 만든다.
	sort.Strings(result)
	if *flagDebug {
		fmt.Println("DB에서 가지고온 마감일 리스트")
		fmt.Println(result)
		fmt.Println()
	}
	var before string
	var datelist []string
	for _, r := range result {
		if r != "" {
			date := ToNormalTime(r)
			if date == before {
				break
			} else {
				datelist = append(datelist, date)
			}
			before = date
		}
	}
	sort.Strings(datelist) //기존 OpenPipelineIO 2.0의 4자리 수를 위하여 정렬한다. 추후 이 줄은 사라진다.
	if *flagDebug {
		fmt.Println("마감일을 Tag형태로 바꾼 리스트")
		fmt.Println(datelist)
		fmt.Println()
	}
	return datelist, nil
}

// SearchAllShot 함수는 모든 Shot 데이터를 반환한다. Excel을 뽑을 때 사용한다.
func SearchAllShot(session *mgo.Session, project, sortkey string) ([]Item, error) {
	results := []Item{}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	queries := []bson.M{}
	queries = append(queries, bson.M{"project": project, "type": "org"})
	queries = append(queries, bson.M{"project": project, "type": "left"})
	q := bson.M{"$or": queries}
	err := c.Find(q).Sort(sortkey).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchAllAsset 함수는 모든 Asset 데이터를 반환한다. Excel을 뽑을 때 사용한다.
func SearchAllAsset(session *mgo.Session, project, sortkey string) ([]Item, error) {
	results := []Item{}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Find(bson.M{"project": project, "type": "asset"}).Sort(sortkey).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchAll 함수는 모든 Shot, Asset 데이터를 반환한다. Excel을 뽑을 때 사용한다.
func SearchAll(session *mgo.Session, project, sortkey string) ([]Item, error) {
	results := []Item{}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	queries := []bson.M{}
	queries = append(queries, bson.M{"project": project, "type": "org"})
	queries = append(queries, bson.M{"project": project, "type": "left"})
	queries = append(queries, bson.M{"project": project, "type": "asset"})
	q := bson.M{"$or": queries}
	err := c.Find(q).Sort(sortkey).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchTag 함수는 태그를 검색할때 사용한다.
// SearchTags라고 이름을 붙히지 않은 이유는 OpenPiplineIO 자료구조의 필드명이 Tag이기 때문이다.
// 미래에 Tag필드를  Tags 필드로 바꾼후 이 함수의 이름을 SearchTags로 바꿀 예정이다.
func SearchTag(session *mgo.Session, op SearchOption) ([]Item, error) {
	return SearchKey(session, op, "tag")
}

// SearchAssettags 함수는 검색옵션으로 에셋태그를 검색할때 사용한다.
func SearchAssettags(session *mgo.Session, op SearchOption) ([]Item, error) {
	return SearchKey(session, op, "assettags")
}

// SearchKey 함수는 Item.{key} 필드의 값과 검색어가 정확하게 일치하는 항목들만 검색한다.
func SearchKey(session *mgo.Session, op SearchOption, key string) ([]Item, error) {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	query := []bson.M{}
	// 아래 단어는 OpenPipelineIO에 버튼으로 되어있는 태그단어를 클릭시 작동되는 예약어이다.
	switch op.Searchword {
	case "2d", "2D":
		query = append(query, bson.M{"shottype": &bson.RegEx{Pattern: "2d", Options: "i"}})
	case "3d", "3D":
		query = append(query, bson.M{"shottype": &bson.RegEx{Pattern: "3d", Options: "i"}})
	case "asset":
		query = append(query, bson.M{"type": &bson.RegEx{Pattern: "asset", Options: "i"}})
	default:
		query = append(query, bson.M{key: op.Searchword})
	}

	var results []Item

	status := []bson.M{}

	q := bson.M{"$and": []bson.M{
		{"$or": query},
		{"$or": status},
	}}
	err := c.Find(q).Sort(op.Sortkey).All(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SearchStatusNum 함수는 검색된 결과에 대한 상태별 갯수를 검색한다.
func SearchStatusNum(op SearchOption, items []Item) (Infobarnum, error) {
	var results Infobarnum
	results.Search = len(items)
	results.StatusNum = make(map[string]int) // statusV2의 갯수를 처리하기 위해 StatusNum 맵을 초기화한다.
	for _, item := range items {
		if item.Shottype == "2D" || item.Shottype == "2d" {
			results.Shot2d++
		}
		if item.Shottype == "3D" || item.Shottype == "3d" {
			results.Shot3d++
		}
		if item.Type == "asset" {
			results.Assets++
		}
		if item.Type == "org" || item.Type == "left" {
			results.Shot++
		}
		if op.Task == "" {
			results.StatusNum[item.StatusV2]++
		} else {
			// task가 존재하지 않으면 넘긴다.
			if _, ok := item.Tasks[op.Task]; !ok {
				continue
			}
			results.StatusNum[item.Tasks[op.Task].StatusV2]++
		}
	}
	return results, nil
}

// Totalnum 함수는 프로젝트의 전체샷에 대한 상태 갯수를 검색한다.
func Totalnum(session *mgo.Session, project string) (Infobarnum, error) {
	if project == "" {
		return Infobarnum{}, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")

	var results Infobarnum
	totalnum, err := c.Find(bson.M{"$or": []bson.M{{"project": project, "type": "org"}, {"project": project, "type": "left"}}}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Total = totalnum
	return results, nil
}

// TotalTaskStatusnum 함수는 프로젝트의 전체샷에 대한 상태 갯수를 검색한다. // legacy
func TotalTaskStatusnum(session *mgo.Session, project, task string) (Infobarnum, error) {
	if project == "" {
		return Infobarnum{}, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")

	var results Infobarnum
	results.StatusNum = make(map[string]int)

	// statusv2
	statuslist, err := AllStatus(session)
	if err != nil {
		return Infobarnum{}, err
	}
	for _, s := range statuslist {
		query := bson.M{"$and": []bson.M{
			{"tasks." + task + ".statusv2": s.ID},
			{"$or": []bson.M{{"type": "org"}, {"type": "left"}, {"type": "asset"}}},
		}}
		num, err := c.Find(query).Count()
		if err != nil {
			continue
		}
		results.StatusNum[s.ID] = num
	}
	// 전체 아이템 갯수를 구한다.
	totalnum, err := c.Find(bson.M{"$or": []bson.M{{"type": "org"}, {"type": "left"}, {"type": "asset"}}}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Total = totalnum
	// 샷 갯수를 구한다.
	shotnum, err := c.Find(bson.M{"$or": []bson.M{{"type": "org"}, {"type": "left"}}}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot = shotnum
	// 샷2D 갯수를 구한다.
	queryShot2D := bson.M{"$and": []bson.M{
		{"shottype": "2d"},
		{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
	}}
	shot2dNum, err := c.Find(queryShot2D).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot2d = shot2dNum
	// 샷3D 갯수를 구한다.
	queryShot3D := bson.M{"$and": []bson.M{
		{"shottype": "3d"},
		{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
	}}
	shot3dNum, err := c.Find(queryShot3D).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot3d = shot3dNum
	// 에셋 갯수를 구한다.
	assetnum, err := c.Find(bson.M{"$or": []bson.M{{"type": "asset"}}}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Assets = assetnum
	// 진행률을 구한다
	results.calculatePercent()
	return results, nil
}

// TotalTaskAndUserStatusnum 함수는 프로젝트,테스크,유저의 샷에 대한 상태 갯수를 검색한다.
func TotalTaskAndUserStatusnum(session *mgo.Session, project, task, user string) (Infobarnum, error) {
	if project == "" {
		return Infobarnum{}, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")

	var results Infobarnum
	results.StatusNum = make(map[string]int)

	statuslist, err := AllStatus(session)
	if err != nil {
		return Infobarnum{}, err
	}
	for _, s := range statuslist {
		query := bson.M{"$and": []bson.M{
			{"tasks." + task + ".statusv2": s.ID},
			{"tasks." + task + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
			{"$or": []bson.M{{"type": "org"}, {"type": "left"}, {"type": "asset"}}},
		}}
		num, err := c.Find(query).Count()
		if err != nil {
			continue
		}
		results.StatusNum[s.ID] = num
	}
	// 전체 아이템 갯수를 구한다.
	queryTotal := bson.M{"$and": []bson.M{
		{"tasks." + task + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
		{"$or": []bson.M{{"type": "org"}, {"type": "left"}, {"type": "asset"}}},
	}}
	totalnum, err := c.Find(queryTotal).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Total = totalnum
	// 샷 갯수를 구한다.
	queryShot := bson.M{"$and": []bson.M{
		{"tasks." + task + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
		{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
	}}
	shotnum, err := c.Find(queryShot).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot = shotnum
	// 샷2D 갯수를 구한다.
	queryShot2D := bson.M{"$and": []bson.M{
		{"shottype": "2d"},
		{"tasks." + task + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
		{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
	}}
	shot2dNum, err := c.Find(queryShot2D).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot2d = shot2dNum
	// 샷3D 갯수를 구한다.
	queryShot3D := bson.M{"$and": []bson.M{
		{"shottype": "3d"},
		{"tasks." + task + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
		{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
	}}
	shot3dNum, err := c.Find(queryShot3D).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot3d = shot3dNum
	// 에셋 갯수를 구한다.
	queryAsset := bson.M{"$and": []bson.M{
		{"tasks." + task + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
		{"$or": []bson.M{{"type": "asset"}}},
	}}
	assetnum, err := c.Find(queryAsset).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Assets = assetnum
	// 진행률을 구한다
	results.calculatePercent()
	return results, nil
}

// TotalUserStatusnum 함수는 프로젝트,테스크,유저의 샷에 대한 상태 갯수를 검색한다.
func TotalUserStatusnum(session *mgo.Session, project, user string) (Infobarnum, error) {
	if project == "" {
		return Infobarnum{}, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")

	var results Infobarnum
	results.StatusNum = make(map[string]int)
	tasks, err := AllTaskSettings(session)
	if err != nil {
		return Infobarnum{}, err
	}

	// statusv2
	statuslist, err := AllStatus(session)
	if err != nil {
		return Infobarnum{}, err
	}
	for _, status := range statuslist {
		var querys []bson.M
		for _, t := range tasks {
			query := bson.M{"$and": []bson.M{
				{"tasks." + t.Name + ".statusv2": status.ID},
				{"tasks." + t.Name + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
				{"$or": []bson.M{{"type": "org"}, {"type": "left"}, {"type": "asset"}}},
			}}
			querys = append(querys, query)
		}
		num, err := c.Find(bson.M{"$or": querys}).Count()
		if err != nil {
			continue
		}
		results.StatusNum[status.ID] = num
	}

	// 전체 아이템 갯수를 구한다.
	var querysTotal []bson.M
	for _, t := range tasks {
		queryTotal := bson.M{"$and": []bson.M{
			{"tasks." + t.Name + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
			{"$or": []bson.M{{"type": "org"}, {"type": "left"}, {"type": "asset"}}},
		}}
		querysTotal = append(querysTotal, queryTotal)
	}
	totalnum, err := c.Find(bson.M{"$or": querysTotal}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Total = totalnum

	// 샷 갯수를 구한다.
	var querysShot []bson.M
	for _, t := range tasks {
		queryShot := bson.M{"$and": []bson.M{
			{"tasks." + t.Name + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
			{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
		}}
		querysShot = append(querysShot, queryShot)
	}
	shotnum, err := c.Find(bson.M{"$or": querysShot}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot = shotnum

	// 샷2D 갯수를 구한다.
	var querysShot2D []bson.M
	for _, t := range tasks {
		queryShot2D := bson.M{"$and": []bson.M{
			{"shottype": "2d"},
			{"tasks." + t.Name + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
			{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
		}}
		querysShot2D = append(querysShot2D, queryShot2D)
	}
	shot2dNum, err := c.Find(bson.M{"$or": querysShot2D}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot2d = shot2dNum

	// 샷3D 갯수를 구한다.
	var querysShot3D []bson.M
	for _, t := range tasks {
		queryShot3D := bson.M{"$and": []bson.M{
			{"shottype": "3d"},
			{"tasks." + t.Name + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
			{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
		}}
		querysShot3D = append(querysShot3D, queryShot3D)
	}
	shot3dNum, err := c.Find(bson.M{"$or": querysShot3D}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot3d = shot3dNum

	// 에셋 갯수를 구한다.
	var querysAsset []bson.M
	for _, t := range tasks {
		queryAsset := bson.M{"$and": []bson.M{
			{"tasks." + t.Name + ".user": &bson.RegEx{Pattern: user, Options: "i"}},
			{"$or": []bson.M{{"type": "asset"}}},
		}}
		querysAsset = append(querysAsset, queryAsset)
	}
	assetnum, err := c.Find(bson.M{"$or": querysAsset}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Assets = assetnum

	// 진행률을 구한다
	results.calculatePercent()
	return results, nil
}

// TotalStatusnum 함수는 프로젝트의 전체샷에 대한 상태 갯수를 검색한다.
func TotalStatusnum(session *mgo.Session, project string) (Infobarnum, error) {
	if project == "" {
		return Infobarnum{}, nil
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")

	var results Infobarnum
	results.StatusNum = make(map[string]int)

	// statusv2
	statuslist, err := AllStatus(session)
	if err != nil {
		return Infobarnum{}, err
	}
	for _, s := range statuslist {
		query := bson.M{"$and": []bson.M{
			{"statusv2": s.ID},
			{"$or": []bson.M{{"type": "org"}, {"type": "left"}, {"type": "asset"}}},
		}}
		num, err := c.Find(query).Count()
		if err != nil {
			continue
		}
		results.StatusNum[s.ID] = num
	}
	// 전체 아이템 갯수를 구한다.
	totalnum, err := c.Find(bson.M{"$or": []bson.M{{"type": "org"}, {"type": "left"}, {"type": "asset"}}}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Total = totalnum
	// 샷 갯수를 구한다.
	shotnum, err := c.Find(bson.M{"$or": []bson.M{{"type": "org"}, {"type": "left"}}}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot = shotnum
	// 샷2D 갯수를 구한다.
	queryShot2D := bson.M{"$and": []bson.M{
		{"shottype": "2d"},
		{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
	}}
	shot2dNum, err := c.Find(queryShot2D).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot2d = shot2dNum
	// 샷3D 갯수를 구한다.
	queryShot3D := bson.M{"$and": []bson.M{
		{"shottype": "3d"},
		{"$or": []bson.M{{"type": "org"}, {"type": "left"}}},
	}}
	shot3dNum, err := c.Find(queryShot3D).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Shot3d = shot3dNum
	// 에셋 갯수를 구한다.
	assetnum, err := c.Find(bson.M{"$or": []bson.M{{"type": "asset"}}}).Count()
	if err != nil {
		return Infobarnum{}, err
	}
	results.Assets = assetnum
	// 진행률을 구한다
	results.calculatePercent()
	return results, nil
}

// setTaskMov함수는 해당 샷에 mov를 설정하는 함수이다.
func setTaskMov(session *mgo.Session, project, name, task, mov string) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", err
	}
	c := session.DB(*flagDBName).C("items")
	typ, err := Type(session, project, name)
	if err != nil {
		return "", err
	}
	id := name + "_" + typ
	err = HasTask(session, id, task)
	if err != nil {
		return id, err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".mov": mov, "tasks." + task + ".mdate": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// setTaskExpectDay함수는 해당 샷에 예상일을 설정하는 함수이다.
func setTaskExpectDay(session *mgo.Session, project, id, task string, expectDay int) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = HasTask(session, id, task)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".expectday": expectDay}})
	if err != nil {
		return err
	}
	return nil
}

// setTaskResultDay함수는 해당 샷에 예상일을 설정하는 함수이다.
func setTaskResultDay(session *mgo.Session, project, id, task string, resultDay int) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = HasTask(session, id, task)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".resultday": resultDay}})
	if err != nil {
		return err
	}
	return nil
}

// setTaskUserComment함수는 해당 아이템의 Task에 UserComment를 설정하는 함수이다.
func setTaskUserComment(session *mgo.Session, id, task, comment string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := HasTask(session, id, task)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".usercomment": comment}})
	if err != nil {
		return err
	}
	return nil
}

// Type 함수는 이름을 이용해서 Type값을 반환한다.
// 일반샷은 org를 반환한다.
// 입체샷인경우 left를 반환한다.
// 에셋은 asset을 반환한다.
func Type(session *mgo.Session, project, name string) (string, error) {
	c := session.DB(*flagDBName).C("items")
	var items []Item
	err := c.Find(bson.M{"$or": []bson.M{{"name": name, "type": "org"}, {"name": name, "type": "left"}, {"name": name, "type": "asset"}}}).All(&items)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", errors.New(name + "에 해당하는 org,left,asset 타입을 DB에서 찾을 수 없습니다.")
	}
	if len(items) != 1 {
		return "", errors.New(name + "값이 DB에서 고유하지 않습니다.")
	}
	return items[0].Type, nil
}

func GetID(session *mgo.Session, project, name string) (string, error) {
	c := session.DB(*flagDBName).C("items")
	var items []Item
	err := c.Find(bson.M{"$and": []bson.M{{"name": name, "project": project}, {"name": name, "project": project}}}).All(&items)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", errors.New(name + "에 해당하는 id를 DB에서 찾을 수 없습니다.")
	}
	if len(items) != 1 {
		return "", errors.New(name + "값이 DB에서 고유하지 않습니다.")
	}
	return items[0].ID, nil
}

// SetTimecode 함수는 item에 Timecode를 설정한다.
// ScanTimecodeIn,ScanTimecodeOut,JustTimecodeIn,JustTimecoeOut 문자를 key로 사용할 수 있다.
func SetTimecode(session *mgo.Session, project, name, key, timecode string) error {
	key = strings.ToLower(key)
	if !(key == "scantimecodein" ||
		key == "scantimecodeout" ||
		key == "justtimecodein" ||
		key == "justtimecodeout") {
		return errors.New("scantimecodein, scantimecodeout, justtimecodein, justtimecodeout 키값만 사용가능")
	}
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{key: timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	// 우리회사는 현재 timecode와 keycode를 혼용해서 사용중이다.
	// 원래는 Timecode가 맞지만 현재 DB가 keycode로 되어있어 아직은 아래줄이 필요하다.
	key = strings.Replace(key, "timecode", "keycode", -1)
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{key: timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetUseType 함수는 item에 UseType string을 설정한다.
func SetUseType(session *mgo.Session, project, id, usetype string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"usetype": usetype, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetFrame 함수는 item에 프레임을 설정한다.
// ScanIn,ScanOut,ScanFrame,PlateIn,PlateOut,JustIn,JustOut,HandleIn,HandleOut 문자를 key로 사용할 수 있다.
func SetFrame(session *mgo.Session, id, key string, frame int) error {
	if frame == -1 {
		return nil
	}
	if !(key == "scanin" ||
		key == "scanout" ||
		key == "scanframe" ||
		key == "platein" ||
		key == "plateout" ||
		key == "justin" ||
		key == "justout" ||
		key == "handlein" ||
		key == "handleout") {
		return errors.New("scanin, scanout, scanframe, platein, plateout, justin, justout, handlein, handleout 키값만 사용가능합니다")
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{key: frame, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCameraPubPath 함수는 해당 카메라 퍼블리쉬 경로를 설정한다.
func SetCameraPubPath(session *mgo.Session, project, id, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"productioncam.pubpath": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCameraPubTask 함수는 해당 카메라 퍼블리쉬 팀을 설정한다.
func SetCameraPubTask(session *mgo.Session, project, id, task string) error {
	if !(task == "" || task == "mm" || task == "layout" || task == "ani") {
		return errors.New("none(빈문자열), mm, layout, ani 팀만 카메라 publish가 가능합니다")
	}
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"productioncam.pubtask": task, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCameraLensmm 함수는 해당 아이템에 카메라 렌즈mm를 설정한다.
func SetCameraLensmm(session *mgo.Session, project, id, lensmm string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"productioncam.lensmm": lensmm, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCameraProjection 함수는 샷에 Projection 카메라 사용여부를 체크한다.
func SetCameraProjection(session *mgo.Session, project, id string, isProjection bool) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"productioncam.projection": isProjection, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetObjectID 함수는 Item에 Object In, Out 값을 설정한다.
func SetObjectID(session *mgo.Session, project, name string, in, out int) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ != "asset" {
		return errors.New("asset 타입이 아닙니다")
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"objectidin": in, "objectidout": out, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}


// SetSeason 함수는 item에 season 값을 셋팅한다.
func SetSeason(session *mgo.Session, project, id, season string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"season": season, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetEpisode 함수는 item에 episode 값을 셋팅한다.
func SetEpisode(session *mgo.Session, project, id, episode string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"episode": episode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetOverscanRatio 함수는 item에 OverscanRatio 값을 셋팅한다.
func SetOverscanRatio(session *mgo.Session, project, id string, ratio float64) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"overscanratio": ratio, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetPlatePath 함수는 item에 PlatePath값을 셋팅한다.
func SetPlatePath(session *mgo.Session, project, id, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"platepath": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetThummov 함수는 item에 Thummov값을 셋팅한다.
func SetThummov(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"thummov": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetBeforemov 함수는 item에 Before mov값을 셋팅한다.
func SetBeforemov(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"beforemov": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetAftermov 함수는 item에 After mov값을 셋팅한다.
func SetAftermov(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": name + "_" + typ}, bson.M{"$set": bson.M{"aftermov": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetEditmov 함수는 item에 Edit(편집본) mov값을 셋팅한다.
func SetEditmov(session *mgo.Session, project, id, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"editmov": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskStatus 함수는 item에 task의 status 값을 셋팅한다.
func SetTaskStatus(session *mgo.Session, project, id, task, status string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	statusNum := ""

	if statusNum == "" {
		return errors.New("올바른 status가 아닙니다")
	}
	if _, found := item.Tasks[strings.ToLower(task)]; !found {
		return errors.New("task가 존재하지 않습니다")
	}
	t := item.Tasks[task]
	t.StatusV2 = status
	item.Tasks[task] = t

	c := session.DB(*flagDBName).C("items")
	item.Updatetime = time.Now().Format(time.RFC3339)
	globalStatus, err := AllStatus(session)
	if err != nil {
		return err
	}
	item.updateStatusV2(globalStatus)
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// SetTaskPipelinestep 함수는 item에 task의 pipelinestep 값을 셋팅한다.
func SetTaskPipelinestep(session *mgo.Session, project, id, task, pipelinestep string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	t := item.Tasks[task]
	item.Tasks[task] = t
	c := session.DB(*flagDBName).C("items")
	item.Updatetime = time.Now().Format(time.RFC3339)
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// SetTaskStatusV2 함수는 item에 task의 status 값을 셋팅한다.
func SetTaskStatusV2(session *mgo.Session, id, task, status string) error {
	session.SetMode(mgo.Monotonic, true)
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	if _, found := item.Tasks[strings.ToLower(task)]; !found {
		return fmt.Errorf("%s 에 %s task가 존재하지 않습니다", id, task)
	}
	t := item.Tasks[task]
	t.StatusV2 = status
	item.Tasks[task] = t

	// 앞으로 바뀔 상태
	c := session.DB(*flagDBName).C("items")
	// 아이템 업데이트 시간을 변경한다.
	item.Updatetime = time.Now().Format(time.RFC3339)
	// 입력받은 상태가 글로벌 status에 존재하는지 체크한다.
	globalStatus, err := AllStatus(session)
	if err != nil {
		return err
	}
	hasStatus := false
	for _, s := range globalStatus {
		if s.ID == status {
			hasStatus = true
			break
		}
	}
	if !hasStatus {
		return fmt.Errorf("%s status가 존재하지 않습니다", status)
	}
	// 아이템의 statusV2를 업데이트한다.
	item.updateStatusV2(globalStatus)
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// HasTask 함수는 item에 task가 존재하는 체크한다.
func HasTask(session *mgo.Session, id, task string) error {
	session.SetMode(mgo.Monotonic, true)
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	if _, found := item.Tasks[task]; !found {
		return fmt.Errorf("%s 에 %s Task가 존재하지 않습니다", id, task)
	}
	return nil
}

// AddTask 함수는 item에 task를 추가한다.
func AddTask(session *mgo.Session, project, id, task, status, pipelinestep string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	taskname := strings.ToLower(task)
	// 기존에 Task가 없다면 추가한다.
	if _, found := item.Tasks[task]; !found {
		t := Task{}
		t.Title = taskname
		t.StatusV2 = status
		item.Tasks[task] = t
	} else {
		return fmt.Errorf("이미 %s 에 %s Task가 존재합니다", id, taskname)
	}
	c := session.DB(*flagDBName).C("items")
	item.Updatetime = time.Now().Format(time.RFC3339)
	globalStatus, err := AllStatus(session)
	if err != nil {
		return err
	}
	item.updateStatusV2(globalStatus)
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// RmTask 함수는 item에 task를 제거한다.
func RmTask(session *mgo.Session, project, id, taskname string) error {
	session.SetMode(mgo.Monotonic, true)
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	delete(item.Tasks, taskname)
	c := session.DB(*flagDBName).C("items")
	item.Updatetime = time.Now().Format(time.RFC3339)
	status, err := AllStatus(session)
	if err != nil {
		return err
	}
	item.updateStatusV2(status)
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// SetTaskUserV2 함수는 item에 task의 user 값을 셋팅한다.
func SetTaskUserV2(session *mgo.Session, id, task, user string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasTask(session, id, task)
	if err != nil {
		return err
	}
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": item.ID}, bson.M{"$set": bson.M{"tasks." + task + ".user": user, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskUserID 함수는 item에 task의 userid 값을 셋팅한다.
func SetTaskUserID(session *mgo.Session, id, task, userid string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasTask(session, id, task)
	if err != nil {
		return err
	}
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": item.ID}, bson.M{"$set": bson.M{"tasks." + task + ".userid": userid, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskDate 함수는 item에 task에 마감일을 셋팅한다.
func SetTaskDate(session *mgo.Session, id, task, date string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".date": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskDuration 함수는 item에 start, end 일을 셋팅한다.
func SetTaskDuration(session *mgo.Session, project, id, task, start, end string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	startTime, err := ditime.ToFullTime(10, start)
	if err != nil {
		return err
	}
	endTime, err := ditime.ToFullTime(19, end)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".start": startTime, "tasks." + task + ".date": endTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskDuration1st 함수는 item에 1차 마감일 start, end 일을 셋팅한다.
func SetTaskDuration1st(session *mgo.Session, project, id, task, start, end string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	startTime, err := ditime.ToFullTime(10, start)
	if err != nil {
		return err
	}
	endTime, err := ditime.ToFullTime(19, end)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".start": startTime, "tasks." + task + ".end": endTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskDuration2nd 함수는 item에 2차 마감일 start, end 일을 셋팅한다.
func SetTaskDuration2nd(session *mgo.Session, project, id, task, start, end string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	startTime, err := ditime.ToFullTime(10, start)
	if err != nil {
		return err
	}
	endTime, err := ditime.ToFullTime(19, end)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".startdate2nd": startTime, "tasks." + task + ".date": endTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetDeadline2D 함수는 item에 2D마감일을 셋팅한다.
func SetDeadline2D(session *mgo.Session, project, id, date string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	fullTime := ""
	if date != "" {
		fullTime, err = ditime.ToFullTime(19, date)
		if err != nil {
			return err
		}
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"ddline2d": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetDeadline3D 함수는 item에 3D마감일을 셋팅한다.
func SetDeadline3D(session *mgo.Session, project, id, date string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	fullTime := ""
	if date != "" {
		fullTime, err = ditime.ToFullTime(19, date)
		if err != nil {
			return err
		}
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"ddline3d": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskStartdate 함수는 item에 task의 startdate 값을 셋팅한다.
func SetTaskStartdate(session *mgo.Session, id, task, date string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := HasTask(session, id, task)
	if err != nil {
		return err
	}
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".start": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskUserNote 함수는 item에 task의 user note 값을 셋팅한다.
func SetTaskUserNote(session *mgo.Session, project, name, task, usernote string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	err = HasTask(session, id, task)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".usernote": usernote, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetTaskPredate 함수는 item에 task의 predate 값을 셋팅한다.
func SetTaskPredate(session *mgo.Session, id, task, date string) (string, error) {
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(*flagDBName).C("items")
	err := HasTask(session, id, task)
	if err != nil {
		return id, err
	}
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return id, err
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tasks." + task + ".end": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// SetShotType 함수는 item에 shot type을 셋팅한다.
func SetShotType(session *mgo.Session, project, name, shottype string) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", err
	}
	if typ == "asset" {
		return "", fmt.Errorf("%s 는 asset type 입니다. 변경할 수 없습니다", name)
	}
	id := name + "_" + typ
	err = validShottype(strings.ToLower(shottype))
	if err != nil {
		return id, err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"shottype": strings.ToLower(shottype), "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, err
	}
	return id, nil
}

// SetOutputName 함수는 item에 Outputname 을 셋팅한다.
func SetOutputName(session *mgo.Session, project, name, outputname string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 는 %s type 입니다. 변경할 수 없습니다", name, typ)
	}
	if outputname == "" {
		return errors.New("outputname 이 빈 문자열 입니다")
	}
	id := name + "_" + typ
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"outputname": outputname, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetRetimePlate 함수는 item에 RetimePlate를 셋팅한다.
func SetRetimePlate(session *mgo.Session, project, name, retimeplate string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 는 %s type 입니다. retime plate를 설정할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"retimeplate": retimeplate, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetOCIOcc 함수는 item에 OCIO .cc를 셋팅한다.
func SetOCIOcc(session *mgo.Session, project, name, path string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	if typ == "asset" {
		return fmt.Errorf("%s 는 %s type 입니다. 설정할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"ociocc": path, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetRollmedia 함수는 item에 Setellite Rollmedia를 셋팅한다.
func SetRollmedia(session *mgo.Session, project, name, rollmedia string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"rollmedia": rollmedia, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetScanname 함수는 item에 Scanname을 셋팅한다.
func SetScanname(session *mgo.Session, project, id, scanname string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"scanname": scanname, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetRnum 함수는 샷에 롤넘버를 설정한다.
func SetRnum(session *mgo.Session, project, id, rnum string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	item.Rnum = rnum
	err = setItem(session, item)
	if err != nil {
		return err
	}
	return nil
}

// SetAssetType 함수는 item에 assettype을 셋팅한다.
func SetAssetType(session *mgo.Session, project, name, assettype string) (string, string, string, error) {
	_, err := validAssettype(assettype)
	if err != nil {
		return "", "", assettype, err
	}
	session.SetMode(mgo.Monotonic, true)
	err = HasProject(session, project)
	if err != nil {
		return "", "", assettype, err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", "", assettype, err
	}
	if typ != "asset" {
		return "", "", assettype, fmt.Errorf("%s 아이템은 %s 타입입니다. 처리할 수 없습니다", name, typ)
	}
	id := name + "_" + typ
	i, err := getItem(session, id)
	if err != nil {
		return id, "", assettype, err
	}
	beforeType := i.Assettype
	i.Assettype = assettype
	i.setAssettags()
	err = setItem(session, i)
	if err != nil {
		return id, beforeType, assettype, err
	}
	return id, beforeType, assettype, nil
}

// SetScanTimecodeIn 함수는 item에 Scan Timecode In을 셋팅한다.
func SetScanTimecodeIn(session *mgo.Session, id, timecode string) error {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"scantimecodein": timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetScanTimecodeOut 함수는 item에 Scan Timecode Out을 셋팅한다.
func SetScanTimecodeOut(session *mgo.Session, id, timecode string) error {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"scantimecodeout": timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetJustTimecodeIn 함수는 item에 Just Timecode In을 셋팅한다.
func SetJustTimecodeIn(session *mgo.Session, id, timecode string) error {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"justtimecodein": timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetJustTimecodeOut 함수는 item에 Just Timecode In을 셋팅한다.
func SetJustTimecodeOut(session *mgo.Session, id, timecode string) error {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"justtimecodeout": timecode, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}


// SetFindate 함수는 item에 최종 데이터 아웃풋 날짜를 셋팅한다.
func SetFindate(session *mgo.Session, project, name, date string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"findate": fullTime, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// SetCrowdAsset 함수는 item에 crowdtype을 설정한다.
func SetCrowdAsset(session *mgo.Session, project, name string) (string, bool, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", false, err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return "", false, err
	}
	id := name + "_" + typ
	c := session.DB(*flagDBName).C("items")
	item, err := getItem(session, id)
	if err != nil {
		return id, item.CrowdAsset, err
	}
	invertBool := !item.CrowdAsset
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"crowdasset": invertBool, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return id, invertBool, err
	}
	return id, invertBool, nil
}

// AddTag 함수는 item에 tag를 셋팅한다.
func AddTag(session *mgo.Session, project, id, inputTag string) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", err
	}
	i, err := getItem(session, id)
	if err != nil {
		return id, err
	}
	rmspaceTag := strings.Replace(inputTag, " ", "", -1) // 태그는 공백을 제거한다.
	for _, tag := range i.Tag {
		if rmspaceTag == tag {
			return id, errors.New(inputTag + "태그는 이미 존재하고 있습니다 추가할 수 없습니다")
		}
	}
	newTags := append(i.Tag, rmspaceTag)
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"tag": newTags, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return i.Name, err
	}
	return i.Name, nil
}

// AddAssetTag 함수는 item에 tag를 셋팅한다.
func AddAssetTag(session *mgo.Session, project, id, assettag string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	i, err := getItem(session, id)
	if err != nil {
		return err
	}
	rmspaceTag := strings.Replace(assettag, " ", "", -1) // 태그는 공백을 제거한다.
	for _, tag := range i.Assettags {
		if rmspaceTag == tag {
			return errors.New(assettag + "태그는 이미 존재하고 있습니다 추가할 수 없습니다")
		}
	}
	newTags := append(i.Assettags, rmspaceTag)
	c := session.DB(*flagDBName).C("items")
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"assettags": newTags, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	return nil
}

// RenameTag 함수는 item의 Tag를 리네임한다.
func RenameTag(session *mgo.Session, project, before, after string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	var items []Item
	err := c.Find(bson.M{}).All(&items)
	if err != nil {
		return err
	}
	for _, i := range items {
		beforeTags := i.Tag
		var newTags []string
		for _, t := range i.Tag {
			if t == before {
				newTags = append(newTags, after)
			} else {
				newTags = append(newTags, t)
			}
		}
		if !reflect.DeepEqual(beforeTags, newTags) {
			err = c.Update(bson.M{"id": i.ID}, bson.M{"$set": bson.M{"tag": newTags, "updatetime": time.Now().Format(time.RFC3339)}})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetTags 함수는 item에 tag를 교체한다.
func SetTags(session *mgo.Session, project, name string, tags []string) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	i, err := getItem(session, id)
	if err != nil {
		return err
	}
	i.Tag = tags
	// 만약 태그에 권정보가 없더라도 권관련 태그는 날아가면 안된다. setItem을 이용한다.
	err = setItem(session, i)
	if err != nil {
		return err
	}
	return nil
}

// RmAssetTag 함수는 item에 assettag를 삭제한다.
func RmAssetTag(session *mgo.Session, project, id, inputTag string, isContain bool) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	i, err := getItem(session, id)
	if err != nil {
		return err
	}
	var newTags []string
	for _, tag := range i.Assettags {
		if isContain {
			if strings.Contains(tag, inputTag) {
				continue
			}
		}
		if inputTag == tag {
			continue
		}
		newTags = append(newTags, tag)
	}
	i.Assettags = newTags
	// 만약 태그에 권정보가 없더라도 권관련 태그는 날아가면 안된다. setItem을 이용한다.
	err = setItem(session, i)
	if err != nil {
		return err
	}
	return nil
}

// SetNote 함수는 item에 작업내용을 추가한다. Name,노트내용,에러를 반환한다.
func SetNote(session *mgo.Session, project, id, userID, text string, overwrite bool) (string, string, error) {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return "", "", err
	}
	c := session.DB(*flagDBName).C("items")
	i, err := getItem(session, id)
	if err != nil {
		return "", "", err
	}
	var note string
	if overwrite {
		note = text
	} else {
		if strings.HasSuffix(i.Note.Text, "\n") {
			note = text + i.Note.Text
		} else {
			note = text + "\n " + i.Note.Text
		}
	}
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"note.text": note, "note.author": userID, "note.date": time.Now().Format(time.RFC3339), "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return i.Name, "", err
	}
	return i.Name, note, nil
}

// RmReference 함수는 item에서 레퍼런스를 삭제합니다.
func RmReference(session *mgo.Session, id, title string) error {
	session.SetMode(mgo.Monotonic, true)
	i, err := getItem(session, id)
	if err != nil {
		return err
	}
	var newReferences []Source
	for _, ref := range i.References {
		if ref.Title == title {
			continue
		}
		newReferences = append(newReferences, ref)
	}
	i.References = newReferences
	err = setItem(session, i)
	if err != nil {
		return err
	}
	return nil
}

// GetTask 함수는 item의 Task 정보를 반환한다.
func GetTask(session *mgo.Session, id, task string) (Task, error) {
	session.SetMode(mgo.Monotonic, true)
	i, err := getItem(session, id)
	if err != nil {
		return Task{}, err
	}
	if _, found := i.Tasks[task]; !found {
		return Task{}, errors.New("task가 존재하지 않습니다")
	}
	return i.Tasks[task], nil
}

// GetShottype 함수는 item의 Shottype 정보를 반환한다.
func GetShottype(session *mgo.Session, project, name string) (string, error) {
	session.SetMode(mgo.Monotonic, true)
	typ, err := Type(session, project, name)
	if err != nil {
		return "", err
	}
	id := name + "_" + typ
	i, err := getItem(session, id)
	if err != nil {
		return "", err
	}
	return i.Shottype, nil
}

// setTaskPublish함수는 해당 샷 Task에 Publish를 설정하는 함수이다.
func addTaskPublish(session *mgo.Session, project, name, task, key string, p Publish) error {
	session.SetMode(mgo.Monotonic, true)
	err := HasProject(session, project)
	if err != nil {
		return err
	}
	c := session.DB(*flagDBName).C("items")
	typ, err := Type(session, project, name)
	if err != nil {
		return err
	}
	id := name + "_" + typ
	err = HasTask(session, id, task)
	if err != nil {
		return err
	}
	err = c.Update(
		bson.M{"id": id},
		bson.M{"$push": bson.M{fmt.Sprintf("tasks.%s.publishes.%s", task, key): p}})
	if err != nil {
		return err
	}
	return nil
}

// rmTaskPublishKey 함수는 item > tasks > publishes 를 제거한다.
func rmTaskPublishKey(session *mgo.Session, project, id, taskname, key string) error {
	session.SetMode(mgo.Monotonic, true)
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	_, ok := item.Tasks[taskname].Publishes[key]
	if ok {
		delete(item.Tasks[taskname].Publishes, key)
	} else {
		return fmt.Errorf("no publish key: %s", key)
	}
	c := session.DB(*flagDBName).C("items")
	item.Updatetime = time.Now().Format(time.RFC3339)
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// rmTaskPublish 함수는 item > tasks > publishes > 하나의 아이템을 제거한다.
func rmTaskPublish(session *mgo.Session, project, id, taskname, key, createtime, path string) error {
	session.SetMode(mgo.Monotonic, true)
	item, err := getItem(session, id)
	if err != nil {
		return err
	}
	var keepList []Publish
	pubList := item.Tasks[taskname].Publishes[key]
	for _, p := range pubList {
		if p.Createtime == createtime && p.Path == path {
			continue // 삭제데이터의 조건이 맞다면 keepList에 넣지 않는다.
		}
		keepList = append(keepList, p)
	}
	item.Tasks[taskname].Publishes[key] = keepList // 퍼블리쉬 리스트를 교체한다.
	c := session.DB(*flagDBName).C("items")
	item.Updatetime = time.Now().Format(time.RFC3339)
	err = c.Update(bson.M{"id": item.ID}, item)
	if err != nil {
		return err
	}
	return nil
}

// HasItem 함수는 입력받은 project에 해당 id를 가진 item이 존재하는지 체크한다. mongoDB의 objectID가 아닌 OpenPipelineIO 내에서 정의된 id를 사용한다.
func HasItem(session *mgo.Session, project, id string) error {
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(*flagDBName).C("items")
	num, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return err
	}
	if num > 0 {
		return nil
	}
	return errors.New(project + " 프로젝트에 해당 Item이 존재하지 않습니다.")
}
