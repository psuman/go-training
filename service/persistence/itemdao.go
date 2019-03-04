package persistence

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

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
	Close() error
}

//MongoItemDao retrieves item with productId from mongo database
type MongoItemDao struct {
	mongoClient *mongo.Client
	logger      log.Logger
}

// Initialize Initialized connection to mongodb
func (dao MongoItemDao) Initialize(connUri string, logger log.Logger) MongoItemDao {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, connUri)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	dao.mongoClient = client
	dao.logger = logger
	return dao
}

// Close closes mongo db connection
func (dao MongoItemDao) Close() error {
	err := dao.mongoClient.Disconnect(context.TODO())
	return err
}

//FindItem retrieves item from Mongo database
func (dao MongoItemDao) FindItem(ProductID string) (common.ProductDetails, error) {
	dao.logger.Log("ProdID", ProductID)
	collection := dao.mongoClient.Database("test").Collection("products")
	filter := bson.M{"prodid": ProductID}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var doc common.ProductDetails

	err := collection.FindOne(ctx, filter).Decode(&doc)

	if err != nil {
		dao.logger.Log("dao_error", err.Error())
		return common.ProductDetails{}, err
	}

	dao.logger.Log("Retreived_ProdId", doc.ProdID)

	return doc, nil
}

//AddItem adda item to Mongo database
func (dao MongoItemDao) AddItem(productDetails common.ProductDetails) (string, error) {
	collection := dao.mongoClient.Database("test").Collection("products")
	doc := bson.D{{"prodid", productDetails.ProdID}, {"name", productDetails.Name}, {"desc", productDetails.Desc}, {"quantity", productDetails.Quantity}}

	res, err := collection.InsertOne(context.TODO(), doc)

	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil

}
