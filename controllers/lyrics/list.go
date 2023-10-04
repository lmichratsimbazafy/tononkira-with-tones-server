package lyrics

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/helpers"
	"lmich.com/tononkira/models"
)

type LyricsListParams struct {
	Author string `form:"author"`
	Title  string `form:"title"`
	helpers.DefaultPaginationParams
}

type LyricsListUriParams struct {
	ID string `uri:"id"`
}

func List(c *gin.Context) {
	var uriParams LyricsListUriParams
	var listParams LyricsListParams
	paginationOptions := helpers.PaginationOptions{}
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}

	if err := c.ShouldBindQuery(&listParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	filter := bson.D{}
	if uriParams.ID != "" {
		authorId, err := primitive.ObjectIDFromHex(uriParams.ID)
		if err != nil {
			panic(err)
		}
		filter = append(filter, bson.E{Key: "authors", Value: bson.M{"$in": []primitive.ObjectID{authorId}}})
	}
	if listParams.Page >= 0 {
		paginationOptions.Page = listParams.Page
	}
	if listParams.Limit >= 0 {
		paginationOptions.Limit = listParams.Limit
	}
	paginationOptions.Filter = filter
	lyricsModel := config.GetCollections().LyricsModel
	results := helpers.Paginate[models.Lyrics](lyricsModel, paginationOptions)
	var apiLyrics []models.ApiLyrics
	for _, song := range *results.Items {
		apiLyrics = append(apiLyrics, song.ToApi())
	}
	c.IndentedJSON(http.StatusOK, helpers.ToApi[models.Lyrics, models.ApiLyrics](results, &apiLyrics))
}
