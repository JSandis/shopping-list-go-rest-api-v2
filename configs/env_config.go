package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoDBURI() string {
	loadGoDotEnv()
	return os.Getenv("MONGO_URI")
}

func EnvMongoDBName() string {
	loadGoDotEnv()
	return os.Getenv("MONGO_DB_NAME")
}

func EnvMongoDBCollectionNameListItem() string {
	loadGoDotEnv()
	return os.Getenv("MONGO_COLLECTION_NAME_LIST_ITEM")
}

func EnvMongoDBCollectionNameList() string {
	loadGoDotEnv()
	return os.Getenv("MONGO_COLLECTION_NAME_LIST")
}

func EnvMongoDBCollectionNameUser() string {
	loadGoDotEnv()
	return os.Getenv("MONGO_COLLECTION_NAME_USER")
}

func EnvPort() string {
	loadGoDotEnv()
	return os.Getenv("PORT")
}

func EnvSecretKey() string {
	loadGoDotEnv()
	return os.Getenv("SECRET_KEY")
}

func loadGoDotEnv() {
	error := godotenv.Load()
	if error != nil {
		log.Fatal("Error loading .env file")
	}
}
