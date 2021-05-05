package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"github.com/Mersock/golang-echo-mongodb-restful-api/handlers"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	c   *mongo.Client
	db  *mongo.Database
	col *mongo.Collection
	cfg config.Properties
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read : %v", err)
	}

	connURI := fmt.Sprintf("mongodb://%s:%s@mongo/?authSource=admin", cfg.DBUser, cfg.DBPass)

	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connURI))

	if err != nil {
		log.Fatalf("Unable to conntect to database : %s", err)
	}

	db = c.Database(cfg.DBName)
	col = db.Collection(cfg.CollectionName)
}

func main() {
	e := echo.New()
	h := handlers.ProductHandlers{Col: col}
	e.Pre(middleware.RemoveTrailingSlash())

	e.POST("/products", h.CreateProducts)

	e.Logger.Infof("Listening on %s", cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
