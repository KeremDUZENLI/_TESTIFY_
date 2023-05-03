package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DbHost string
	DbPort string
	DbUser string
	DbPass string
	DbName string
)

func Load(args ...int) {
	var subFolder int
	if args == nil {
		subFolder = 0
	} else {
		subFolder = args[0]
	}

	envFounder(subFolder)

	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPass = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")
}

func envFounder(subFolder int) {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}

	for subFolder > 0 {
		path += "/.."
		subFolder--
	}

	path += "/.env"
	if err := godotenv.Load(path); err != nil {
		log.Fatal(err.Error())
	}
}
