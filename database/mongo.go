// Package database handles mongo database calls and functionality
package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	openai "github.com/sashabaranov/go-openai"
)

var client *mongo.Database

func InitMongo(ctx context.Context) (*mongo.Database, error) {
	if client != nil {
		return client, nil
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI not set")
	}

  db, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect error %w", err)
	} 

	if err := db.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	client = db.Database("ai-final")
	return client, nil
}

type ConversationListItem struct {
	ID string
	Title string
}

func GetConversationsByID(ctx context.Context, userID string) ([]ConversationListItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	InitMongo(ctx)

	coll := client.Collection("conversations")

	cursor, err := coll.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, fmt.Errorf("find error: %w", err)
	}
	defer cursor.Close(ctx)

	var items []ConversationListItem
	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}

		if oid, ok := raw["_id"].(bson.ObjectID); ok {
			items = append(items, ConversationListItem{
				ID: oid.Hex(),
				Title: fmt.Sprint(raw["title"]),
			})
		}
	}
	return items, nil
}

func AddNewConversation(title, userID string, messages []openai.ChatCompletionMessage) (string, error) {
	coll := client.Collection("conversations")

	result, err := coll.InsertOne(context.TODO(), bson.M{
		"userId": userID,
		"title": title,
		"messages": messages,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create conversation: %w", err)
	}

	idRaw, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("could not parse inserted user _id")
	}

	return idRaw.Hex(), nil
}

func AddMessageToConversation(ctx context.Context, title, chatID, userID string, userMessage, botMessage openai.ChatCompletionMessage) error {
	messages := []openai.ChatCompletionMessage{
		userMessage,
		botMessage,
	}
	oid, err := bson.ObjectIDFromHex(chatID)
	if err != nil {
		return fmt.Errorf("error parsing chatID")
	}
	_, err = client.Collection("conversations").UpdateOne(ctx, 
		bson.M{"_id": oid, "userId": userID},
		bson.M{"$push": bson.M{"messages": bson.M{"$each": messages}}},
		)
	if err != nil {
		return fmt.Errorf("error adding message to conversation: %w", err)
	}
	return nil
}

type StoredMessage struct {
	Role string `bson:"role" json:"role"`
	Content string `bson:"content" json:"content"`
}

type ConversationDoc struct {
	ID bson.ObjectID `bson:"_id"`
	UserID string `bson:"userId"`
	Messages []StoredMessage `bson:"messages"`
}

func GetConversation(ctx context.Context, userID, chatID string) ([]openai.ChatCompletionMessage, error) {
	var convo ConversationDoc
	oid, err := bson.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, fmt.Errorf("error parsing chatID")
	}
	err = client.Collection("conversations").FindOne(ctx, bson.M{
		"_id": oid,
		"userId": userID,
	}).Decode(&convo)	
	if err != nil {
		return nil, fmt.Errorf("error getting conversation")
	}

	converted := make([]openai.ChatCompletionMessage, 0, len(convo.Messages))
	for _, msg := range convo.Messages {
		converted = append(converted, openai.ChatCompletionMessage{
			Role: msg.Role,
			Content: msg.Content,
		})
	}
	return converted, nil
}
