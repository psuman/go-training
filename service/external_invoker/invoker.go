package external_invoker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	common "github.com/psuman/go-training/service/common"
)

type ExternalFindItemRequest struct {
	ProdID string `json:"productId"`
}

type ExternalFindItemResponse struct {
	ProdDetails common.ProductDetails `json:"productDetails"`
	Err         string                `json:"err,omitempty"`
}

type ExternalFindItemServiceInvoker interface {
	Invoke(req ExternalFindItemRequest) (ExternalFindItemResponse, error)
}

type ExternalFindItemServiceInvokerImpl struct {
	serviceUrl string
	httpClient *http.Client
}

func (invoker ExternalFindItemServiceInvokerImpl) Initialize(serviceUrl string) ExternalFindItemServiceInvoker {
	httpClient := &http.Client{}
	return ExternalFindItemServiceInvokerImpl{serviceUrl: serviceUrl, httpClient: httpClient}
}

func (invoker ExternalFindItemServiceInvokerImpl) Invoke(req ExternalFindItemRequest) (ExternalFindItemResponse, error) {

	reqPayload, err := json.Marshal(req)

	if err != nil {
		return ExternalFindItemResponse{}, err
	}

	httpReq, err := http.NewRequest("POST", invoker.serviceUrl, bytes.NewBuffer([]byte(reqPayload)))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := invoker.httpClient.Do(httpReq)

	if err != nil {
		return ExternalFindItemResponse{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var responseObj ExternalFindItemResponse

	if err := json.Unmarshal([]byte(body), &responseObj); err != nil {
		return ExternalFindItemResponse{}, err
	}

	return responseObj, nil
}
