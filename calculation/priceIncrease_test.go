package calculation

import (
	"testify/common/env"
	"testify/common/helper"
	"testify/database"
	"testify/model"

	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type IntTestSuite struct {
	suite.Suite
	databasePostgre *sql.DB
	priceIncrease   PriceIncrease
}

func TestIntTestSuite(t *testing.T) {
	suite.Run(t, &IntTestSuite{})
}

func (its *IntTestSuite) SetupSuite() {
	env.Load(1)

	its.databasePostgre = database.DbConnect("postgres_test")
	priceProvider := model.NewPriceProvider(its.databasePostgre)
	its.priceIncrease = NewPriceIncrease(priceProvider)

	database.DbCreateTable()
}

func (its *IntTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestPriceIncrease" {
		database.DbSeedTable()
	}
}

func (its *IntTestSuite) TearDownTest() {
	cleanTable(its)
}

func (its *IntTestSuite) TearDownSuite() {
	dropTable(its)
	closeConnection(its)
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

// ----------------------------------------------------------------
func cleanTable(its *IntTestSuite) {
	_, err := its.databasePostgre.Exec(`DELETE FROM stockprices`)
	helper.ErrorSuite(err)
}

func dropTable(its *IntTestSuite) {
	_, err := its.databasePostgre.Exec(`DROP TABLE stockprices`)
	if err != nil {
		helper.ErrorSuite(err)
	}
}

func closeConnection(its *IntTestSuite) {
	err := its.databasePostgre.Close()
	helper.ErrorSuite(err)
}
