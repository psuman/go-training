package persistence

import common "github.com/psuman/go-training/service/common"

//ItemDao retrieves item with productId from database
type ItemDao interface {
	FindItem(ProductID string) (common.ProductDetails, error)
}

//MongoItemDao retrieves item with productId from mongo database
type MongoItemDao struct{}

//FindItem retrieves item from Mongo database
func (MongoItemDao) FindItem(ProductID string) (common.ProductDetails, error) {
	return common.ProductDetails{ProdID: "a124", ProdName: "samsung", ProdDesc: "new samsung", Quantity: 5}, nil
}
