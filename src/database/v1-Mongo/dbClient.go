package dbClient

import (
	"context"
	"fmt"
	"log"
	"time"

	envHelper "back-end/mapSearchService/env"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx, Cancel = context.WithTimeout(context.Background(), 36000*time.Second)

func Connect() *mongo.Client {    
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)      
	atlasUri := envHelper.GetEnvVar("ATLAS_URI")
	connectionOptions := options.Client().ApplyURI(atlasUri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(Ctx, connectionOptions)

	if err != nil {
		Disconnect(client)
		panic(err)
	}
	
	if err := client.Database("DevCluster").RunCommand(Ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatal("There was a connection error", err)
		Disconnect(client)
		panic(err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	return client
}

func Disconnect(client *mongo.Client) {
	client.Disconnect(Ctx)
}