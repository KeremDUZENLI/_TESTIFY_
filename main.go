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
	db := setDatabase()
	pp := model.NewPriceProvider(db)
	calculator := calculation.NewPriceIncrease(pp)

	increase, err := calculator.PriceIncrease()
	helper.ErrorLog(err)
	fmt.Println(increase)
}

func setDatabase() *sql.DB {
	env.Load()

	db := database.DbConnect()

	database.DbCreateTable(db)
	database.DbSeedTable(db)
	database.DbCreateExtra(db)

	return db
}
