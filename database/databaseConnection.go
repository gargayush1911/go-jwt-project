package database

import (
	"fmt"
	"os"
	"time"
	"log"
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBinstance() *mongo.Client{
	err := godotenv.Load(".env")
	if err!=nil{
		log.Fatal("error loading .env file")
	}

	MongoDb:= os.Getenv("MONGODB_URL")
	client,err:= mongo.Connect(options.Client().ApplyURI(MongoDb))
	if err!=nil{
		log.Fatal(err)	
	}
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionname string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionname)
	return collection
}