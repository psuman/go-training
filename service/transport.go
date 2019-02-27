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
func MakeFindItemEndPoint(svc FindItemService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(findItemRequest)
		res, err := svc.FindItem(req.ProdID)

		if err != nil {
			return common.ProductDetails{}, err
		}

		return res, nil
	}
}

//DecodeFindItemRequest decodes find item request
func DecodeFindItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request findItemRequest
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
	ProdID string `json:"productId"`
}

type findItemResponse struct {
	ProdDetails common.ProductDetails `json:"productDetails"`
	Err         string                `json:"err,omitempty"`
}
