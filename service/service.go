package service

import (
	"errors"

	"github.com/go-kit/kit/log"

	cache "github.com/psuman/go-training/service/cache"
	common "github.com/psuman/go-training/service/common"

	persistence "github.com/psuman/go-training/service/persistence"
)

// ErrEmpty thrown when productId is empty
var ErrEmpty = errors.New("empty Product ID")

// ItemService finds item service retreives item with given product id
// When failed to retrieve item it will return an error
type ItemService interface {
	FindItem(findItemRequest) findItemResponse
	AddItem(addItemRequest) addItemResponse
	Close()
}

// ItemCatalogService is the implementation of FindItemService
type ItemCatalogService struct {
	CacheFinder cache.CacheFinder
	ItemDao     persistence.ItemDao
	Logger      log.Logger
}

func (svc ItemCatalogService) AddItem(req addItemRequest) addItemResponse {

	id, err := svc.ItemDao.AddItem(req.ProdDetails)

	if err != nil {
		return addItemResponse{Err: err.Error()}
	}

	return addItemResponse{Id: id}
}

// FindItem retrieves item details from redis cache if exists. If not loads item from mongo and cache it in Redis
// and return item details as response
func (svc ItemCatalogService) FindItem(req findItemRequest) findItemResponse {
	if req.ProdID == "" {
		return findItemResponse{Err: "ProductId is empty"}
	}

	var itemFromCache common.ProductDetails

	itemFromCache, _ = svc.CacheFinder.FindItemInCache(req.ProdID)

	if itemFromCache.ProdID == "" {
		itemFromDb, err := svc.ItemDao.FindItem(req.ProdID)
		if err != nil {
			return findItemResponse{Err: "Product not found"}
		}

		svc.CacheFinder.PutItemInCache(req.ProdID, itemFromDb)

		return findItemResponse{ProdDetails: itemFromDb}
	}

	return findItemResponse{ProdDetails: itemFromCache}

}

func (svc ItemCatalogService) Close() {
	err := svc.CacheFinder.Close()
	if err != nil {
		svc.Logger.Log("failed to close cache connection: [error=%s]", err.Error())
	}
	err = svc.ItemDao.Close()
	if err != nil {
		svc.Logger.Log("failed to close mongo connection: [error=%s]", err.Error())
	}
}
