package calculation

import (
	"database/sql"
	"testify/database"
	"testify/model"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_PriceIncrease(t *testing.T) {
	asserts := assert.New(t)

	db := database.DbConnect()
	dbSetup(t, db)

	mp := model.NewPriceProvider(db)
	priceIncrease := NewPriceIncrease(mp)
	percentage, err := priceIncrease.PriceIncrease()

	asserts.Nil(err)
	/*
	   if err != nil {
	       t.Logf("err must be nil, but was %s", err.Error())
	       t.Fail()
	   }
	*/

	asserts.Equal(25.0, percentage)
	/*
		if percentage != 25 {
			t.Logf("price increase must be 25, but was %f", percentage)
			t.Fail()
		}
	*/
}

func dbSetup(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`CREATE DATABASE IF NOT EXISTS stockprices_test`)
	if err != nil {
		t.Logf(err.Error())
	}

	database.DbCreateTable(db)
	database.DbSeedTable(db)
}
