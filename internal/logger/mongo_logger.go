package logger

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoLogger struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func New(mongoURI, database string) (*MongoLogger, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	collection := client.Database(database).Collection("logs")

	return &MongoLogger{
		client:     client,
		collection: collection,
	}, nil
}

func (ml *MongoLogger) Log(level, service, message string) {
	if ml == nil || ml.collection == nil {
		log.Printf("[%s] %s: %s", level, service, message)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logEntry := bson.M{
		"level":     level,
		"service":   service,
		"message":   message,
		"timestamp": time.Now(),
	}

	_, err := ml.collection.InsertOne(ctx, logEntry)
	if err != nil {
		log.Printf("Failed to write log to MongoDB: %v", err)
	}
}

func (ml *MongoLogger) Info(service, message string) {
	ml.Log("INFO", service, message)
}

func (ml *MongoLogger) Warn(service, message string) {
	ml.Log("WARN", service, message)
}

func (ml *MongoLogger) Error(service, message string) {
	ml.Log("ERROR", service, message)
}

func (ml *MongoLogger) Close() {
	if ml != nil && ml.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ml.client.Disconnect(ctx)
	}
}
