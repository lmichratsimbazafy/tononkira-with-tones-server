package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
)

type Program struct {
	ID        primitive.ObjectID   `bson:"_id,omitenmpty" json:"id"`
	Date      time.Time            `bson:"date" json:"date"`
	SongsList []primitive.ObjectID `bson:"songsList" json:"songsList"`
	CreatedAt time.Time            `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updated_at"`
}

type ApiProgram struct {
	ID        primitive.ObjectID `json:"id"`
	Date      time.Time          `json:"date"`
	SongsList []Lyrics           `json:"songsList"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

func (p *Program) GetSongs() []Lyrics {
	coll := config.GetCollections().ProgramModel
	var songs []Lyrics
	filter := bson.D{{Key: "_id", Value: bson.M{"$in": p.SongsList}}}
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

func (p *Program) ToApi() ApiProgram {
	songs := p.GetSongs()
	return ApiProgram{
		ID:        p.ID,
		Date:      p.Date,
		SongsList: songs,
	}
}

func (p *Program) MarshalBSON() ([]byte, error) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()

	type my Program
	return bson.Marshal((*my)(p))
}
