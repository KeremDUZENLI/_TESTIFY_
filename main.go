package main

import (
	"database/sql"
	"testify/calculation"
	"testify/common/env"
	"testify/common/helper"
	"testify/database"
	"testify/model"

	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	databasePostgre := setDatabase()
	priceProvider := model.NewPriceProvider(databasePostgre)
	priceIncrease := calculation.NewPriceIncrease(priceProvider)

	increase, err := priceIncrease.PriceIncrease()
	helper.ErrorLog(err)
	fmt.Println(increase)
}

func setDatabase() *sql.DB {
	env.Load()

	db := database.DbConnect()
	database.DbCreateTable()
	database.DbSeedTable()
	database.DbCreateExtra()

	return db
}
