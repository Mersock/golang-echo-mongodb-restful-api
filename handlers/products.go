package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/Mersock/golang-echo-mongodb-restful-api/dbiface"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	v = validator.New()
)

type ProductValidator struct {
	validator *validator.Validate
}

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"product_name" bson:"product_name" validate:"required,max=10"`
	Price       int                `json:"price" bson:"price" validate:"required,max=2000"`
	Currency    string             `json:"currency" bson:"currency" validate:"required,len=3"`
	Discount    int                `json:"discount" bson:"discount"`
	Vendor      string             `json:"vendor" bson:"vendor" validate:"required"`
	Accessories []string           `json:"accessories,omitempty" bson:"accessories,omitempty"`
	IsEssential bool               `json:"is_essential" bson:"is_essential"`
}

type ProductHandlers struct {
	Col dbiface.CollectionAPI
}

func (p *ProductValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}

func findProducts(ctx context.Context, q url.Values, collection dbiface.CollectionAPI) ([]Product, error) {
	var products []Product
	filter := make(map[string]interface{})

	for k, v := range q {
		filter[k] = v[0]
	}

	if id, ok := filter["_id"]; ok {
		docID, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return products, err
		}
		filter["_id"] = docID
	}

	cursor, err := collection.Find(ctx, bson.M(filter))

	if err != nil {
		log.Errorf("Unable to find the product :%v", err)
		return products, err
	}

	err = cursor.All(ctx, &products)
	if err != nil {
		log.Errorf("Unable to read the cursor :%v", err)
		return products, err
	}

	return products, nil
}

func (h *ProductHandlers) GetProducts(c echo.Context) error {
	products, err := findProducts(context.Background(), c.QueryParams(), h.Col)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, products)
}

func insertProducts(ctx context.Context, products []Product, collection dbiface.CollectionAPI) ([]interface{}, error) {
	var insertIds []interface{}
	for _, product := range products {
		product.ID = primitive.NewObjectID()
		insertID, err := collection.InsertOne(ctx, product)
		if err != nil {
			log.Errorf("Unable to insert :%v", err)
			return nil, err
		}
		insertIds = append(insertIds, insertID.InsertedID)
	}
	return insertIds, nil
}

func (h *ProductHandlers) CreateProducts(c echo.Context) error {
	var products []Product
	c.Echo().Validator = &ProductValidator{validator: v}

	if err := c.Bind(&products); err != nil {
		log.Errorf("Unable to bind: %v", err)
		return err
	}

	for _, product := range products {
		if err := c.Validate(product); err != nil {
			log.Errorf("Unable to validate the product %v %+v", product, err)
			return err
		}

	}

	IDs, err := insertProducts(context.Background(), products, h.Col)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, IDs)
}

func modifyProduct(ctx context.Context, id string, reqBody io.ReadCloser, collection dbiface.CollectionAPI) (Product, error) {
	var product Product
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return product, err
	}
	filter := bson.M{"_id": docId}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(&product); err != nil {
		return product, err
	}

	if err := json.NewDecoder(reqBody).Decode(&product); err != nil {
		return product, err
	}

	if err := v.Struct(product); err != nil {
		return product, err
	}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": product})

	if err != nil {
		return product, err
	}

	return product, nil
}

func (h *ProductHandlers) UpdateProducts(c echo.Context) error {
	var product Product
	product, err := modifyProduct(context.Background(), c.Param("id"), c.Request().Body, h.Col)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, product)
}
