package models

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func init() {
	_ = loadDb()
}

func loadDb() error {
	err := godotenv.Load(".env")
	if err != nil {
		godotenv.Load("../.env")
	}

	dbAddress := os.Getenv("database_url")
	usersdbName := os.Getenv("database_name")

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, err2 := mongo.Connect(ctx, options.Client().ApplyURI(dbAddress))
	db = client.Database(usersdbName)
	return err2
}

func GetDB() *mongo.Database {
	return db
}


