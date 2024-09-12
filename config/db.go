package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBCollections struct {
	AuthorModel    *mongo.Collection
	LyricsModel    *mongo.Collection
	ProgramModel   *mongo.Collection
	UserModel      *mongo.Collection
	RoleModel      *mongo.Collection
	AuthTokenModel *mongo.Collection
}

var DatabaseInstance *mongo.Database

type Db struct {
	Client      *mongo.Client
	Collections DBCollections
}

const (
	Authors    = "authors"
	Lyrics     = "lyrics"
	Programs   = "programs"
	Users      = "users"
	Roles      = "roles"
	AuthTokens = "authTokens"
)

func (db *Db) Connect(needLocalDb bool) *mongo.Client {
	env := Env
	var host string

	if needLocalDb {
		host = env.DbHost
	} else {
		host = env.LocalScriptDBHost
	}
	var uri string
	if env.DbUri != "" {
		uri = env.DbUri
	} else {
		uri = "mongodb://" + env.DbUser + ":" + env.DbPassword + "@" + host + ":" + env.DbPort + "/" + env.DbName + "?ssl=false&authSource=admin"
	}
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
	env := Env
	dataBase := db.Client.Database(env.DbName)
	return dataBase
}

func GetCollections() DBCollections {
	return DBCollections{
		AuthorModel:    DatabaseInstance.Collection(Authors),
		LyricsModel:    DatabaseInstance.Collection(Lyrics),
		ProgramModel:   DatabaseInstance.Collection(Programs),
		UserModel:      DatabaseInstance.Collection(Users),
		RoleModel:      DatabaseInstance.Collection(Roles),
		AuthTokenModel: DatabaseInstance.Collection(AuthTokens),
	}
}
