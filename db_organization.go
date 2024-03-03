package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// allDivisionsV2 함수는 DB에서 전체 Divisions 정보를 가지고오는 함수입니다.
func allDivisionsV2(client *mongo.Client) ([]Division, error) {
	var results []Division
	collection := client.Database(*flagDBName).Collection("divisions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item Division
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// allDepartmentsV2 함수는 DB에서 전체 Departments 정보를 가지고오는 함수입니다.
func allDepartmentsV2(client *mongo.Client) ([]Department, error) {
	var results []Department
	collection := client.Database(*flagDBName).Collection("departments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item Department
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// allTeamsV2 함수는 DB에서 전체 Teams 정보를 가지고오는 함수입니다.
func allTeamsV2(client *mongo.Client) ([]Team, error) {
	var results []Team
	collection := client.Database(*flagDBName).Collection("teams")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item Team
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// allRolesV2 함수는 DB에서 전체 Role 정보를 가지고오는 함수입니다.
func allRolesV2(client *mongo.Client) ([]Role, error) {
	var results []Role
	collection := client.Database(*flagDBName).Collection("roles")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item Role
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// allPositionsV2 함수는 DB에서 전체 Position 정보를 가지고오는 함수입니다.
func allPositionsV2(client *mongo.Client) ([]Position, error) {
	var results []Position
	collection := client.Database(*flagDBName).Collection("positions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item Position
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// getDivisionV2 함수는 본부를 가지고오는 함수이다.
func getDivisionV2(client *mongo.Client, id string) (Division, error) {
	collection := client.Database(*flagDBName).Collection("divisions")
	d := Division{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&d)

	if err != nil {
		return d, err
	}
	return d, nil
}

// getDepartmentV2 함수는 부서를 가지고오는 함수이다.
func getDepartmentV2(client *mongo.Client, id string) (Department, error) {
	collection := client.Database(*flagDBName).Collection("departments")
	d := Department{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&d)

	if err != nil {
		return d, err
	}
	return d, nil
}

// getTeamV2 함수는 팀을 가지고오는 함수이다.
func getTeamV2(client *mongo.Client, id string) (Team, error) {
	collection := client.Database(*flagDBName).Collection("teams")
	t := Team{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&t)

	if err != nil {
		return t, err
	}
	return t, nil
}

// getRoleV2 함수는 역할을 가지고오는 함수이다.
func getRoleV2(client *mongo.Client, id string) (Role, error) {
	collection := client.Database(*flagDBName).Collection("roles")
	r := Role{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&r)

	if err != nil {
		return r, err
	}
	return r, nil
}

// getPositionV2 함수는 역할을 가지고오는 함수이다.
func getPositionV2(client *mongo.Client, id string) (Position, error) {
	collection := client.Database(*flagDBName).Collection("positions")
	p := Position{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&p)

	if err != nil {
		return p, err
	}
	return p, nil
}

func rmDivisionV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("divisions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func rmDepartmentV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("departments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func rmTeamV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("teams")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func rmRoleV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("roles")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func rmPositionV2(client *mongo.Client, id string) error {
	collection := client.Database(*flagDBName).Collection("positions")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
