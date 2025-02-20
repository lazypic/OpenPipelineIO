package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CheckBsonMStructure 함수는 bson.M 구조를 체크합니다.
func CheckBsonMStructure(data interface{}) bool {
	// data가 bson.M 타입인지 확인합니다.
	if reflect.TypeOf(data) != reflect.TypeOf(bson.M{}) {
		fmt.Println("Data is not of type bson.M")
		return false
	}

	// bson.M으로 타입 어설션을 수행합니다.
	bsonData, _ := data.(bson.M)

	// bson.M 내의 각 항목을 순회하며 타입을 체크합니다.
	for key, value := range bsonData {
		valueType := reflect.TypeOf(value)
		switch value.(type) {
		case string, int, float64, bool, primitive.Regex:
			// 기본 타입이거나 MongoDB 특정 타입인 경우
			fmt.Printf("Key: %s, Type: %s is valid\n", key, valueType)
		case bson.A, bson.M:
			// 배열이나 다른 문서인 경우 (재귀적 검사를 추가할 수 있습니다.)
			fmt.Printf("Key: %s, Type: %s is a BSON array or document and is valid\n", key, valueType)
		default:
			// 예상치 못한 타입인 경우
			fmt.Printf("Key: %s has an unexpected type: %s\n", key, valueType)
			return false
		}
	}

	return true
}

