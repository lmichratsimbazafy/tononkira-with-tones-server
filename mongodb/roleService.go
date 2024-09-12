package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
