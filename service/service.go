package service

import (
	"errors"
	"fmt"

	cache "github.com/psuman/go-training/service/cache"
	common "github.com/psuman/go-training/service/common"

	persistence "github.com/psuman/go-training/service/persistence"
)

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
		svc.CacheFinder.PutItemInCache(prodID, itemFromDb)
		return itemFromDb, nil
	}

	return itemFromCache, nil

}
