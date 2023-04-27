package main

import (
	"fmt"
	"testify/calculation"
	"testify/database"
	"testify/model"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db := database.DbOperation()
	res := model.NewPriceProvider(db)

	latest, _ := res.Latest()
	latestListe, _ := res.List(time.Now())

	fmt.Println(latest.Price)
	for _, v := range latestListe {
		fmt.Println(v.Price)
	}

	res2 := calculation.NewPriceIncrease(res)
	fmt.Println(res2.PriceIncrease())
}
