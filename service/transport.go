package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	common "github.com/psuman/go-training/service/common"

	"github.com/go-kit/kit/endpoint"
)

//MakeFindItemEndPoint creates end point for find item
func MakeFindItemEndPoint(svc ItemService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(findItemRequest)
		return svc.FindItem(req), nil
	}
}

//MakeAddItemEndPoint creates end point for add item
func MakeAddItemEndPoint(svc ItemService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addItemRequest)
		return svc.AddItem(req), nil
	}
}

//DecodeFindItemRequest decodes find item request
func DecodeFindItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request findItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Printf("error inside transport %s", err.Error())
		return nil, err
	}
	fmt.Println(request.ProdID)

	return request, nil
}

//DecodeAddItemRequest decodes add item request
func DecodeAddItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request addItemRequest
	fmt.Println("inside transport")
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

type findItemRequest struct {
	ProdID string `json:"ProdId"`
}

type findItemResponse struct {
	ProdDetails common.ProductDetails `json:"productDetails"`
	Err         string                `json:"err,omitempty"`
}

type addItemRequest struct {
	ProdDetails common.ProductDetails `json:"productDetails"`
}

type addItemResponse struct {
	Id  string `json:"id"`
	Err string `json:"err,omitempty"`
}
