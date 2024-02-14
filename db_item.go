package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/digital-idea/ditime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func addItemV2(client *mongo.Client, i Item) error {
	err := i.CheckError()
	if err != nil {
		return err
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"id": i.ID})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("같은 이름을 가진 데이터가 있습니다")
	}
	_, err = collection.InsertOne(ctx, i)
	if err != nil {
		return err
	}
	return nil
}

// DistinctDdline 함수는 프로젝트, dict key를 받아서 key에 사용되는 모든 마감일을 반환한다. 예) 태그
func DistinctDdlineV2(client *mongo.Client, project string, key string) ([]string, error) {
	var results []string
	if project == "" || key == "" {
		return results, nil
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "project", Value: project}}
	values, err := collection.Distinct(ctx, key, filter)
	if err != nil {
		return results, err
	}
	for _, value := range values {
		results = append(results, fmt.Sprintf("%v", value))
	}
	sort.Strings(results)

	if *flagDebug {
		fmt.Println("DB에서 가지고온 마감일 리스트")
		fmt.Println(results)
		fmt.Println()
	}
	var before string
	var datelist []string
	for _, r := range results {
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

func DistinctV2(client *mongo.Client, project string, key string) ([]string, error) {
	var results []string
	if project == "" || key == "" {
		return results, nil
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "project", Value: project}}
	values, err := collection.Distinct(ctx, key, filter)
	if err != nil {
		return results, err
	}

	for _, value := range values {
		if value == nil {
			continue
		}
		results = append(results, fmt.Sprintf("%v", value))
	}
	sort.Strings(results)
	return results, nil
}

func AllAssetsV2(client *mongo.Client, project string) ([]string, error) {
	var results []string
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{
		"project": project,
		"type":    "asset",
	}

	values, err := collection.Distinct(ctx, "name", filter)
	if err != nil {
		return results, err
	}
	for _, value := range values {
		if name, ok := value.(string); ok {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results, nil
}

func TotalnumV2(client *mongo.Client, project string) (Infobarnum, error) {
	if project == "" {
		return Infobarnum{}, nil
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var results Infobarnum

	filter := bson.M{"$or": []bson.M{{"project": project, "type": "org"}, {"project": project, "type": "left"}}}
	num, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return Infobarnum{}, err
	}
	results.Total = int(num)
	return results, nil
}

func AddTagV2(client *mongo.Client, id, inputTag string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	i, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	rmspaceTag := strings.Replace(inputTag, " ", "", -1) // 태그는 공백을 제거한다.
	for _, tag := range i.Tag {
		if rmspaceTag == tag {
			return errors.New(inputTag + "태그는 이미 존재하고 있습니다 추가할 수 없습니다")
		}
	}
	newTags := append(i.Tag, rmspaceTag)

	result, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"tag": newTags, "updatetime": time.Now().Format(time.RFC3339)}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id" + id)
	}
	return nil
}

func getItemV2(client *mongo.Client, id string) (Item, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := Item{}
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func rmItemIDV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

func RmTag(client *mongo.Client, id, inputTag string, isContain bool) (string, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	i, err := getItemV2(client, id)
	if err != nil {
		return "", err
	}
	var newTags []string
	for _, tag := range i.Tag {
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
	i.Tag = newTags
	// 만약 태그에 권정보가 없더라도 권관련 태그는 날아가면 안된다. setItem을 이용한다.

	filter := bson.M{"id": id}
	update := bson.D{{Key: "$set", Value: i}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return i.ID, err
	}
	if result.MatchedCount == 0 {
		return i.ID, errors.New("no document found with id" + i.ID)
	}

	return i.ID, nil
}

func AddTaskV2(client *mongo.Client, id, task, status string) error {
	item, err := getItemV2(client, id)
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

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item.Updatetime = time.Now().Format(time.RFC3339)
	globalStatus, err := AllStatusV2(client)
	if err != nil {
		return err
	}
	item.updateStatusV2(globalStatus)

	filter := bson.M{"id": item.ID}
	update := bson.D{{Key: "$set", Value: item}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + item.ID)
	}
	return nil

}

func RmTaskV2(client *mongo.Client, id, taskname string) error {
	item, err := getItemV2(client, id)
	if err != nil {
		return err
	}

	delete(item.Tasks, taskname)

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item.Updatetime = time.Now().Format(time.RFC3339)

	status, err := AllStatusV2(client)
	if err != nil {
		return err
	}
	item.updateStatusV2(status)

	filter := bson.M{"id": item.ID}
	update := bson.D{{Key: "$set", Value: item}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + item.ID)
	}
	return nil
}

// SetNoteV2 함수는 item에 작업내용을 추가한다. Name,노트내용,에러를 반환한다.
func SetNoteV2(client *mongo.Client, id, userID, text string, overwrite bool) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	i, err := getItemV2(client, id)
	if err != nil {
		return err
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

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"note.text": note, "note.author": userID, "note.date": time.Now().Format(time.RFC3339), "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func AddCommentV2(client *mongo.Client, id, userID, authorName, date, text, media, mediatitle string) error {

	i, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	c := Comment{
		Date:       date,
		Author:     userID,
		AuthorName: authorName,
		Text:       text,
		Media:      media,
		MediaTitle: mediatitle,
	}
	i.Comments = append(i.Comments, c)
	err = setItemV2(client, i)
	if err != nil {
		return err
	}
	return nil
}

