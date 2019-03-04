package cache

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/go-redis/redis"
	common "github.com/psuman/go-training/service/common"
)

//CacheFinder retrieves item with productId from cache
type CacheFinder interface {
	FindItemInCache(ProductID string) (common.ProductDetails, error)
	PutItemInCache(ProductId string, ProductDetails common.ProductDetails) error
	Close() error
}

//RedisCacheFinder retrieves item with productId from redis cache
type RedisCacheFinder struct {
	redisClient *redis.Client
	logger      log.Logger
}

func (cacheFinder RedisCacheFinder) Initialize(connUrl string, logger log.Logger) RedisCacheFinder {
	client := redis.NewClient(&redis.Options{
		Addr:     connUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	cacheFinder.redisClient = client
	cacheFinder.logger = logger
	return cacheFinder
}

func (cacheFinder RedisCacheFinder) Close() error {
	fmt.Println("Inside redis cache Close")

	err := cacheFinder.redisClient.Close()
	if err != nil {
		return err
	}

	return nil
}

//FindItemInCache retrieves item from redis cache
func (cacheFinder RedisCacheFinder) FindItemInCache(productID string) (common.ProductDetails, error) {
	fmt.Println("Inside redis cache finder")
	val, _ := cacheFinder.redisClient.Get(productID).Result()
	fmt.Printf("val inside cache finder:%s", val)
	if val == "" {
		return common.ProductDetails{}, errors.New("Missing Key")
	}

	var productDetails common.ProductDetails

	err := json.Unmarshal([]byte(val), &productDetails)

	if err != nil {
		return common.ProductDetails{}, nil
	}

	return productDetails, nil

}

//FindItemInCache retrieves item from redis cache
func (cacheFinder RedisCacheFinder) PutItemInCache(productID string, productDetails common.ProductDetails) error {
	res, err := json.Marshal(productDetails)
	if err != nil {
		return err
	}

	err = cacheFinder.redisClient.Set(productID, res, 0).Err()

	if err != nil {
		return err
	}

	return nil

	// if ProductID == "a123" {
	// 	return common.ProductDetails{ProdID: "a123", ProdName: "iPhone", ProdDesc: "new iPhone", Quantity: 10}, nil
	// } else {
	// 	return common.ProductDetails{}, nil
	// }
}
