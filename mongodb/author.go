package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/domain"
)

type Author struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string               `bson:"name,omitempty" json:"name"`
	Songs     []primitive.ObjectID `bson:"songs,omitempty" json:"songs"`
	CreatedAt time.Time            `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updated_at"`
}

type ApiAuthor struct {
	ID        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Songs     []Lyrics           `json:"songs"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
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

func (a *Author) ToApi() (domain.Author, error) {
	songs := a.GetSongs()
	songsApi := []domain.Lyrics{}
	for _, song := range songs {
		songsApi = append(songsApi, domain.Lyrics{
			ID:        song.ID.Hex(),
			Lyrics:    song.Lyrics,
			Tone:      song.Tone,
			Title:     song.Title,
			CreatedAt: song.CreatedAt,
			UpdatedAt: song.UpdatedAt,
		})
	}
	return domain.Author{
		ID:        a.ID.Hex(),
		Name:      a.Name,
		Songs:     songsApi,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}, nil
}

func (a *Author) MarshalBSON() ([]byte, error) {
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	a.UpdatedAt = time.Now()

	type my Author
	return bson.Marshal((*my)(a))
}
