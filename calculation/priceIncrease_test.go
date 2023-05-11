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
	setDatabase(its)

	priceProvider := model.NewPriceProvider(its.databasePostgre)
	its.priceIncrease = NewPriceIncrease(priceProvider)
}

func (its *IntTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestPriceIncrease" {
		database.DbSeedTable()
	}
}

func (its *IntTestSuite) TearDownTest() {
	cleanTableTest(its)
}

func (its *IntTestSuite) TearDownSuite() {
	dropTableTest(its)
	closeConnection(its)

	dbConnect(its)
	dropDatabaseTest(its)
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
func setDatabase(its *IntTestSuite) {
	dbConnect(its)
	createDatabaseTest(its)

	dbConnectTest(its)
	database.DbCreateTable()
}

func dbConnect(its *IntTestSuite) {
	env.Load(1)
	its.databasePostgre = database.DbConnect()
}

func dbConnectTest(its *IntTestSuite) {
	env.Load(1)
	its.databasePostgre = database.DbConnect("postgres_test")
}

func createDatabaseTest(its *IntTestSuite) {
	_, err := its.databasePostgre.Exec(`CREATE DATABASE postgres_test`)
	helper.ErrorPrint(err)
}

func cleanTableTest(its *IntTestSuite) {
	_, err := its.databasePostgre.Exec(`DELETE FROM stockprices`)
	helper.ErrorSuite(err)
}

func dropTableTest(its *IntTestSuite) {
	_, err := its.databasePostgre.Exec(`DROP TABLE stockprices`)
	helper.ErrorSuite(err)
}

func dropDatabaseTest(its *IntTestSuite) {
	_, err := its.databasePostgre.Exec(`DROP DATABASE postgres_test`)
	helper.ErrorSuite(err)
}

func closeConnection(its *IntTestSuite) {
	err := its.databasePostgre.Close()
	helper.ErrorSuite(err)
}
