// Package auth handles login/signup functionality
package auth

import (
	"ai-final/database"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameTaken = errors.New("username taken")

func Signup(ctx context.Context, username, password string) (string, error) {
	db, err := database.InitMongo(ctx)	
	if err != nil {
		return "", fmt.Errorf("database connection error")
	} 
	
	coll := db.Collection("users")

	var existing bson.M
	err = coll.FindOne(ctx, bson.M{"username": username}).Decode(&existing)
	if err == nil {
		return "", ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("encryption error")
	}

	result, err := coll.InsertOne(context.TODO(), bson.M{
		"username": username,
		"hashedPassword": string(hashedPassword),
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
