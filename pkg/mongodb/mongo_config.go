package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	COLLECTION_CUSTOMER = "customer"
)

func ConnectMongo() mongo.Database {
	uriConnection := fmt.Sprint(
		"mongodb://", os.Getenv("MONGO_USER"), ":", os.Getenv("MONGO_PASSWORD"),
		"@", os.Getenv("MONGO_HOST"), ":", os.Getenv("MONGO_PORT"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uriConnection))
	if err != nil {
		panic(err)
	}

	database := client.Database(os.Getenv("MONGODB_DB"))

	// if err := client.Database(os.Getenv("MONGODB_DB")).RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
	// 	panic(err)
	// }

	fmt.Println("You successfully connected to MongoDB")

	return *database
}
