package author

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/models"
)

func Create(c *gin.Context) {
	name := "rojo"
	var result models.Author
	coll := config.GetCollections().AuthorModel
	err := coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: bson.M{"$regex": name}}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the name %s\n", name)
		return
	}
	if err != nil {
		panic(err)
	}
	// newLyrics := models.Lyrics{Lyrics: []string{"I'm", "yours"}, Author: result.ID}
	// res, err := lyricsColl.InsertOne(context.TODO(), newLyrics)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res)

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	c.IndentedJSON(200, result)
}
