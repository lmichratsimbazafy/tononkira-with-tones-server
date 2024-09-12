package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/domain"
)

type Permissions struct {
	SuperAdmin string `json:"superAdmin"`
	Admin      string `json:"admin"`
	Authors    string `json:"authors"`
	Lyrics     string `json:"lyrics"`
	Users      string `json:"users"`
}

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name,omitempty" json:"name"`
	Slug        string             `bson:"slug,omitempty" json:"slug"`
	Permissions Permissions        `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updated_at"`
}

func (r *Role) ToApi() domain.Role {
	return domain.Role{
		ID:          r.ID.Hex(),
		Name:        r.Name,
		Slug:        r.Slug,
		Permissions: domain.Permissions(r.Permissions),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func (r *Role) Upsert() (*Role, error) {
	roleModel := config.GetCollections().RoleModel
	res := roleModel.FindOne(context.TODO(), bson.D{{Key: "slug", Value: bson.M{"$regex": r.Slug, "$options": "i"}}})

	if res.Err() != nil {
		result, err := roleModel.InsertOne(context.TODO(), r)
		if err != nil {
			panic(err)
		}
		fmt.Println(`Created role `, result)
		r.ID = result.InsertedID.(primitive.ObjectID)
	} else if err := res.Decode(&r); err != nil {
		log.Fatal("error while decoding role data")
	}
	fmt.Println("there is already role with name ", r.Name)
	return r, nil
}
