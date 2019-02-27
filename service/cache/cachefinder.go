package cache

import common "product-catalog-service/service/common"

//CacheFinder retrieves item with productId from cache
type CacheFinder interface {
	FindItemInCache(ProductID string) (common.ProductDetails, error)
}

//RedisCacheFinder retrieves item with productId from redis cache
type RedisCacheFinder struct{}

//FindItemInCache retrieves item from redis cache
func (RedisCacheFinder) FindItemInCache(ProductID string) (common.ProductDetails, error) {
	if ProductID == "a123" {
		return common.ProductDetails{ProdID: "a123", ProdName: "iPhone", ProdDesc: "new iPhone", Quantity: 10}, nil
	} else {
		return common.ProductDetails{}, nil
	}
}
