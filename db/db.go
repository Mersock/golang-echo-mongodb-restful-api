package db

import (
	"context"
	"fmt"
	"log"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(cfg config.Properties) *mongo.Database {
	connURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort)

	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connURI))

	if err != nil {
		log.Fatalf("Unable to conntect to database : %s", err)
	}

	db := c.Database(cfg.DBName)

	return db
}

func NewTest(cfg config.Properties) *mongo.Database {
	connURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort)

	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connURI))

	if err != nil {
		log.Fatalf("Unable to conntect to database : %s", err)
	}

	db := c.Database(cfg.DBTestName)

	return db
}
