package models

import (
	"context"
	"time"

	// "github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lmich.com/tononkira/config"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserName  string             `bson:"userName,required" json:"userName"`
	Password  string             `bson:"password,required" json:"password"`
	Role      primitive.ObjectID `bson:"role,omitempty" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

type ApiUser struct {
	ID        primitive.ObjectID `json:"id"`
	UserName  string             `json:"userName"`
	Role      Role               `json:"role"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

func (u *User) GetRole() (Role, error) {
	coll := config.GetCollections().RoleModel
	var role Role
	filter := bson.D{{Key: "_id", Value: u.Role}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return Role{}, err
	}
	if err = cursor.All(context.TODO(), &role); err != nil {
		return Role{}, err
	}
	cursor.Close(context.Background())
	return role, nil
}

func (u *User) ToApi() (ApiUser, error) {
	role, err := u.GetRole()
	if err != nil {
		return ApiUser{}, err
	}
	return ApiUser{
		ID:        u.ID,
		UserName:  u.UserName,
		Role:      role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

// func (u *User)GenerateTokenPair() (string, error) {
// 	token := jwt.New(jwt.SigningMethodEdDSA)

// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["sub"] = u.ID
// 	claims["userName"] = u.UserName
// 	claims["exp"] = time.Now().Add(time.Minute *15).Unix()

// }
