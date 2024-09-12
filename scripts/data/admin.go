package data

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/constants"
	"lmich.com/tononkira/helpers"
	"lmich.com/tononkira/mongodb"
)

func UpsertSuperAdminAdmin() {
	userModel := config.GetCollections().UserModel
	existing := &mongodb.User{}
	adminUserName := config.Env.AdminUserName
	res := userModel.FindOne(context.TODO(), bson.D{{Key: "userName", Value: adminUserName}})
	if res.Err() != nil {
		fmt.Println("error while getting admin user ", res.Err())
		role := &mongodb.Role{
			Slug: "superAdmin",
			Permissions: mongodb.Permissions{
				SuperAdmin: constants.Write,
			},
		}
		role, err := role.Upsert()
		if err != nil {
			log.Fatalf("error while role upsert %v", err)
		}
		password, err := helpers.EncryptPassword("Password123!")
		if err != nil {
			log.Fatalf("error while encrypting password %v", err)
		}
		admin := &mongodb.User{
			UserName: adminUserName,
			Password: password,
			Role:     role.ID,
		}
		result, err := userModel.InsertOne(context.TODO(), admin)
		if err != nil {
			log.Fatalf("error while creating admin user %v", err)
		}
		fmt.Println(`Created adminUser `, result.InsertedID)
	} else if err := res.Decode(existing); err != nil {
		log.Fatalf("error while decoding existing admin user %v", err)
	}
}
