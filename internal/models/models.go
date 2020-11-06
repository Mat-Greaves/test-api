package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDB(connectionString string) (*mongo.Client, error) {
	var err error
	db, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = db.Connect(ctx)
	if err != nil {
		return nil, err
	}
	err = db.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}
