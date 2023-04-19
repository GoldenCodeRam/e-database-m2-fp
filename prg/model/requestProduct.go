package model

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/goldencoderam/prg/database"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RequestProduct struct {
	ProductValue float32
	RequestId    int
	ProductId    string
}

func (r RequestProduct) NewFakeRecord() (database.Insertable, error) {
	var requestId int
	var productId string

	database.ExecuteRowQuery("SELECT id FROM request ORDER BY random() LIMIT 1", &requestId)
	database.ExecuteRowQuery("SELECT id FROM product ORDER BY random() LIMIT 1", &productId)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URL")).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(context.TODO(), opts)

	coll := client.Database("products").Collection("products")

	// insert products
    productObjectId, err := primitive.ObjectIDFromHex(productId)
    log.Println(err)
	product := coll.FindOne(context.TODO(), bson.D{
		{Key: "_id", Value: productObjectId},
	})

	productStr := struct {
		Id primitive.ObjectID `bson:"_id"`
        Price float32 `bson:"price"`
	}{}
    product.Decode(&productStr)

    defer client.Disconnect(context.Background())

	return RequestProduct{
		ProductValue: productStr.Price,
		RequestId:    r.RequestId,
		ProductId:    productId,
	}, nil
}

func (r RequestProduct) GenerateInsertSql() string {
	return fmt.Sprintf(
		`INSERT INTO request_product (product_value, request_id, product_id)
        VALUES (%f, %d, '%s')`,
		r.ProductValue,
		r.RequestId,
		r.ProductId,
	)
}