// GenQueryV2 함수는 검색옵션을 받아서 검색옵션과 쿼리를 반환한다.
func GenQueryV2(client *mongo.Client, op SearchOption) (SearchOption, bson.M) {
	// 검색어중 연산에 필요한 검색어는 제거한다.
	var words []string
	var selectTasks []string
	// Task 처리
	allTasks, err := TaskSettingNamesV2(client)
	if err != nil {
		log.Println(err)
	}
	if op.Task != "" {
		selectTasks = append(selectTasks, op.Task)
	}

	for _, word := range strings.Split(op.Searchword, " ") {
		// task를 searchbox UX가 아닌 타이핑으로도 선언할 수 있어야 한다.
		if strings.HasPrefix(word, "task:") {
			selectTasks = append(selectTasks, strings.TrimPrefix(word, "task:"))
			continue
		}
		switch word {
		case "":
		case "or", "||":
		case "and", "&&":
		default:
			words = append(words, word)
		}
	}
	for _, word := range strings.Split(op.Searchword, " ") {
		// task를 searchbox UX가 아닌 타이핑으로도 선언할 수 있어야 한다.
		if strings.HasPrefix(word, "task:") {
			selectTasks = append(selectTasks, strings.TrimPrefix(word, "task:"))
			continue
		}
		switch word {
		case "":
		case "or", "||":
		case "and", "&&":
		default:
			words = append(words, word)
		}
	}

	wordQueries := []bson.M{}
	if *flagDebug {
		fmt.Println(words)
	}
	for _, word := range words {
		query := []bson.M{}
		if MatchShortTime.MatchString(word) {
			// 1121 형식의 날짜
			regFullTime := fmt.Sprintf(`^\d{4}-%s-%sT\d{2}:\d{2}:\d{2}[-+]\d{2}:\d{2}$`, word[0:2], word[2:4])
			if len(selectTasks) == 0 {
				for _, task := range allTasks {
					query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".date": &primitive.Regex{Pattern: regFullTime}})
					query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".end": &primitive.Regex{Pattern: regFullTime}})
				}
				query = append(query, bson.M{"ddline2d": &primitive.Regex{Pattern: regFullTime}})
				query = append(query, bson.M{"ddline3d": &primitive.Regex{Pattern: regFullTime}})
			} else {
				for _, task := range selectTasks {
					query = append(query, bson.M{"tasks." + task + ".date": &primitive.Regex{Pattern: regFullTime}})
					query = append(query, bson.M{"tasks." + task + ".end": &primitive.Regex{Pattern: regFullTime}})
				}
			}
			query = append(query, bson.M{"name": &primitive.Regex{Pattern: word}}) // 샷 이름에 숫자가 포함되는 경우도 검색한다.
		} else if MatchNormalTime.MatchString(word) {
			// 데일리 날짜를 검색한다.
			// 2016-11-21 형태는 데일리로 간주합니다.
			// jquery 달력의 기본형식이기도 합니다.
			regFullTime := fmt.Sprintf(`^%sT\d{2}:\d{2}:\d{2}[-+]\d{2}:\d{2}$`, word)
			if len(selectTasks) == 0 {
				for _, task := range allTasks {
					query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".mdate": &primitive.Regex{Pattern: regFullTime}})
				}
			} else {
				for _, task := range selectTasks {
					query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".mdate": &primitive.Regex{Pattern: regFullTime}})
				}
			}
		} else if regexpTimecode.MatchString(word) {
			query = append(query, bson.M{"justtimecodein": word})
			query = append(query, bson.M{"justtimecodeout": word})
			query = append(query, bson.M{"scantimecodein": word})
			query = append(query, bson.M{"scantimecodeout": word})
		} else if strings.HasPrefix(word, "tag:") {
			query = append(query, bson.M{"tag": strings.TrimPrefix(word, "tag:")})
		} else if strings.HasPrefix(word, "assettags:") {
			query = append(query, bson.M{"assettags": strings.TrimPrefix(word, "assettags:")})
		} else if strings.HasPrefix(word, "deadline2d:") {
			if word == "deadline2d:" {
				query = append(query, bson.M{"ddline2d": ""}) // Deadline2D 마감일이 빈 문자열이라면 빈문자열인 데이터만 검색되어야 한다.
			} else {
				query = append(query, bson.M{"ddline2d": &primitive.Regex{Pattern: strings.TrimPrefix(word, "deadline2d:"), Options: "i"}})
			}
		} else if strings.HasPrefix(word, "deadline3d:") {
			if word == "deadline3d:" {
				query = append(query, bson.M{"ddline3d": ""}) // Deadline3D 마감일이 빈 문자열이라면 빈문자열인 데이터만 검색되어야 한다.
			} else {
				query = append(query, bson.M{"ddline3d": &primitive.Regex{Pattern: strings.TrimPrefix(word, "deadline3d:"), Options: "i"}})
			}
		} else if strings.HasPrefix(word, "shottype:") {
			query = append(query, bson.M{"shottype": &primitive.Regex{Pattern: strings.TrimPrefix(word, "shottype:"), Options: "i"}})
		} else if strings.HasPrefix(word, "type:shot") {
			query = append(query, bson.M{"$or": []bson.M{
				{"type": "org"},
				{"type": "main"},
				{"type": "mp"},
				{"type": "left"},
			}})
		} else if strings.HasPrefix(word, "type:asset") {
			query = append(query, bson.M{"type": "asset"})
		} else if strings.HasPrefix(word, "name:") {
			query = append(query, bson.M{"name": &primitive.Regex{Pattern: strings.TrimPrefix(word, "name:"), Options: "i"}})
		} else if strings.HasPrefix(word, "episode:") {
			query = append(query, bson.M{"episode": &primitive.Regex{Pattern: strings.TrimPrefix(word, "episode:"), Options: "i"}})
		} else if strings.HasPrefix(word, "season:") {
			query = append(query, bson.M{"season": &primitive.Regex{Pattern: strings.TrimPrefix(word, "season:"), Options: "i"}})
		} else if strings.HasPrefix(word, "status:") {
			status := strings.TrimPrefix(word, "status:")
			// 검색바에서 task를 선택했다면,
			if len(selectTasks) != 0 {
				for _, task := range selectTasks {
					query = append(query, bson.M{"tasks." + task + ".statusv2": status})
				}
			} else {
				// 검색바에서 Task가 All 이면
				query = append(query, bson.M{"statusv2": status})

			}
		} else if strings.HasPrefix(word, "user:") {
			if len(selectTasks) == 0 {
				if strings.TrimPrefix(word, "user:") == "notassign" {
					for _, task := range allTasks {
						query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".user": ""})
					}
				} else {
					for _, task := range allTasks {
						query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".user": &primitive.Regex{Pattern: strings.TrimPrefix(word, "user:")}})
					}
				}
			} else {
				for _, task := range selectTasks {
					if strings.TrimPrefix(word, "user:") == "notassign" {
						query = append(query, bson.M{"tasks." + task + ".user": ""})
					} else {
						query = append(query, bson.M{"tasks." + task + ".user": &primitive.Regex{Pattern: strings.TrimPrefix(word, "user:")}})
					}
				}
			}
		} else if strings.HasPrefix(word, "usercomment:") {
			userComment := strings.TrimPrefix(word, "usercomment:")
			if len(selectTasks) == 0 {
				for _, task := range allTasks {
					if userComment != "" {
						query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".usercomment": userComment})
					} else {
						query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".usercomment": ""})
					}
				}
			} else {
				for _, task := range selectTasks {
					if userComment != "" {
						query = append(query, bson.M{"tasks." + task + ".usercomment": userComment})
					} else {
						query = append(query, bson.M{"tasks." + task + ".usercomment": ""})
					}
				}
			}
		} else if strings.HasPrefix(word, "rnum:") { // 롤넘버 형태일 때
			query = append(query, bson.M{"rnum": &primitive.Regex{Pattern: strings.TrimPrefix(word, "rnum:"), Options: "i"}})
		} else {
			switch word {
			case "all", "All", "ALL", "올", "미ㅣ", "dhf", "전체":
				query = append(query, bson.M{})
			case "shot", "샷", "전샷", "전체샷":
				query = append(query, bson.M{"type": "org"})
				query = append(query, bson.M{"type": "main"})
				query = append(query, bson.M{"type": "mp"})
				query = append(query, bson.M{"type": "left"})
			case "asset", "assets", "에셋":
				query = append(query, bson.M{"type": "asset"})
			default:
				query = append(query, bson.M{"id": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"comments.text": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"sources.title": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"sources.path": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"references.title": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"references.path": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"note.text": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"tag": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"assettags": &primitive.Regex{Pattern: word, Options: "i"}})
				query = append(query, bson.M{"scanname": &primitive.Regex{Pattern: word, Options: ""}})
				query = append(query, bson.M{"rnum": &primitive.Regex{Pattern: word, Options: ""}})
				// Task가 선언 되어있을 때
				if len(selectTasks) == 0 {
					for _, task := range allTasks {
						query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".user": &primitive.Regex{Pattern: word}})        // 아티스트명을 검색한다.
						query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".usercomment": &primitive.Regex{Pattern: word}}) // UserComment를 검색한다.
					}
				} else {
					for _, task := range selectTasks {
						query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".user": &primitive.Regex{Pattern: word}})        // 아티스트명을 검색한다.
						query = append(query, bson.M{"tasks." + strings.ToLower(task) + ".usercomment": &primitive.Regex{Pattern: word}}) // UserComment를 검색한다.
					}
				}
			}
		}

		wordQueries = append(wordQueries, bson.M{"$or": query})
	}

	statusQueries := []bson.M{}
	if len(selectTasks) == 0 {
		// 검색바가 All 이면 검색바 옵션 True status 리스트만 쿼리에 추가한다.
		for _, status := range op.TrueStatus {
			statusQueries = append(statusQueries, bson.M{"statusv2": status})
		}
	} else {
		// 만약 검색바에서 Task가 선택되어 있다면..
		// op(SearchOption)에서 true 상태 리스트만 가지고 온다.
		// for문을 돌면서 해당 쿼리를 추가한다.
		for _, status := range op.TrueStatus {
			for _, task := range selectTasks {
				statusQueries = append(statusQueries, bson.M{"tasks." + task + ".statusv2": status})
			}
		}
	}
	// 각 단어에 대한 쿼리를 and 로 검색할지 or 로 검색할지 결정한다.
	expression := "$and"
	for _, word := range strings.Split(op.Searchword, " ") {
		if word == "or" || word == "||" {
			expression = "$or"
		}
	}
	queries := []bson.M{
		{expression: wordQueries},
	}
	// 상태 쿼리가 존재하면 상태에 대해서 or 처리한다.
	if len(statusQueries) != 0 {
		queries = append(queries, bson.M{"$or": statusQueries})
	}

	// 프로젝트 쿼리에 대해서 and 처리를 진행한다.
	projectQueries := []bson.M{}
	if op.Project != "" { // 빈문자열일 때 전체 프로젝트를 검색한다.
		projectQueries = append(projectQueries, bson.M{"project": op.Project})
	}
	if len(projectQueries) != 0 {
		queries = append(queries, bson.M{"$and": projectQueries})
	}

	// 최종쿼리를 지정한다.
	q := bson.M{"$and": queries}

	// 정렬설정
	switch op.Sortkey {
	// 스캔길이, 스캔날짜는 역순으로 정렬한다.
	// 스캔길이는 보통 난이도를 결정하기 때문에 역순(긴 길이순)을 매니저인 팀장,실장은 우선적으로 봐야한다.
	// 스캔날짜는 IO팀에서 최근 등록한 데이터를 많이 검토하기 때문에 역순(최근등록순)으로 봐야한다.
	case "scanframe", "scantime":
		op.Sortkey = "-" + op.Sortkey
	case "taskdate":
		if len(selectTasks) != 0 {
			op.Sortkey = "tasks." + op.Task + ".date"
		}
	case "taskpredate":
		if len(selectTasks) != 0 {
			op.Sortkey = "tasks." + op.Task + ".end"
		}
	case "": // 기본적으로 id로 정렬한다.
		op.Sortkey = "id"
	}

	return op, q
}

