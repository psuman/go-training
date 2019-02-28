package persistence

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	common "github.com/psuman/go-training/service/common"
	"go.mongodb.org/mongo-driver/mongo"
)

//ItemDao retrieves item with productId from database
type ItemDao interface {
	FindItem(ProductID string) (common.ProductDetails, error)
}

//MongoItemDao retrieves item with productId from mongo database
type MongoItemDao struct{}

//NewClient creates and returns connection to MongoDb
func (MongoItemDao) NewClient(connUri string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, connUri)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}
	return client, nil

}

//FindItem retrieves item from Mongo database
func (MongoItemDao) FindItem(ProductID string) (common.ProductDetails, error) {
	return common.ProductDetails{ProdID: "a124", ProdName: "samsung", ProdDesc: "new samsung", Quantity: 5}, nil
}
