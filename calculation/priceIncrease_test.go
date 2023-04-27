package calculation

import (
	"database/sql"
	"testify/database"
	"testify/model"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func Test_PriceIncrease(t *testing.T) {
	requires := require.New(t)

	db := database.DbConnect()
	dbSetup(t, db)

	mp := model.NewPriceProvider(db)
	priceIncrease := NewPriceIncrease(mp)
	percentage, err := priceIncrease.PriceIncrease()

	requires.Nil(err)
	requires.Equal(25.0, percentage)
}

func dbSetup(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`CREATE DATABASE IF NOT EXISTS stockprices_test`)
	if err != nil {
		t.Logf(err.Error())
	}

	database.DbCreateTable(db)
	database.DbSeedTable(db)
}
