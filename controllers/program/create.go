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
	"lmich.com/tononkira/models"
)

type ProgramCreatePayload struct {
	Date      time.Time            `json:"date"`
	SongsList []primitive.ObjectID `json:"songsList"`
}

func Create(c *gin.Context) {
	var body ProgramCreatePayload
	var program *models.Program
	programModel := config.GetCollections().ProgramModel
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	program = &models.Program{
		Date:      body.Date,
		SongsList: body.SongsList,
	}
	result, err := programModel.InsertOne(context.TODO(), program)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	err = programModel.FindOne(context.TODO(), bson.D{{Key: "_id", Value: result.InsertedID.(primitive.ObjectID)}}).Decode(program)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	programApi := program.ToApi()
	c.IndentedJSON(200, programApi)
}
