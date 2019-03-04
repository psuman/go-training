package service

import (
	"errors"
	"time"

	"github.com/go-kit/kit/log"

	cache "github.com/psuman/go-training/service/cache"
	common "github.com/psuman/go-training/service/common"

	external_invoker "github.com/psuman/go-training/service/external_invoker"

	persistence "github.com/psuman/go-training/service/persistence"

	"github.com/go-kit/kit/metrics"
)

// ErrEmpty thrown when productId is empty
var ErrEmpty = errors.New("empty Product ID")

type MetricsMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Next           ItemService
}

func (mw MetricsMiddleware) FindItem(req findItemRequest) findItemResponse {
	defer func(begin time.Time) {
		lvs := []string{"method", "FindItem", "error", "false"}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res := mw.Next.FindItem(req)
	return res
}

func (mw MetricsMiddleware) AddItem(req addItemRequest) addItemResponse {
	defer func(begin time.Time) {
		lvs := []string{"method", "AddItem", "error", "false"}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res := mw.Next.AddItem(req)
	return res
}

func (mw MetricsMiddleware) Close() {
	mw.Next.Close()
}

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
	ExtService  external_invoker.ExternalFindItemServiceInvoker
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
	svc.Logger.Log("req", req.ProdID)
	if req.ProdID == "" {
		return findItemResponse{Err: "ProductId is empty"}
	}

	var itemFromCache common.ProductDetails

	itemFromCache, err := svc.CacheFinder.FindItemInCache(req.ProdID)

	if err != nil {
		itemFromDb, err := svc.ItemDao.FindItem(req.ProdID)
		svc.Logger.Log("ITEM_LOADED_FROM_DB", itemFromDb.ProdID)

		if err != nil {
			extReq := external_invoker.ExternalFindItemRequest{ProdID: req.ProdID}

			extRes, err := svc.ExtService.Invoke(extReq)

			if err != nil {
				return findItemResponse{Err: "Product not found"}
			}

			return findItemResponse{ProdDetails: extRes.ProdDetails}
		}

		svc.CacheFinder.PutItemInCache(req.ProdID, itemFromDb)
		return findItemResponse{ProdDetails: itemFromDb}
	}

	svc.Logger.Log("ITEM_FOUND_IN_CACHE", itemFromCache.ProdID)

	return findItemResponse{ProdDetails: itemFromCache}

}

// Close closes cache and database connections
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
