package program

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/helpers"
	"lmich.com/tononkira/models"
)

type ProgramDetailsUriParams struct {
	ID   string `uri:"id"`
	Date string `uri:"date"`
}

func Details(c *gin.Context) {
	var uriParams ProgramDetailsUriParams
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	programModel := config.GetCollections().ProgramModel
	var filter bson.D
	var err error

	if uriParams.ID != "" {
		programId, err := primitive.ObjectIDFromHex(uriParams.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
			return
		}
		filter = bson.D{{Key: "_id", Value: programId}}
	} else if uriParams.Date != "" {
		date, err := time.Parse("2006-01-02", uriParams.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
			return
		}
		filter = bson.D{{Key: "date", Value: bson.M{"$gte": helpers.GetStartOfDay(date)}}, {Key: "date", Value: bson.M{"$lt": helpers.GetEndOfDay(date)}}}
	}
	program := &models.Program{}
	err = programModel.FindOne(context.TODO(), filter).Decode(program)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	programApi := program.ToApi()
	c.IndentedJSON(200, programApi)
}
