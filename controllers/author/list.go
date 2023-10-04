package author

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/helpers"
	"lmich.com/tononkira/models"
)

type AuthorListParams struct {
	Name string `form:"name"`
	helpers.DefaultPaginationParams
}

func List(c *gin.Context) {
	var params AuthorListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)})
		return
	}
	authorModel := config.GetCollections().AuthorModel
	paginationOptions := helpers.PaginationOptions{}

	filter := bson.D{}
	if params.Name != "" {
		filter = append(filter, bson.E{Key: "name", Value: bson.M{"$regex": params.Name, "$options": "i"}})
	}
	if params.Page >= 0 {
		paginationOptions.Page = params.Page
	}
	if params.Limit >= 0 {
		paginationOptions.Limit = params.Limit
	}
	paginationOptions.Filter = filter
	results := helpers.Paginate[models.Author](authorModel, paginationOptions)

	var authorsApi []models.ApiAuthor
	for _, author := range *results.Items {
		authorsApi = append(authorsApi, author.ToApi())
	}

	c.IndentedJSON(200, helpers.ToApi[models.Author, models.ApiAuthor](results, &authorsApi))
}
