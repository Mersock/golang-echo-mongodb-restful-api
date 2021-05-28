package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestProduct(t *testing.T) {
	var docID string

	t.Run("test create product", func(t *testing.T) {
		var IDs []string
		body := `[{
			"product_name": "test",
			"price": 33,
			"currency": "THB",
			"vendor": "Google",
			"accessories": [
				"gift coupon"
			]
		}]`
		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		h.Col = productCol
		err := h.CreateProducts(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)

		err = json.Unmarshal(res.Body.Bytes(), &IDs)
		assert.Nil(t, err)
		docID = IDs[0]
		t.Logf("IDs: %#+v\n", IDs)
		for _, ID := range IDs {
			assert.NotNil(t, ID)
		}
	})

	t.Run("get product", func(t *testing.T) {
		var products []Product
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		h.Col = productCol
		test := h.GetProducts(c)
		assert.Nil(t, test)
		assert.Equal(t, http.StatusOK, res.Code)

		err := json.Unmarshal(res.Body.Bytes(), &products)
		assert.Nil(t, err)
		for _, product := range products {
			assert.Equal(t, "test", product.Name)
		}
	})

	t.Run("get product with query params", func(t *testing.T) {
		var products []Product
		req := httptest.NewRequest(http.MethodGet, "/products?currency=THB&vendor=google", nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		h.Col = productCol
		test := h.GetProducts(c)
		assert.Nil(t, test)
		assert.Equal(t, http.StatusOK, res.Code)

		err := json.Unmarshal(res.Body.Bytes(), &products)
		assert.Nil(t, err)
		for _, product := range products {
			assert.Equal(t, "test", product.Name)
		}
	})

	t.Run("get single product", func(t *testing.T) {
		var product Product
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products/%s", docID), nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		c.SetParamNames("id")
		c.SetParamValues(docID)
		h.Col = productCol
		test := h.GetProduct(c)
		assert.Nil(t, test)
		assert.Equal(t, http.StatusOK, res.Code)

		err := json.Unmarshal(res.Body.Bytes(), &product)
		assert.Nil(t, err)
		assert.Equal(t, "THB", product.Currency)
	})

	t.Run("put product", func(t *testing.T) {
		var product Product
		body := `
		{
			"product_name": "test",
			"price": 250,
			"currency": "THB",
			"vendor": "LINE",
			"accessories": [
				"gift","subscription"
			]
		}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/products/%s", docID), strings.NewReader(body))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		c.SetParamNames("id")
		c.SetParamValues(docID)
		h.Col = productCol
		test := h.UpdateProducts(c)
		assert.Nil(t, test)
		assert.Equal(t, http.StatusOK, res.Code)

		err := json.Unmarshal(res.Body.Bytes(), &product)
		assert.Nil(t, err)
		assert.Equal(t, "THB", product.Currency)
	})

	t.Run("delete single product", func(t *testing.T) {
		var delCount int64
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/products/%s", docID), nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		c.SetParamNames("id")
		c.SetParamValues(docID)
		h.Col = productCol
		test := h.DeleteProduct(c)
		assert.Nil(t, test)
		assert.Equal(t, http.StatusOK, res.Code)

		err := json.Unmarshal(res.Body.Bytes(), &delCount)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), delCount)
	})
}
