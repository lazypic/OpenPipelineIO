package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// searchReviewV2 함수는 Review를 검색한다.
func searchReviewV2(client *mongo.Client, searchword string) ([]Review, error) {
	var results []Review
	if searchword == "" {
		return results, nil
	}
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	allQueries := []bson.M{}
	for _, word := range strings.Split(searchword, " ") {
		if len(word) < 2 { // 한글자의 단어는 무시한다.
			continue
		}
		orQueries := []bson.M{}
		if strings.HasPrefix(word, "createtime:") {
			if strings.TrimPrefix(word, "createtime:") != "" {
				orQueries = append(orQueries, bson.M{"createtime": &primitive.Regex{Pattern: strings.TrimPrefix(word, "createtime:")}})
				orQueries = append(orQueries, bson.M{"createtime": &primitive.Regex{Pattern: strings.TrimPrefix(word, "createtime:")}})
			}
		} else if strings.HasPrefix(word, "status:") {
			orQueries = append(orQueries, bson.M{"status": &primitive.Regex{Pattern: strings.TrimPrefix(word, "status:")}})
		} else if strings.HasPrefix(word, "itemstatus:") {
			orQueries = append(orQueries, bson.M{"itemstatus": &primitive.Regex{Pattern: strings.TrimPrefix(word, "itemstatus:")}})
		} else if strings.HasPrefix(word, "project:") {
			orQueries = append(orQueries, bson.M{"project": &primitive.Regex{Pattern: strings.TrimPrefix(word, "project:")}})
		} else if strings.HasPrefix(word, "name:") {
			orQueries = append(orQueries, bson.M{"name": &primitive.Regex{Pattern: strings.TrimPrefix(word, "name:")}})
		} else if strings.HasPrefix(word, "task:") {
			orQueries = append(orQueries, bson.M{"task": &primitive.Regex{Pattern: strings.TrimPrefix(word, "task:")}})
		} else {
			orQueries = append(orQueries, bson.M{"project": &primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"name": &primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"task": &primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"updatetime": &primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"author": &primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"path": &primitive.Regex{Pattern: word}})
			orQueries = append(orQueries, bson.M{"description": &primitive.Regex{Pattern: word}})
		}
		if len(orQueries) == 0 {
			return results, nil
		}
		allQueries = append(allQueries, bson.M{"$or": orQueries})
	}
	q := bson.M{"$and": allQueries}
	findOptions := options.Find().SetSort(bson.D{{Key: "createtime", Value: -1}})
	cur, err := collection.Find(ctx, q, findOptions)
	if err != nil {
		return results, err
	}
	for cur.Next(ctx) {
		var result Review
		err := cur.Decode(&result)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}

// RmProjectReviewV2 함수는 해당 프로젝트의 Review 데이터를 DB에서 삭제한다.
func RmProjectReviewV2(client *mongo.Client, project string) error {
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.DeleteMany(ctx, bson.M{"project": &primitive.Regex{Pattern: project}})
	if err != nil {
		return err
	}
	return nil
}

func addReviewV2(client *mongo.Client, r Review) error {
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"_id": r.ID})
	if err != nil {
		return err
	}

	if num != 0 {
		return errors.New(r.ID.Hex() + " ID를 가진 문서가 이미 DB에 존재합니다.")
	}
	r.Createtime = time.Now().Format(time.RFC3339)
	_, err = collection.InsertOne(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

func getReviewV2(client *mongo.Client, id string) (Review, error) {
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := Review{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return r, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func setReviewProcessStatusV2(client *mongo.Client, id, status string) error {
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(ctx, objID, bson.M{"$set": bson.M{"processstatus": status}})
	if err != nil {
		return err
	}
	return nil
}

func setErrReviewV2(client *mongo.Client, id, log string) error {
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(ctx, objID, bson.M{"$set": bson.M{"processstatus": "error", "log": log}})
	if err != nil {
		return err
	}
	return nil
}

func setReviewPathV2(client *mongo.Client, id, path string) error {
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(ctx, objID, bson.M{"$set": bson.M{"path": path}})
	if err != nil {
		return err
	}
	return nil
}

func setReviewItemV2(client *mongo.Client, r Review) error {
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.D{{Key: "$set", Value: r}}

	_, err := collection.UpdateByID(ctx, r.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func setReviewItemStatusV2(client *mongo.Client, id, itemstatus string) error {
	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"itemstatus": itemstatus}}

	_, err = collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return err
	}
	return nil
}

func addReviewCommentV2(client *mongo.Client, id string, cmt Comment) error {
	if cmt.Text == "" {
		return errors.New("need comment")
	}

	collection := client.Database(*flagDBName).Collection("review")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 먼저 문서를 조회하여 'comments' 필드 상태 확인
	var result bson.M
	if err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result); err != nil {
		return err
	}

	// 'comments' 필드가 null이라면 빈 배열로 초기화하는 업데이트 실행
	if result["comments"] == nil {
		update := bson.M{"$set": bson.M{"comments": []Comment{cmt}}}
		if _, err := collection.UpdateByID(ctx, objID, update); err != nil {
			return err
		}
	} else {
		// 'comments' 필드가 null이 아니라면 댓글 추가
		update := bson.M{"$push": bson.M{"comments": cmt}}
		if _, err := collection.UpdateByID(ctx, objID, update); err != nil {
			return err
		}
	}
	return nil
}
