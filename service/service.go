package service

import (
	"errors"
	"fmt"
)

import common "product-catalog-service/service/common"
import cache "product-catalog-service/service/cache"
import persistence "product-catalog-service/service/persistence"

// ErrEmpty thrown when productId is empty
var ErrEmpty = errors.New("empty Product ID")

// FindItemService finds item service retreives item with given product id
// When failed to retrieve item it will return an error
type FindItemService interface {
	FindItem(string) (common.ProductDetails, error)
}

// FindItemInCatalogService is the implementation of FindItemService
type FindItemInCatalogService struct {
	CacheFinder cache.CacheFinder
	ItemDao     persistence.ItemDao
}

// FindItem retrieves item details from redis cache if exists. If not loads item from mongo and cache it in Redis
// and return item details as response
func (svc FindItemInCatalogService) FindItem(prodID string) (common.ProductDetails, error) {
	fmt.Println("inside find item in catalog service")
	if prodID == "" {
		return common.ProductDetails{}, ErrEmpty
	}

	var itemFromCache common.ProductDetails
	var itemFromDb common.ProductDetails

	itemFromCache, _ = svc.CacheFinder.FindItemInCache(prodID)

	if itemFromCache.ProdID == "" {
		itemFromDb, _ = svc.ItemDao.FindItem(prodID)
		return itemFromDb, nil
	}

	return itemFromCache, nil

}
