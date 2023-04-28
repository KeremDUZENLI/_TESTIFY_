package calculation

import (
	"database/sql"
	"testify/database"
	"testify/helper"
	"testify/model"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type IntTestSuite struct {
	suite.Suite
	db            *sql.DB
	priceIncrease PriceIncrease
}

func TestIntTestSuite(t *testing.T) {
	suite.Run(t, &IntTestSuite{})
}

func (its *IntTestSuite) SetupSuite() {
	db := database.DbConnect()

	dbCreateTableTest(db)
	database.DbCreateTable(db)

	mp := model.NewPriceProvider(db)
	priceIncrease := NewPriceIncrease(mp)

	its.db = db
	its.priceIncrease = priceIncrease
}

func (its *IntTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestPriceIncrease_Error" {
		return
	}

	database.DbSeedTable(its.db)
}

func (its *IntTestSuite) TearDownSuite() {
	tearDownDatabase(its)
}

func (its *IntTestSuite) TearDownTest() {
	cleanTable(its)
}

func (its *IntTestSuite) TestPriceIncrease_Error() {
	percentage, err := its.priceIncrease.PriceIncrease()

	its.EqualError(err, "not enough data")
	its.Equal(0.0, percentage)
}

func (its *IntTestSuite) TestPriceIncrease() {
	percentage, err := its.priceIncrease.PriceIncrease()

	its.Nil(err)
	its.Equal(25.0, percentage)
}

func dbCreateTableTest(db *sql.DB) {
	println("setting up database")

	_, err := db.Exec(`CREATE DATABASE stockprices_test`)
	helper.ErrorIts(err)
}

func cleanTable(its *IntTestSuite) {
	println("cleaning up table")

	_, err := its.db.Exec(`DELETE FROM stockprices`)
	helper.ErrorIts(err)
}

func tearDownDatabase(its *IntTestSuite) {
	println("drop table & database")

	_, err := its.db.Exec(`DROP TABLE stockprices`)
	helper.ErrorIts(err)

	_, err = its.db.Exec(`DROP DATABASE stockprices_test`)
	helper.ErrorIts(err)

	err = its.db.Close()
	helper.ErrorIts(err)
}
