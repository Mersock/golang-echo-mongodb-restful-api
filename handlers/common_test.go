package handlers

import (
	"context"
	"os"
	"testing"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"github.com/Mersock/golang-echo-mongodb-restful-api/db"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	cfg        config.Properties
	h          ProductHandlers
	uh         UsersHandler
	productCol *mongo.Collection
	usersCol   *mongo.Collection
	mainDB     *mongo.Database
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read : %v", err)
	}
	mainDB = db.NewTest(cfg)
	productCol = mainDB.Collection(cfg.ProductCollection)
	usersCol = mainDB.Collection(cfg.UserCollection)

	isUserIndexUnique := true
	indexModel := mongo.IndexModel{
		Keys: bson.M{"username": 1},
		Options: &options.IndexOptions{
			Unique: &isUserIndexUnique,
		},
	}
	_, err := usersCol.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Unable to create an index : %+v", err)
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	testCode := m.Run()
	usersCol.Drop(ctx)
	productCol.Drop(ctx)
	os.Exit(testCode)
}
