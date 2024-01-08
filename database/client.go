package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database

func Connect() {
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	Database = c.Database("MonitorService")
	if err != nil {
		panic(err)
	}
	fmt.Println("App is connected to MongoDB")
}

// reusable get collection method taking string as parameter, the collection needs to have a representation in the database
func GetCollection(collection string) *mongo.Collection {
	col := Database.Collection(collection)
	return col
}

func Close() {
	if Database != nil {
		Database.Client().Disconnect(context.TODO())
		Database = nil
		fmt.Println("Database connection closed")
	}
}
