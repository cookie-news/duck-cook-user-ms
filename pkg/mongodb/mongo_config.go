package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	COLLECTION_CUSTOMER = "customer"
)

func ConnectMongo() mongo.Database {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		panic(err)
	}

	database := client.Database(os.Getenv("MONGODB_DB"))

	if err := client.Database(os.Getenv("MONGODB_DB")).RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}

	fmt.Println("You successfully connected to MongoDB")

	return *database
}
