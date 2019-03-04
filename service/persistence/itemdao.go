package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	common "github.com/psuman/go-training/service/common"
	"go.mongodb.org/mongo-driver/mongo"
)

//ItemDao retrieves item with productId from database
type ItemDao interface {
	FindItem(ProductID string) (common.ProductDetails, error)
	AddItem(productDetails common.ProductDetails) (string, error)
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
	fmt.Println("Inside Mongo Dao")
	collection := dao.mongoClient.Database("test").Collection("products")
	filter := bson.D{{"ProdID", ProductID}}
	var result common.ProductDetails
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println("Could not find item in mongo")
		return common.ProductDetails{}, err
	}

	return result, nil
}

//AddItem adda item to Mongo database
func (dao MongoItemDao) AddItem(productDetails common.ProductDetails) (string, error) {
	fmt.Println("Inside Mongo Dao")
	collection := dao.mongoClient.Database("test").Collection("products")
	doc := bson.D{{"ProdID", productDetails.ProdID}, {"ProdName", productDetails.ProdName}, {"ProdDesc", productDetails.ProdDesc}, {"Quantity", productDetails.Quantity}}

	res, err := collection.InsertOne(context.TODO(), doc)

	if err != nil {
		fmt.Println("Failed to add item to mongo")
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil

}
