package data

import (
	"log"

	"lmich.com/tononkira/constants"
	"lmich.com/tononkira/mongodb"
)

func UpsertSuperRoles() {
	roles := []mongodb.Role{
		{
			Name: "SuperAdmin",
			Slug: "superAdmin",
			Permissions: mongodb.Permissions{
				SuperAdmin: constants.Write,
				Authors:    constants.Write,
			},
		},
		{
			Name: "Admin",
			Slug: "admin",
			Permissions: mongodb.Permissions{
				Authors: constants.Write,
				Lyrics:  constants.Write,
				Users:   constants.Write,
			},
		},
	}
	for _, role := range roles {
		if _, err := role.Upsert(); err != nil {
			log.Fatalf("error while upsert %v", err)
		}
	}
}
