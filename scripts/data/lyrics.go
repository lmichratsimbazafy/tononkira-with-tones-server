package data

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
	"lmich.com/tononkira/mongodb"
)

type JSONLyrics struct {
	Authors []string `json:"authors"`
	Lyrics  []string `json:"lyrics"`
	Tone    string   `json:"tone"`
	Title   string   `json:"title"`
}

func GetLyricsFromFile() []JSONLyrics {
	dir := helpers.GetCurrentDirPath()
	path := filepath.Join(dir, "..", "static", "lyrics")
	var songList []JSONLyrics
	err := filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			fmt.Printf("visited file or dir: %q\n", path)
			data := helpers.ReadJSON[[]JSONLyrics](path)
			songList = append(songList, data...)
		}
		return err
	})
	if err != nil {
		panic("Error while walking through directory")
	}
	return songList
}

func upsertOneLyrics(lyrics *mongodb.Lyrics, lyricsModel *mongo.Collection) *mongodb.Lyrics {

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

func UpsertLyrics(input []JSONLyrics) {
	authorModel := config.GetCollections().AuthorModel
	lyricsModel := config.GetCollections().LyricsModel
	for _, v := range input {
		var authorIds []primitive.ObjectID
		for _, author := range v.Authors {
			authorIds = append(authorIds, upsertAuthor(author, authorModel).ID)
		}
		lyrics := &mongodb.Lyrics{
			Authors: authorIds,
			Lyrics:  v.Lyrics,
			Tone:    v.Tone,
			Title:   v.Title,
		}
		upsertOneLyrics(lyrics, lyricsModel)
	}
}
