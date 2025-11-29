// Package auth handles login/signup functionality
package auth

import (
	"ai-final/database"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func Signup(ctx context.Context, username, password string) (string, error) {
	db, err := database.InitMongo(ctx)	
	if err != nil {
		return "", fmt.Errorf("database connection error")
	} 
	
	coll := db.Collection("users")

	result, err := coll.InsertOne(context.TODO(), bson.M{
		"username": username,
		"hashedPassword": password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	idRaw, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("could not parse inserted user _id")
	}

	return idRaw.Hex(), nil
}
