package env

import (
	"os"
	"testify/common/helper"

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
	helper.ErrorLog(err)

	for i := subFolder; i > 0; i-- {
		path += "/.."
	}

	err = godotenv.Load(path + "/.env")
	helper.ErrorLog(err)
}
