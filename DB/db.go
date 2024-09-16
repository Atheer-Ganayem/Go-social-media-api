package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB     *mongo.Database
	client *mongo.Client
	err error
)

func Init() {
	clinetOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err = mongo.Connect(context.TODO(), clinetOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB")
	DB = client.Database("go-social-media")
}

func Disconnect() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}