package persistence

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	common "github.com/psuman/go-training/service/common"
	"go.mongodb.org/mongo-driver/mongo"
)

//ItemDao retrieves item with productId from database
type ItemDao interface {
	FindItem(ProductID string) (common.ProductDetails, error)
}

//MongoItemDao retrieves item with productId from mongo database
type MongoItemDao struct {
	mongoClient *mongo.Client
}

// Initialize Initialized connection to mongodb
func Initialize(connUri string) MongoItemDao {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, connUri)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	mongoDao := MongoItemDao{}
	mongoDao.mongoClient = client
	return mongoDao
}

// Close closes mongo db connection
func (dao MongoItemDao) Close() {
	err := dao.mongoClient.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
}

//FindItem retrieves item from Mongo database
func (dao MongoItemDao) FindItem(ProductID string) (common.ProductDetails, error) {
	collection := dao.mongoClient.Database("test").Collection("products")
	filter := bson.D{{"productId", ProductID}}
	var result common.ProductDetails
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return common.ProductDetails{}, err
	}

	return result, nil
}
