package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lmich.com/tononkira/domain"
)

// type PaginateResult[V interface{}] struct {
// 	Items        *[]V  `json:"items"`
// 	TotalItems   int64 `json:"totalItems"`
// 	ItemsPerPage int64 `json:"itemsPerPage"`
// 	Page         int64 `json:"page"`
// 	ItemCount    int64 `json:"itemCount"`
// }

// type PaginationOptions struct {
// 	Page      int64
// 	Limit     int64
// 	Filter    bson.D
// 	Sort      bson.D
// 	Collation options.Collation
// }

type DefaultPaginationParams struct {
	Search string `form:"search"`
	Page   int64  `form:"page"`
	Limit  int64  `form:"limit"`
}

type ModelApi[V interface{}] interface {
	ToApi() (V, error)
}

func Paginate[O interface{}, V ModelApi[O]](coll *mongo.Collection, paginateOption domain.PaginationOptions) (domain.PaginateResult[O], error) {
	filter := bson.D{}
	var page int64 = 0
	findOptions := &options.FindOptions{}
	if paginateOption.Page >= 0 {
		page = paginateOption.Page
	}
	if paginateOption.Limit >= 0 {
		findOptions.SetLimit(paginateOption.Limit)
	}
	if paginateOption.Page > 0 && paginateOption.Limit > 0 {
		findOptions.SetSkip(paginateOption.Page * paginateOption.Limit)
	}
	if paginateOption.Sort != nil {
		findOptions.SetSort(paginateOption.Sort)
	}
	// if paginateOption.Collation.(*options.Collation) != nil {
	// 	findOptions.SetCollation(paginateOption.Collation.(*options.Collation).Locale)
	// }
	if paginateOption.Filter != nil {
		filter = paginateOption.Filter.(primitive.D)
	}
	totalItem, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	cursor, err := coll.Find(context.TODO(), filter, findOptions)
	if err != nil {
		panic(err)
	}
	results := []V{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	cursor.Close(context.Background())
	itemsPerPage := *findOptions.Limit
	if *findOptions.Limit > 0 {
		itemsPerPage = int64(len(results))
	}
	resApi := []O{}
	for _, r := range results {
		api, err := r.ToApi()
		if err != nil {
			return domain.PaginateResult[O]{}, err
		}
		resApi = append(resApi, api)
	}
	res := domain.PaginateResult[O]{
		Items:        &resApi,
		TotalItems:   totalItem,
		ItemsPerPage: itemsPerPage,
		Page:         page,
		ItemCount:    int64(len(results)),
	}
	return res, err
}

func ToApi[I interface{}, O interface{}](input domain.PaginateResult[I], outPut *[]O) domain.PaginateResult[O] {
	return domain.PaginateResult[O]{
		Items:        outPut,
		TotalItems:   input.ItemCount,
		ItemsPerPage: input.ItemsPerPage,
		Page:         input.Page,
		ItemCount:    int64(len(*outPut)),
	}
}
