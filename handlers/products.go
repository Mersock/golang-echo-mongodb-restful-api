package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/Mersock/golang-echo-mongodb-restful-api/dbiface"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func insertProducts(ctx context.Context, products []Product, collection dbiface.CollectionAPI) ([]interface{}, error) {
	var insertIds []interface{}
	for _, product := range products {
		product.ID = primitive.NewObjectID()
		insertID, err := collection.InsertOne(ctx, product)
		if err != nil {
			log.Printf("Unable to insert :%v", err)
			return nil, err
		}
		insertIds = append(insertIds, insertID.InsertedID)
	}
	return insertIds, nil
}

func (h *ProductHandlers) CreateProducts(c echo.Context) error {
	var products []Product

	if err := c.Bind(&products); err != nil {
		log.Printf("Unable to bind: %v", err)
		return err
	}

	IDs, err := insertProducts(context.Background(), products, h.Col)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, IDs)
}
