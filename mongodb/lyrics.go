package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/domain"
)

type MongoLyricsService struct {
	Collection *mongo.Collection
}

type Lyrics struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Lyrics    []string             `bson:"lyrics,required" json:"lyrics"`
	Authors   []primitive.ObjectID `bson:"authors,omitempty" json:"authors"`
	Tone      string               `bson:"tone,required" json:"tone"`
	Title     string               `bson:"title,required" json:"title"`
	CreatedAt time.Time            `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updated_at"`
}

func (l *Lyrics) GetAuthors() ([]Author, error) {
	coll := config.GetCollections().AuthorModel
	var authors []Author
	filter := bson.D{{Key: "_id", Value: bson.M{"$in": l.Authors}}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Fatalln("error while fetching authors", err)
		return nil, err
	}
	if err = cursor.All(context.TODO(), &authors); err != nil {
		panic(err)
	}
	cursor.Close(context.Background())
	return authors, nil
}

func (l *Lyrics) ToApi() (domain.Lyrics, error) {
	authors, err := l.GetAuthors()
	if err != nil {
		return domain.Lyrics{}, err
	}
	authorsApi := []domain.Author{}
	for _, author := range authors {

		authorsApi = append(authorsApi, domain.Author{
			ID:        author.ID.Hex(),
			Name:      author.Name,
			CreatedAt: author.CreatedAt,
			UpdatedAt: author.UpdatedAt,
		})
	}
	return domain.Lyrics{
		ID:        l.ID.Hex(),
		Lyrics:    l.Lyrics,
		Tone:      l.Tone,
		Title:     l.Title,
		Authors:   authorsApi,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}, nil
}
func (m *MongoLyricsService) List(params domain.PaginationOptions) (domain.PaginateResult[domain.Lyrics], error) {
	return Paginate[domain.Lyrics, *Lyrics](m.Collection, params)
}

func (m *MongoLyricsService) Create(l *domain.Lyrics) error {
	_, err := m.Collection.InsertOne(context.TODO(), l)
	return err
}
