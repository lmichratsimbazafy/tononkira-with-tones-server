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

type ProgramUpdatePayload struct {
	Date      time.Time            `json:"date"`
	SongsList []primitive.ObjectID `json:"songsList"`
}
type ProgramUpdateUriParams struct {
	ID string `uri:"id"`
}

func Update(c *gin.Context) {
	var body *ProgramCreatePayload
	program := &models.Program{}
	var uriParams ProgramUpdateUriParams
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	programModel := config.GetCollections().ProgramModel
	programId, err := primitive.ObjectIDFromHex(uriParams.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	err = programModel.FindOne(context.TODO(), bson.D{{Key: "_id", Value: programId}}).Decode(program)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	updateDocument := bson.D{}
	if body.SongsList != nil {
		updateDocument = append(updateDocument, bson.E{Key: "songsList", Value: body.SongsList})
	}
	if !body.Date.IsZero() {
		updateDocument = append(updateDocument, bson.E{Key: "date", Value: body.Date})
	}
	_, err = programModel.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: programId}}, bson.D{{Key: "$set", Value: updateDocument}})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	err = programModel.FindOne(context.TODO(), bson.D{{Key: "_id", Value: programId}}).Decode(program)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	programApi := program.ToApi()
	c.IndentedJSON(200, programApi)
}
