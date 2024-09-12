package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/domain"
	"lmich.com/tononkira/helpers"
)

type LyricsListParams struct {
	Author   string    `form:"author"`
	Title    string    `form:"title"`
	FromDate time.Time `form:"fromDate"`
	helpers.DefaultPaginationParams
}

type LyricsListUriParams struct {
	ID string `uri:"id"`
}
type LyricshHandler struct {
	LyricsService domain.LyricsService
}

func (l *LyricshHandler) List(c *gin.Context) {
	var uriParams LyricsListUriParams
	var listParams LyricsListParams
	paginationOptions := domain.PaginationOptions{}
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
	if !listParams.FromDate.IsZero() {
		filter = append(filter, bson.E{Key: "created_at", Value: bson.M{"$gt": listParams.FromDate}})
	}
	if listParams.Page >= 0 {
		paginationOptions.Page = listParams.Page
	}
	if listParams.Limit >= 0 {
		paginationOptions.Limit = listParams.Limit
	}
	paginationOptions.Filter = filter
	results, err := l.LyricsService.List(paginationOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}
