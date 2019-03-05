package service

import (
	"errors"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	external_invoker "github.com/psuman/go-training/service/external_invoker"

	common "github.com/psuman/go-training/service/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCacheFinder struct {
	mock.Mock
}

func (mock *MockCacheFinder) FindItemInCache(ProductID string) (common.ProductDetails, error) {
	args := mock.Called(ProductID)
	prodDetails, _ := args[0].(common.ProductDetails)
	err, _ := args[1].(error)

	return prodDetails, err
}

func (cf MockCacheFinder) PutItemInCache(ProductId string, ProductDetails common.ProductDetails) error {
	return nil
}

func (cf MockCacheFinder) Close() error {
	return nil

}

type MockItemDao struct {
	mock.Mock
}

func (mock *MockItemDao) FindItem(ProductID string) (common.ProductDetails, error) {
	args := mock.Called(ProductID)
	prodDetails, _ := args[0].(common.ProductDetails)
	err, _ := args[1].(error)

	return prodDetails, err

}

func (mock *MockItemDao) AddItem(productDetails common.ProductDetails) (string, error) {
	return "", nil

}

func (mock *MockItemDao) Close() error {
	return nil
}

type MockExtService struct {
	mock.Mock
}

func (mock *MockExtService) Invoke(req external_invoker.ExternalFindItemRequest) (external_invoker.ExternalFindItemResponse, error) {

	args := mock.Called(req)
	svcResp, _ := args[0].(external_invoker.ExternalFindItemResponse)
	err, _ := args[1].(error)

	return svcResp, err
}

func TestFindItem_WehnNotFoundInCache_ShouldLoadFromDB(t *testing.T) {

	cacheFinder := new(MockCacheFinder)
	itemDao := new(MockItemDao)
	extSvc := new(MockExtService)

	cacheFinder.On("FindItemInCache", "a123").Return(common.ProductDetails{}, errors.New("not in cache"))
	itemDao.On("FindItem", "a123").Return(common.ProductDetails{ProdID: "a123", Name: "a1", Desc: "desc", Quantity: 1}, nil)

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	svc := ItemCatalogService{CacheFinder: cacheFinder,
		ItemDao: itemDao, ExtService: extSvc, Logger: logger}

	res := svc.FindItem(findItemRequest{ProdID: "a123"})

	assert.Equal(t, "a123", res.ProdDetails.ProdID)
	cacheFinder.AssertExpectations(t)
	itemDao.AssertExpectations(t)
	extSvc.AssertExpectations(t)

}

func TestFindItem_WehnFoundInCache_ShouldReturnFromCache(t *testing.T) {

	cacheFinder := new(MockCacheFinder)
	itemDao := new(MockItemDao)
	extSvc := new(MockExtService)

	cacheFinder.On("FindItemInCache", "a123").Return(common.ProductDetails{ProdID: "a123", Name: "a1", Desc: "desc", Quantity: 1}, nil)

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	svc := ItemCatalogService{CacheFinder: cacheFinder,
		ItemDao: itemDao, ExtService: extSvc, Logger: logger}

	res := svc.FindItem(findItemRequest{ProdID: "a123"})

	assert.Equal(t, "a123", res.ProdDetails.ProdID)
	cacheFinder.AssertExpectations(t)
	itemDao.AssertExpectations(t)
	extSvc.AssertExpectations(t)

}

func TestFindItem_WehnNotInCacheAndNotInDB_ShouldReturnFromExternalService(t *testing.T) {

	cacheFinder := new(MockCacheFinder)
	itemDao := new(MockItemDao)
	extSvc := new(MockExtService)

	extSvcResp := external_invoker.ExternalFindItemResponse{ProdDetails: common.ProductDetails{ProdID: "ext123", Name: "a1", Desc: "desc", Quantity: 1}, Err: ""}

	cacheFinder.On("FindItemInCache", "ext123").Return(common.ProductDetails{}, errors.New("not in cache"))
	itemDao.On("FindItem", "ext123").Return(common.ProductDetails{}, errors.New("not in db"))
	extSvc.On("Invoke", mock.Anything).Return(extSvcResp, "")

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	svc := ItemCatalogService{CacheFinder: cacheFinder,
		ItemDao: itemDao, ExtService: extSvc, Logger: logger}

	res := svc.FindItem(findItemRequest{ProdID: "ext123"})

	assert.Equal(t, "ext123", res.ProdDetails.ProdID)
	cacheFinder.AssertExpectations(t)
	itemDao.AssertExpectations(t)
	extSvc.AssertExpectations(t)

}
