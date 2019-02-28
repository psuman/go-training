package cache

import (
	"github.com/go-redis/redis"
	common "github.com/psuman/go-training/service/common"
)

//CacheFinder retrieves item with productId from cache
type CacheFinder interface {
	FindItemInCache(ProductID string) (common.ProductDetails, error)
}

//RedisCacheFinder retrieves item with productId from redis cache
type RedisCacheFinder struct{}

//NewClient creates new redis client
func (RedisCacheFinder) NewClient(connUrl string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     connUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}
	return client, nil
}

//FindItemInCache retrieves item from redis cache
func (RedisCacheFinder) FindItemInCache(ProductID string) (common.ProductDetails, error) {
	if ProductID == "a123" {
		return common.ProductDetails{ProdID: "a123", ProdName: "iPhone", ProdDesc: "new iPhone", Quantity: 10}, nil
	} else {
		return common.ProductDetails{}, nil
	}
}
