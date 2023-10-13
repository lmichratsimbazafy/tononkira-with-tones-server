package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
)

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name,omitempty" json:"name"`
	Slug        string             `bson:"slug,omitempty" json:"slug"`
	Permissions Permissions        `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updated_at"`
}

type ApiRole struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Permissions Permissions        `json:"permissions"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

func (role *Role) Upsert() error {
	var _role = &Role{}
	roleModel := config.GetCollections().RoleModel
	err := roleModel.FindOne(context.TODO(), bson.D{{Key: "slug", Value: role.Slug}}).Decode(_role)
	if err != nil {
		return err
	}
	if _role != nil {
		if role.Permissions != _role.Permissions {
			fmt.Printf("update permissions of role %s", _role.Slug)
			_role.Permissions = role.Permissions
			if _, err := roleModel.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: _role.ID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "permissions", Value: role.Permissions}}}}); err != nil {
				return err
			}
		}
		return nil
	} else {
		fmt.Printf("create new role of role %s", role.Slug)
		if _, err := roleModel.InsertOne(context.TODO(), role); err != nil {
			return err
		}
		return nil
	}
}

func (r *Role) ToApi() ApiRole {
	return ApiRole{
		ID:          r.ID,
		Name:        r.Name,
		Slug:        r.Slug,
		Permissions: r.Permissions,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
