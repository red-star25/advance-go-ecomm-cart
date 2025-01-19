package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func DBSet() *mongo.Client {
	client, err := mongo.NewClient
}

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {}