func setItemV2(client *mongo.Client, i Item) error {
	i.Updatetime = time.Now().Format(time.RFC3339)
	status, err := AllStatusV2(client)
	if err != nil {
		return err
	}
	i.updateStatusV2(status)
	i.setRnumTag() // 롤넘버에 따른 테그 셋팅
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": i.ID}
	update := bson.D{{Key: "$set", Value: i}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + i.ID)
	}
	return nil
}

func EditCommentV2(client *mongo.Client, id, date, authorName, text, mediatitle, media string) error {
	i, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	var comments []Comment
	for _, c := range i.Comments {
		if c.Date == date {
			c.AuthorName = authorName
			c.Text = text
			c.MediaTitle = mediatitle
			c.Media = media
			comments = append(comments, c)
			continue
		}
		comments = append(comments, c)
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"comments": comments, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + i.ID)
	}
	return nil
}

func RmCommentV2(client *mongo.Client, id, userID, date string) (string, string, error) {
	i, err := getItemV2(client, id)
	if err != nil {
		return id, "", err
	}
	var newComments []Comment
	var removeText string
	for _, comment := range i.Comments {
		if comment.Date == date {
			removeText = comment.Text
			continue
		}
		newComments = append(newComments, comment)
	}
	i.Comments = newComments
	err = setItemV2(client, i)
	if err != nil {
		return id, "", err
	}
	return id, removeText, nil
}

func AddAssetTagV2(client *mongo.Client, id, assettag string) error {
	i, err := getItemV2(client, id)
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

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"assettags": newTags, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + i.ID)
	}
	return nil
}

func RmAssetTagV2(client *mongo.Client, id, inputTag string, isContain bool) error {
	i, err := getItemV2(client, id)
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
	err = setItemV2(client, i)
	if err != nil {
		return err
	}
	return nil
}

