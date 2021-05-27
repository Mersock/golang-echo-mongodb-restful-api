package main

import (
	"context"
	"fmt"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"github.com/Mersock/golang-echo-mongodb-restful-api/db"
	"github.com/Mersock/golang-echo-mongodb-restful-api/handlers"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mainDB        *mongo.Database
	productsCol   *mongo.Collection
	usersCol      *mongo.Collection
	cfg           config.Properties
	CorrelationID = "X-Correlation-ID"
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read : %v", err)
	}
	mainDB = db.New(cfg)
	productsCol = mainDB.Collection(cfg.ProductCollection)
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

func addCorrelationID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Request().Header.Get(CorrelationID)
		var newID string
		if id == "" {
			newID = random.String(12)

		} else {
			newID = id
		}
		c.Request().Header.Set(CorrelationID, newID)
		c.Response().Header().Set(CorrelationID, newID)
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(addCorrelationID)
	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(cfg.JwtTokenSecret),
		TokenLookup: "header:x-auth-token",
	})
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${remote_ip} ${header:X-Correlation-ID} ${host} ${method} ${uri} ${user_agent} ` +
			`${status} ${error} ${latency_human}` + "\n",
	}))

	h := &handlers.ProductHandlers{Col: productsCol}
	e.GET("/products/:id", h.GetProduct)
	e.PUT("/products/:id", h.UpdateProducts, middleware.BodyLimit("1M"), jwtMiddleware)
	e.POST("/products", h.CreateProducts, middleware.BodyLimit("1M"), jwtMiddleware)
	e.GET("/products", h.GetProducts)
	e.DELETE("/products/:id", h.DeleteProduct, jwtMiddleware)

	uh := handlers.UsersHandler{Col: usersCol}
	e.POST("/users", uh.CreateUser)
	e.POST("/auth", uh.AuthUser)

	e.Logger.Infof("Listening on %s", cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
