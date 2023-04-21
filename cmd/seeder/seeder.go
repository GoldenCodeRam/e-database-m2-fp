package main

import (
	"context"
	"log"
	"math/rand"
	"os"

	"github.com/goldencoderam/prg/database"
	"github.com/goldencoderam/prg/model"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URL")).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(context.TODO(), opts)

	coll := client.Database("products").Collection("products")

	// insert products
	products, err := coll.Find(context.TODO(), bson.D{})
	for products.Next(context.TODO()) {
		oid := struct {
			Id primitive.ObjectID `bson:"_id"`
		}{}
		products.Decode(&oid)
		database.ExecuteFakeInsert(model.Product{
			Id: oid.Id.Hex(),
		}, 1)
	}

    client.Disconnect(context.Background())

	// insert people
	database.ExecuteFakeInsert(new(model.Person), 50)

	// insert status types
	err = database.ExecuteInsert("INSERT INTO status_type VALUES ('REQUESTED')")
	err = database.ExecuteInsert("INSERT INTO status_type VALUES ('ON DELIVERY')")
	err = database.ExecuteInsert("INSERT INTO status_type VALUES ('CANCELLED')")
	err = database.ExecuteInsert("INSERT INTO status_type VALUES ('DELIVERED')")
	log.Println(err)

	// insert requests
	err = database.ExecuteFakeInsert(new(model.Request), 500)
	log.Println(err)

	// insert delivery details
	err = database.ExecuteFakeInsert(new(model.DeliveryDetail), 1000)
	log.Println(err)

	// insert request product
	database.ExecuteQuery("select id from request", func(rows pgx.Rows) error {
		for rows.Next() {
			var requestId int
			rows.Scan(&requestId)
			err = database.ExecuteFakeInsert(model.RequestProduct{
				RequestId: requestId,
			}, rand.Intn(10)+1)
			log.Println(err)
		}
		return nil
	})

    database.ExecuteFakeInsert(new(model.Coupon), 200)
}
