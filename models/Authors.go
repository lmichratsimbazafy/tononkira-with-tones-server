package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
)

type Author struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name  string               `bson:"name,omitempty" json:"name"`
	Songs []primitive.ObjectID `bson:"songs,omitempty" json:"songs"`
}

type ApiAuthor struct {
	ID    primitive.ObjectID `json:"id"`
	Name  string             `json:"name"`
	Songs []Lyrics           `json:"songs"`
}

func (a *Author) GetSongs() []Lyrics {
	coll := config.GetCollections().LyricsModel
	var songs []Lyrics
	filter := bson.D{{Key: "authors", Value: bson.M{"$in": []primitive.ObjectID{a.ID}}}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatalln("error while fetching songs", err)
	}
	if err = cursor.All(context.TODO(), &songs); err != nil {
		panic(err)
	}
	cursor.Close(context.Background())
	return songs
}

func (a *Author) ToApi() ApiAuthor {
	songs := a.GetSongs()
	return ApiAuthor{
		ID:    a.ID,
		Name:  a.Name,
		Songs: songs,
	}
}
