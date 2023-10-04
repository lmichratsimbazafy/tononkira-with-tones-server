package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
)

type Lyrics struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Lyrics  []string             `bson:"lyrics,required" json:"lyrics"`
	Authors []primitive.ObjectID `bson:"authors,omitempty" json:"authors"`
	Tone    string               `bson:"tone,required" json:"tone"`
	Title   string               `bson:"title,required" json:"title"`
}

type JSONLyrics struct {
	Authors []string `json:"authors"`
	Lyrics  []string `json:"lyrics"`
	Tone    string   `json:"tone"`
	Title   string   `json:"title"`
}

type ApiLyrics struct {
	ID      primitive.ObjectID `json:"id"`
	Lyrics  []string           `json:"lyrics"`
	Authors []Author           `json:"authors"`
	Tone    string             `json:"tone"`
	Title   string             `json:"title"`
}

func (l *Lyrics) GetAuthors() []Author {
	coll := config.GetCollections().AuthorModel
	var authors []Author
	filter := bson.D{{Key: "_id", Value: bson.M{"$in": l.Authors}}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatalln("error while fetching authors", err)
	}
	if err = cursor.All(context.TODO(), &authors); err != nil {
		panic(err)
	}
	cursor.Close(context.Background())
	return authors
}

func (l *Lyrics) ToApi() ApiLyrics {
	authors := l.GetAuthors()
	return ApiLyrics{
		ID:      l.ID,
		Lyrics:  l.Lyrics,
		Tone:    l.Tone,
		Title:   l.Title,
		Authors: authors,
	}
}
