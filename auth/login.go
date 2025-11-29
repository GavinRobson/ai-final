// Package auth handles login/signup functionality
package auth

import (
	"ai-final/database"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrNoUsersFound = errors.New("no users found")

func Login(ctx context.Context, username, password string) (string, error) {
	db, err := database.InitMongo(ctx)	
	if err != nil {
		return "", fmt.Errorf("database connection error")
	} 
	
	coll := db.Collection("users")

	var userDoc struct {
		ID bson.ObjectID `bson:"_id"`
		Username string `bson:"username"`
		HashedPassword string `bson:"hashedPassword"`
	}
	
	err = coll.FindOne(context.TODO(), bson.M{
		"username": username,
		"hashedPassword": password,
	}).Decode(&userDoc)

	if err == mongo.ErrNoDocuments {
		return "", ErrNoUsersFound
	} 

	if err != nil {
		return "", fmt.Errorf("internal error")
	}

	return userDoc.ID.Hex(), nil
}
