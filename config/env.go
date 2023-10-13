package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	DbUri             string
	DbUser            string
	DbPassword        string
	DbHost            string
	DbPort            string
	DbName            string
	ApiPort           string
	LocalScriptDBHost string
}

func Getenv() *Environment {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("You must set your 'PORT' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	dbUri := os.Getenv("DB_URI")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	if dbUri == "" {
		if dbUser == "" {
			log.Fatal("You must set your 'DB_USER' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}
		if dbPassword == "" {
			log.Fatal("You must set your 'DB_PASSWORD' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}
		if dbHost == "" {
			log.Fatal("You must set your 'DB_HOST' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}
		if dbPort == "" {
			log.Fatal("You must set your 'DB_PORT' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("You must set your 'DB_NAME' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	localScriptDBHost := os.Getenv("LOCAL_SCRIPT_DB_HOST")
	if localScriptDBHost == "" {
		log.Fatal("You must set your 'LOCAL_SCRIPT_DB_HOST' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	return &Environment{
		DbUser:            dbUser,
		DbPassword:        dbPassword,
		DbHost:            dbHost,
		DbPort:            dbPort,
		DbName:            dbName,
		ApiPort:           port,
		LocalScriptDBHost: localScriptDBHost,
		DbUri:             dbUri,
	}
}
