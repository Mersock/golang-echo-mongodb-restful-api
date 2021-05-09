package handlers

import (
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"github.com/Mersock/golang-echo-mongodb-restful-api/db"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	cfg config.Properties
	h   ProductHandlers
	col *mongo.Collection
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration cannot be read : %v", err)
	}
	col = db.New(cfg)
}

func TestProduct(t *testing.T) {
	t.Run("test create product", func(t *testing.T) {
		body := `[{
			"product_name": "test",
			"price": 33,
			"currency": "THB",
			"vendor": "Google",
			"accessories": [
				"gift coupon"
			]
		}]`
		req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		h.Col = col
		test := h.CreateProducts(c)
		assert.Nil(t, test)
	})
}
