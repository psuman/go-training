package external_invoker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	common "github.com/psuman/go-training/service/common"
)

type ExternalFindItemRequest struct {
	ProdID string
}

type ExternalFindItemResponse struct {
	ProdDetails common.ProductDetails `json:"productDetails"`
	Err         string                `json:"err,omitempty"`
}

type ExternalFindItemServiceInvoker interface {
	Invoke(req ExternalFindItemRequest) (ExternalFindItemResponse, error)
}

type ExternalFindItemServiceInvokerImpl struct {
	ServiceUrl string
	Timeout    int32
	HttpClient *http.Client
}

func (invoker ExternalFindItemServiceInvokerImpl) Invoke(req ExternalFindItemRequest) (ExternalFindItemResponse, error) {
	var jsonStr = fmt.Sprintf("{'productId':'%s'}", req.ProdID)
	httpReq, err := http.NewRequest("POST", invoker.ServiceUrl, bytes.NewBuffer([]byte(jsonStr)))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := invoker.HttpClient.Do(httpReq)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var responseObj ExternalFindItemResponse

	if err := json.Unmarshal([]byte(body), &responseObj); err != nil {
		panic(err)
	}

	return responseObj, nil
}
