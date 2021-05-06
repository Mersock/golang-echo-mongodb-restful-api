package db

import (
	"context"
	"fmt"
	"log"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(cfg config.Properties) *mongo.Collection {
	connURI := fmt.Sprintf("mongodb://%s:%s@mongo/?authSource=admin", cfg.DBUser, cfg.DBPass)

	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connURI))

	if err != nil {
		log.Fatalf("Unable to conntect to database : %s", err)
	}

	db := c.Database(cfg.DBName)
	col := db.Collection(cfg.CollectionName)

	return col
}

func TestDB(cfg config.Properties) *mongo.Collection {
	connURI := fmt.Sprintf("mongodb://%s:%s@mongo/?authSource=admin", cfg.DBUser, cfg.DBPass)

	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connURI))

	if err != nil {
		log.Fatalf("Unable to conntect to database : %s", err)
	}

	db := c.Database(cfg.DBTestName)
	col := db.Collection(cfg.CollectionName)

	return col
}

func DropTestDB(col mongo.Collection) error {
	if err := col.Database().Drop(context.Background()); err != nil {
		return err
	}
	return nil
}
