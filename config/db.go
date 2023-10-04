package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBCollections struct {
	AuthorModel *mongo.Collection
	LyricsModel *mongo.Collection
}

var DatabaseInstance *mongo.Database

type Db struct {
	Client      *mongo.Client
	Collections DBCollections
}

func (db *Db) Connect(needLocalDb bool) *mongo.Client {
	env := Getenv()
	var host string

	if needLocalDb == true {
		host = env.LocalScriptDBHost
	} else {
		host = env.DbHost
	}
	uri := "mongodb://" + env.DbUser + ":" + env.DbPassword + "@" + host + ":" + env.DbPort + "/" + env.DbName + "?ssl=false&authSource=admin"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

func (db *Db) Disconnect() {
	fmt.Println("...disconnect from MongoDb")
	if err := db.Client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (db *Db) GetDbInstance() *mongo.Database {
	env := Getenv()
	dataBase := db.Client.Database(env.DbName)
	return dataBase
}

func GetCollections() DBCollections {
	return DBCollections{
		AuthorModel: DatabaseInstance.Collection("authors"),
		LyricsModel: DatabaseInstance.Collection("lyrics"),
	}
}
