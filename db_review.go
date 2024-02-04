package main

import (
	"context"
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
