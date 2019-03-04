package ext_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/go-kit/kit/log"
)

// ErrEmpty thrown when productId is empty
var ErrEmpty = errors.New("empty Product ID")

type ProductDetails struct {
	ProdID   string
	Name     string
	Desc     string
	Quantity int32
}

// ItemService finds item service retreives item with given product id
// When failed to retrieve item it will return an error
type ItemService interface {
	FindItem(extFindItemRequest) extFindItemResponse
}

// ItemCatalogService is the implementation of FindItemService
type ItemCatalogService struct {
	Logger log.Logger
}

// FindItem retrieves item details from redis cache if exists. If not loads item from mongo and cache it in Redis
// and return item details as response
func (svc ItemCatalogService) FindItem(req extFindItemRequest) extFindItemResponse {
	svc.Logger.Log("Request", req.ProdID)
	prodDetails := ProductDetails{ProdID: "ext123", Name: "extName", Desc: "extDesc", Quantity: 1}
	return extFindItemResponse{ProdDetails: prodDetails}
}

//MakeFindItemEndPoint creates end point for find item
func MakeFindItemEndPoint(svc ItemService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(extFindItemRequest)
		return svc.FindItem(req), nil
	}
}

//DecodeFindItemRequest decodes find item request
func DecodeFindItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request extFindItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Printf("error inside transport %s", err.Error())
		return nil, err
	}

	return request, nil
}

//EncodeResponse encodes response
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type extFindItemRequest struct {
	ProdID string `json:"ProdId"`
}

type extFindItemResponse struct {
	ProdDetails ProductDetails `json:"productDetails"`
	Err         string         `json:"err,omitempty"`
}
