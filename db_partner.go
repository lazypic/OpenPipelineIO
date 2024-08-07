package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addPartner(client *mongo.Client, s Partner) error {
	if s.Name == "" {
		return errors.New("need name")
	}
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	num, err := collection.CountDocuments(ctx, bson.M{"name": s.Name})
	if err != nil {
		return err
	}
	if num != 0 {
		return errors.New("A name with the same name already exists, so it cannot be created")
	}
	_, err = collection.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func getPartner(client *mongo.Client, id string) (Partner, error) {
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s := Partner{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return s, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func rmPartner(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}

func setPartner(client *mongo.Client, s Partner) error {
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": s.ID},
		bson.D{{Key: "$set", Value: s}},
	)
	if err != nil {
		return err
	}
	return nil
}

func allPartners(client *mongo.Client) ([]Partner, error) {
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []Partner
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func allPartnersCodename(client *mongo.Client) ([]string, error) {
	collection := client.Database(*flagDBName).Collection("partner")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var results []string
	values, err := collection.Distinct(ctx, "codename", bson.D{})
	if err != nil {
		return results, err
	}
	for _, value := range values {
		if fmt.Sprintf("%v", value) == "" {
			continue
		}
		results = append(results, fmt.Sprintf("%v", value))
	}
	sort.Strings(results)
	return results, nil
}
