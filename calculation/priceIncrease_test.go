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
	its.db = database.DbConnect("stockprices_test")
	database.DbCreateTable(its.db)

	priceProvider := model.NewPriceProvider(its.db)
	its.priceIncrease = NewPriceIncrease(priceProvider)
}

func (its *IntTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestPriceIncrease_Error" {
		return
	}
	database.DbSeedTable(its.db)
}

func (its *IntTestSuite) TearDownTest() {
	cleanTable(its)
}

func (its *IntTestSuite) TearDownSuite() {
	tearDownDatabase(its)
}

func (its *IntTestSuite) TestPriceIncrease() {
	percentage, err := its.priceIncrease.PriceIncrease()

	its.Nil(err)
	its.Equal(25.0, percentage)
}

func (its *IntTestSuite) TestPriceIncrease_Error() {
	percentage, err := its.priceIncrease.PriceIncrease()

	its.EqualError(err, "not enough data")
	its.Equal(0.0, percentage)
}

func cleanTable(its *IntTestSuite) {
	_, err := its.db.Exec(`DELETE FROM stockprices`)
	helper.ErrorIts(err)
}

func tearDownDatabase(its *IntTestSuite) {
	_, err := its.db.Exec(`DROP TABLE stockprices`)
	if err != nil {
		helper.ErrorIts(err)
	}

	err = its.db.Close()
	if err != nil {
		helper.ErrorIts(err)
	}
}
