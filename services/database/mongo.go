package database

import (
	"api/migrations"
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Client

var MongoDatabase *mongo.Database

// ConnectMongo initialize MongoDB connection
func ConnectMongo() {
	if nil == Mongo {
		client := mongoConnect()
		if nil != client {
			Mongo = client
			// Run migrations if enabled
			if viper.GetBool("ENABLE_MIGRATION") {
				mongoMigrate()
			}
		}
	}
}

func mongoConnect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB configuration from viper
	host := viper.GetString("MONGO_HOST")
	if host == "" {
		host = "localhost"
	}
	port := viper.GetString("MONGO_PORT")
	if port == "" {
		port = "27017"
	}
	username := viper.GetString("MONGO_USERNAME")
	password := viper.GetString("MONGO_PASSWORD")
	database := viper.GetString("MONGO_DATABASE")

	// Build MongoDB connection URI
	var uri string
	if username != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
			username, password, host, port, database)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s/%s", host, port, database)
	}

	// Create MongoDB client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to ping MongoDB: %v", err))
	}

	Mongo = client
	MongoDatabase = client.Database(database)

	return client
}

// mongoMigrate runs MongoDB migrations
func mongoMigrate() {
	if nil == Mongo {
		return
	}

	if err := migrations.RunMongoMigrations(Mongo); err != nil {
		panic(fmt.Sprintf("Failed to run MongoDB migrations: %v", err))
	}
}
