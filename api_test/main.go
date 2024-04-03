package main

import (
	"testing-api-gateway/api_test/handlers"
	"testing-api-gateway/api_test/storage/kv"
	"log"
	"net/http"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // connect to MongoDB
    collection := client.Database("testdb").Collection("testcollection")

	kv.Init(kv.NewMongo(collection))

	router := gin.New()

	router.POST("/user/register", handlers.RegisterUser)
	router.POST("/user/verify/:code", handlers.Verify)
	router.GET("/user/get", handlers.GetUser)
	router.POST("/user/create", handlers.CreateUser)
	router.GET("/users", handlers.ListUsers)

	log.Fatal(http.ListenAndServe(":9999", router))
}
