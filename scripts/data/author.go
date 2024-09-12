package data

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/helpers"
	"lmich.com/tononkira/mongodb"
)

func GetAuthorsFromFile() [][]string {
	dir := helpers.GetCurrentDirPath()
	path := filepath.Join(dir, "..", "static", "authors.csv")

	fmt.Println(path)
	return helpers.ReadCsv(path)
}

func upsertAuthor(name string, authorModel *mongo.Collection) *mongodb.Author {
	var author *mongodb.Author
	res := authorModel.FindOne(context.TODO(), bson.D{{Key: "name", Value: bson.M{"$regex": name, "$options": "i"}}})

	if res.Err() != nil {
		author = &mongodb.Author{Name: name}
		result, err := authorModel.InsertOne(context.TODO(), author)
		if err != nil {
			panic(err)
		}
		fmt.Println(`Created author `, result)
		author.ID = result.InsertedID.(primitive.ObjectID)
	} else if err := res.Decode(&author); err != nil {
		log.Fatal("error while decoding author data")
	}
	fmt.Println("there is already author with name ", name)
	return author
}

func UpsertAuthors(input [][]string) {
	for i, v := range input {
		authorModel := config.GetCollections().AuthorModel
		if i > 0 {
			fmt.Println(i, v[0])
			upsertAuthor(v[0], authorModel)
		}
	}
}
