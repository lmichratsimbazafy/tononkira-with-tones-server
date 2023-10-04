package helpers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginateResult[V interface{}] struct {
	Items        *[]V  `json:"items"`
	TotalItems   int64 `json:"totalItems"`
	ItemsPerPage int64 `json:"itemsPerPage"`
	Page         int64 `json:"page"`
	ItemCount    int64 `json:"itemCount"`
}

type PaginationOptions struct {
	Page      int64
	Limit     int64
	Filter    bson.D
	Sort      bson.D
	Collation options.Collation
}

type DefaultPaginationParams struct {
	Search string `form:"search"`
	Page   int64  `form:"page"`
	Limit  int64  `form:"limit"`
}

func Paginate[V interface{}](coll *mongo.Collection, paginateOption PaginationOptions) PaginateResult[V] {
	filter := bson.D{}
	var page int64 = 0
	options := &options.FindOptions{}
	if paginateOption.Page >= 0 {
		page = paginateOption.Page
	}
	if paginateOption.Limit >= 0 {
		options.SetLimit(paginateOption.Limit)
	}
	if paginateOption.Page > 0 && paginateOption.Limit > 0 {
		options.SetSkip(paginateOption.Page * paginateOption.Limit)
	}
	if paginateOption.Sort != nil {
		options.SetSort(paginateOption.Sort)
	}
	if paginateOption.Collation.Locale != "" {
		options.SetCollation(&paginateOption.Collation)
	}
	if paginateOption.Filter != nil {
		filter = paginateOption.Filter
	}
	totalItem, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	cursor, err := coll.Find(context.TODO(), filter, options)
	if err != nil {
		panic(err)
	}
	results := []V{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	cursor.Close(context.Background())
	itemsPerPage := *options.Limit
	if *options.Limit > 0 {
		itemsPerPage = int64(len(results))
	}

	res := PaginateResult[V]{
		Items:        &results,
		TotalItems:   totalItem,
		ItemsPerPage: itemsPerPage,
		Page:         page,
		ItemCount:    int64(len(results)),
	}
	return res
}

func ToApi[I interface{}, O interface{}](input PaginateResult[I], outPut *[]O) PaginateResult[O] {
	return PaginateResult[O]{
		Items:        outPut,
		TotalItems:   input.ItemCount,
		ItemsPerPage: input.ItemsPerPage,
		Page:         input.Page,
		ItemCount:    int64(len(*outPut)),
	}
}
