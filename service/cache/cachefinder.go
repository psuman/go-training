package cache

import (
	"encoding/json"
	"strings"

	"github.com/go-redis/redis"
	common "github.com/psuman/go-training/service/common"
)

//CacheFinder retrieves item with productId from cache
type CacheFinder interface {
	FindItemInCache(ProductID string) (common.ProductDetails, error)
}

//RedisCacheFinder retrieves item with productId from redis cache
type RedisCacheFinder struct {
	redisClient *redis.Client
}

func Initialize(connUrl string) RedisCacheFinder {
	client := redis.NewClient(&redis.Options{
		Addr:     connUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	cacheFinder := RedisCacheFinder{}
	cacheFinder.redisClient = client

	return cacheFinder
}

func (cacheFinder RedisCacheFinder) close() error {
	err := cacheFinder.redisClient.Close()
	if err != nil {
		return err
	}

	return nil
}

//FindItemInCache retrieves item from redis cache
func (cacheFinder RedisCacheFinder) FindItemInCache(ProductID string) (common.ProductDetails, error) {
	val, err := cacheFinder.redisClient.Get(ProductID).Result()
	if err != nil {
		return common.ProductDetails{}, nil
	}

	var productDetails common.ProductDetails
	decoder := json.NewDecoder(strings.NewReader(val))
	decoder.Decode(&productDetails)
	return productDetails, nil

	// if ProductID == "a123" {
	// 	return common.ProductDetails{ProdID: "a123", ProdName: "iPhone", ProdDesc: "new iPhone", Quantity: 10}, nil
	// } else {
	// 	return common.ProductDetails{}, nil
	// }
}
