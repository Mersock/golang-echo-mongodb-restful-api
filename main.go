package main

import (
	"fmt"
	"log"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"github.com/Mersock/golang-echo-mongodb-restful-api/db"
	"github.com/Mersock/golang-echo-mongodb-restful-api/handlers"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	cfg config.Properties
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read : %v", err)
	}
}

func main() {
	e := echo.New()
	col := db.New()
	h := handlers.ProductHandlers{Col: col}

	e.Pre(middleware.RemoveTrailingSlash())

	e.POST("/products", h.CreateProducts)

	e.Logger.Infof("Listening on %s", cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
