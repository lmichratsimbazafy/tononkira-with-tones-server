package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/helpers"
	"lmich.com/tononkira/models"
)

func getAuthorsFromFile() [][]string {
	dir := helpers.GetCurrentDirPath()
	path := filepath.Join(dir, "..", "static", "authors.csv")

	fmt.Println(path)
	return helpers.ReadCsv(path)
}

func getLyricsFromFile() []models.JSONLyrics {
	dir := helpers.GetCurrentDirPath()
	path := filepath.Join(dir, "..", "static", "lyrics")
	var songList []models.JSONLyrics
	err := filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			fmt.Printf("visited file or dir: %q\n", path)
			data := helpers.ReadJSON[[]models.JSONLyrics](path)
			songList = append(songList, data...)
		}
		return err
	})
	if err != nil {
		panic("Error while walking through directory")
	}
	return songList
}
func upsertAuthors(input [][]string) {
	for i, v := range input {
		authorModel := config.GetCollections().AuthorModel
		if i > 0 {
			fmt.Println(i, v[0])
			upsertAuthor(v[0], authorModel)
		}
	}
}

func upsertAuthor(name string, authorModel *mongo.Collection) *models.Author {
	var author *models.Author
	res := authorModel.FindOne(context.TODO(), bson.D{{Key: "name", Value: bson.M{"$regex": name, "$options": "i"}}})

	if res.Err() != nil {
		author = &models.Author{Name: name}
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

func upsertLyrics(input []models.JSONLyrics) {
	songList := getLyricsFromFile()
	authorModel := config.GetCollections().AuthorModel
	lyricsModel := config.GetCollections().LyricsModel
	for _, v := range songList {
		var authorIds []primitive.ObjectID
		for _, author := range v.Authors {
			authorIds = append(authorIds, upsertAuthor(author, authorModel).ID)
		}
		lyrics := &models.Lyrics{
			Authors: authorIds,
			Lyrics:  v.Lyrics,
			Tone:    v.Tone,
			Title:   v.Title,
		}
		upsertOneLyrics(lyrics, lyricsModel)
	}
}

func upsertOneLyrics(lyrics *models.Lyrics, lyricsModel *mongo.Collection) *models.Lyrics {

	res := lyricsModel.FindOne(context.TODO(), bson.D{{Key: "authors", Value: bson.M{"$in": lyrics.Authors}}, {Key: "title", Value: lyrics.Title}})

	if res.Err() != nil {
		result, err := lyricsModel.InsertOne(context.TODO(), lyrics)
		if err != nil {
			panic(err)
		}
		fmt.Println(`Created lyrics `, result)
	} else if err := res.Decode(&lyrics); err != nil {
		log.Fatalf("error while decoding lyrics %s", lyrics.Title)
	}
	fmt.Printf("there is already lyrics with author %s and with title %s ", lyrics.Authors[len(lyrics.Authors)-1], lyrics.Title)
	return lyrics
}

func seed() {
	setupOption := new(config.SetupOptions)
	setupOption.Run = func() {

		authorData := getAuthorsFromFile()
		upsertAuthors(authorData)

		lyricsData := getLyricsFromFile()
		upsertLyrics(lyricsData)
	}
	setupOption.NeedLocalDb = true

	config.Setup(setupOption)

}

func main() {
	seed()
}