func SetShotTypeV2(client *mongo.Client, id, shottype string) error {
	err := validShottype(strings.ToLower(shottype))
	if err != nil {
		return err
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"shottype": strings.ToLower(shottype), "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func GetTaskV2(client *mongo.Client, id, task string) (Task, error) {
	i, err := getItemV2(client, id)
	if err != nil {
		return Task{}, err
	}
	if _, found := i.Tasks[task]; !found {
		return Task{}, errors.New("no task")
	}
	return i.Tasks[task], nil
}

func HasTaskV2(client *mongo.Client, id, task string) error {
	item, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	if _, found := item.Tasks[task]; !found {
		return fmt.Errorf("%s 에 %s Task가 존재하지 않습니다", id, task)
	}
	return nil
}

func SetTaskUserIDV2(client *mongo.Client, id, task, userid string) error {
	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}
	item, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": item.ID}
	update := bson.M{"$set": bson.M{"tasks." + task + ".userid": userid, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetTaskUserV3(client *mongo.Client, id, task, user string) error {
	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}
	item, err := getItemV2(client, id)
	if err != nil {
		return err
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": item.ID}
	update := bson.M{"$set": bson.M{"tasks." + task + ".user": user, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

// SetTaskStatusV3 함수는 item에 task의 status 값을 셋팅한다.
func SetTaskStatusV3(client *mongo.Client, id, task, status string) error {
	item, err := getItemV2(client, id)
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
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 아이템 업데이트 시간을 변경한다.
	item.Updatetime = time.Now().Format(time.RFC3339)
	// 입력받은 상태가 글로벌 status에 존재하는지 체크한다.
	globalStatus, err := AllStatusV2(client)
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

	filter := bson.M{"id": item.ID}
	update := bson.D{{Key: "$set", Value: item}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func setTaskUserCommentV2(client *mongo.Client, id, task, comment string) error {
	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"tasks." + task + ".usercomment": comment}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func setTaskMovV2(client *mongo.Client, id, task, mov string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"tasks." + task + ".mov": mov, "tasks." + task + ".mdate": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func UseTypesV2(client *mongo.Client, id string) (string, []string, error) {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []string
	items := []Item{}

	item, err := getItemV2(client, id)
	if err != nil {
		return item.UseType, results, err
	}

	cursor, err := collection.Find(ctx, bson.M{"name": item.Name, "project": item.Project})
	if err != nil {
		return item.UseType, results, err
	}
	err = cursor.All(ctx, &items)
	if err != nil {
		return item.UseType, results, err
	}

	for _, i := range items {
		results = append(results, i.Type)
	}
	return item.UseType, results, nil
}

func SetUseTypeV2(client *mongo.Client, id, usetype string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"usetype": usetype, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

// SetTaskUserNoteV2 함수는 item에 task의 user note 값을 셋팅한다.
func SetTaskUserNoteV2(client *mongo.Client, id, task, usernote string) error {

	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"tasks." + task + ".usernote": usernote, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetRnumV2(client *mongo.Client, id, rnum string) error {
	item, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	item.Rnum = rnum
	err = setItemV2(client, item)
	if err != nil {
		return err
	}
	return nil
}

func SetDeadline3DV2(client *mongo.Client, id, date string) error {
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"ddline3d": fullTime, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetDeadline2DV2(client *mongo.Client, id, date string) error {
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"ddline2d": fullTime, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetImageSizeV2(client *mongo.Client, id, key, size string) error {
	if !(key == "platesize" || key == "undistortionsize" || key == "rendersize") {
		return errors.New("잘못된 key값입니다")
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{key: size, "updatetime": time.Now().Format(time.RFC3339)}}
	if key == "undistortionsize" {
		update = bson.M{"$set": bson.M{"undistortionsize": size, "updatetime": time.Now().Format(time.RFC3339)}}
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetOverscanRatioV2(client *mongo.Client, id string, ratio float64) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"overscanratio": ratio, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetFrameV2(client *mongo.Client, id, key string, frame int) error {
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

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{key: frame, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetScanTimecodeInV2(client *mongo.Client, id, timecode string) error {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"scantimecodein": timecode, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetScanTimecodeOutV2(client *mongo.Client, id, timecode string) error {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"scantimecodeout": timecode, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetJustTimecodeInV2(client *mongo.Client, id, timecode string) error {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"justtimecodein": timecode, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetJustTimecodeOutV2(client *mongo.Client, id, timecode string) error {
	if !(regexpTimecode.MatchString(timecode) || timecode == "") {
		return fmt.Errorf("%s 문자열은 00:00:00:00 형식의 문자열이 아닙니다", timecode)
	}
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"justtimecodeout": timecode, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetTaskStart(client *mongo.Client, id, task, date string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"tasks." + task + ".start": fullTime, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetTaskEnd(client *mongo.Client, id, task, date string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := HasTaskV2(client, id, task)
	if err != nil {
		return err
	}
	fullTime, err := ditime.ToFullTime(19, date)
	if err != nil {
		return err
	}

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"tasks." + task + ".end": fullTime, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetCameraPubTaskV2(client *mongo.Client, id, task string) error {
	if !(task == "" || task == "mm" || task == "layout" || task == "ani") {
		return errors.New("none(빈문자열), mm, layout, ani 팀만 카메라 publish가 가능합니다")
	}

	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"productioncam.pubtask": task, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetCameraLensmmV2(client *mongo.Client, id, lensmm string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"productioncam.lensmm": lensmm, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetCameraPubPathV2(client *mongo.Client, id, path string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"productioncam.pubpath": path, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetCameraProjectionV2(client *mongo.Client, id string, isProjection bool) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"productioncam.projection": isProjection, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func AddSourceV2(client *mongo.Client, id, author, title, path string) error {
	i, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	for _, i := range i.Sources {
		if i.Title == title {
			return errors.New(title + " already exists")
		}
	}
	s := Source{}
	s.Date = time.Now().Format(time.RFC3339)
	s.Author = author
	s.Title = title
	s.Path = path
	i.Sources = append(i.Sources, s)
	err = setItemV2(client, i)
	if err != nil {
		return err
	}
	return nil
}

func RmSourceV2(client *mongo.Client, id, title string) error {
	i, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	var newSources []Source
	for _, source := range i.Sources {
		if source.Title == title {
			continue
		}
		newSources = append(newSources, source)
	}
	i.Sources = newSources
	err = setItemV2(client, i)
	if err != nil {
		return err
	}
	return nil
}

func AddReferenceV2(client *mongo.Client, id, author, title, path string) error {
	i, err := getItemV2(client, id)
	if err != nil {
		return err
	}
	r := Source{}
	r.Date = time.Now().Format(time.RFC3339)
	r.Author = author
	r.Title = title
	r.Path = path
	i.References = append(i.References, r)
	err = setItemV2(client, i)
	if err != nil {
		return err
	}
	return nil
}

func RmReferenceV2(client *mongo.Client, id, title string) error {
	i, err := getItemV2(client, id)
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
	err = setItemV2(client, i)
	if err != nil {
		return err
	}
	return nil
}

func SetEpisodeV2(client *mongo.Client, id, episode string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"episode": episode, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetSeqV2(client *mongo.Client, id, seq string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"seq": seq, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetPlatePathV2(client *mongo.Client, id, path string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"platepath": path, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetThummovV2(client *mongo.Client, id, path string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"thummov": path, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetBeforemovV2(client *mongo.Client, id, path string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"beforemov": path, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}

func SetAftermovV2(client *mongo.Client, id, path string) error {
	collection := client.Database(*flagDBName).Collection("items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"aftermov": path, "updatetime": time.Now().Format(time.RFC3339)}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no document found with id: " + id)
	}
	return nil
}
