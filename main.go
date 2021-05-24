package main

import (
	"fmt"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"github.com/Mersock/golang-echo-mongodb-restful-api/db"
	"github.com/Mersock/golang-echo-mongodb-restful-api/handlers"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mainDB        *mongo.Database
	productCol    *mongo.Collection
	cfg           config.Properties
	CorrelationID = "X-Correlation-ID"
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read : %v", err)
	}
	mainDB = db.New(cfg)
	productCol = mainDB.Collection(cfg.CollectionName)
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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${remote_ip} ${header:X-Correlation-ID} ${host} ${method} ${uri} ${user_agent} ` +
			`${status} ${error} ${latency_human}` + "\n",
	}))

	h := handlers.ProductHandlers{Col: productCol}
	e.GET("/products/:id", h.GetProduct)
	e.PUT("/products/:id", h.UpdateProducts, middleware.BodyLimit("1M"))
	e.POST("/products", h.CreateProducts, middleware.BodyLimit("1M"))
	e.GET("/products", h.GetProducts)
	e.DELETE("/products/:id", h.DeleteProduct)

	e.Logger.Infof("Listening on %s", cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