// SearchV2 함수는 다음 검색함수이다.
func SearchV2(client *mongo.Client, op SearchOption) ([]Item, error) {
	results := []Item{}
	// 검색어가 없다면 바로 빈 값을 리턴한다.
	if op.Searchword == "" {
		return results, nil
	}

	if len(op.TrueStatus) == 0 {
		// 선택된 상태가 없다면 바로 리턴한다.
		return results, nil
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	o, q := GenQueryV2(client, op)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: o.Sortkey, Value: 1}})
	//findOptions.SetSkip(int64(CachedAdminSetting.ItemNumberOfPage * (op.Page - 1)))
	//findOptions.SetLimit(int64(CachedAdminSetting.ItemNumberOfPage))

	cursor, err := collection.Find(ctx, q, findOptions)

	//cursor, err := collection.Find(ctx, q)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// SearchPage 함수는 페이지로 검색하는 함수이다. "아이템, totalpagenum, 에러" 를 반환한다.
func SearchPageV2(client *mongo.Client, op SearchOption) ([]Item, int, error) {
	results := []Item{}
	// 서비스를 최초에 설치할 때는 Page 갯수가 0이되고 0으로 나누면 서버에서 에러가 나기 때문에 아래 에러처리가 필요하다.
	if CachedAdminSetting.ItemNumberOfPage == 0 {
		return results, 0, errors.New("페이지에 보이는 아이템 갯수가 0이 될 수 없습니다. /adminsetting 에 접속 후 페이지 설정이 필요합니다")
	}
	if op.Page <= 0 {
		op.Page = 1
	}
	// 검색어가 없다면 바로 빈 값을 리턴한다.
	if op.Searchword == "" {
		return results, 0, nil
	}
	if len(op.TrueStatus) == 0 {
		// 선택된 상태가 없다면 바로 리턴한다.
		return results, 0, nil
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	o, q := GenQueryV2(client, op)

	findOptions := options.Find()

	findOptions.SetSort(bson.D{{Key: o.Sortkey, Value: 1}})
	findOptions.SetSkip(int64(CachedAdminSetting.ItemNumberOfPage * (op.Page - 1)))
	findOptions.SetLimit(int64(CachedAdminSetting.ItemNumberOfPage))

	cursor, err := collection.Find(ctx, q, findOptions)

	//cursor, err := collection.Find(ctx, q)
	if err != nil {

		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	totalItemNum, err := collection.CountDocuments(ctx, q)
	if err != nil {
		return results, 0, err
	}
	totalPageNum := int(totalItemNum) / CachedAdminSetting.ItemNumberOfPage
	if int(totalItemNum)%CachedAdminSetting.ItemNumberOfPage != 0 {
		totalPageNum++
	}
	return results, totalPageNum, nil
}
