package dbiface

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	//collectionAPI
	CollectionAPI interface {
		InsertOne(ctx context.Context, document interface{},
			opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
		Find(ctx context.Context, filter interface{},
			opts ...*options.FindOptions) (*mongo.Cursor, error)
	}
)
