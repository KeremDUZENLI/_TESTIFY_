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

type PostgreTestSuite struct {
	suite.Suite
	databasePostgre *sql.DB
	priceIncrease   PriceIncrease
}

func TestPostgreTestSuite(t *testing.T) {
	suite.Run(t, &PostgreTestSuite{})
}

func (pTS *PostgreTestSuite) SetupSuite() {
	databasePostgre_Set(pTS)

	priceProvider := model.NewPriceProvider(pTS.databasePostgre)
	pTS.priceIncrease = NewPriceIncrease(priceProvider)
}

func (pTS *PostgreTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestPriceIncrease" {
		database.DbSeedTable()
	}
}

func (pTS *PostgreTestSuite) TearDownTest() {
	table_CleanTest(pTS)
}

func (pTS *PostgreTestSuite) TearDownSuite() {
	table_DropTest(pTS)
	databasePostgre_Close(pTS)

	databasePostgre_Connect(pTS)
	databasePostgre_Drop(pTS)
	databasePostgre_Close(pTS)
}

func (pTS *PostgreTestSuite) TestPriceIncrease() {
	percentage, err := pTS.priceIncrease.PriceIncrease()

	pTS.Nil(err)
	pTS.Equal(25.0, percentage)
}

func (pTS *PostgreTestSuite) TestPriceIncrease_Error() {
	percentage, err := pTS.priceIncrease.PriceIncrease()

	pTS.EqualError(err, "not enough data")
	pTS.Equal(0.0, percentage)
}

// ----------------------------------------------------------------
func databasePostgre_Set(pTS *PostgreTestSuite) {
	databasePostgre_Connect(pTS)
	databasePostgre_Create(pTS)

	databasePostgre_Connect(pTS, "postgres_test")
	database.DbCreateTable()
}

func databasePostgre_Connect(pTS *PostgreTestSuite, args ...string) {
	env.Load(1)
	if args == nil {
		pTS.databasePostgre = database.DbConnect()
	} else {
		pTS.databasePostgre = database.DbConnect(args[0])
	}
}

func databasePostgre_Create(pTS *PostgreTestSuite) {
	_, err := pTS.databasePostgre.Exec(`CREATE DATABASE postgres_test`)
	helper.ErrorPrint(err)
}

func databasePostgre_Drop(pTS *PostgreTestSuite) {
	_, err := pTS.databasePostgre.Exec(`DROP DATABASE postgres_test`)
	helper.ErrorSuite(err)
}

func databasePostgre_Close(pTS *PostgreTestSuite) {
	err := pTS.databasePostgre.Close()
	helper.ErrorSuite(err)
}

func table_CleanTest(pTS *PostgreTestSuite) {
	_, err := pTS.databasePostgre.Exec(`DELETE FROM stockprices`)
	helper.ErrorSuite(err)
}

func table_DropTest(pTS *PostgreTestSuite) {
	_, err := pTS.databasePostgre.Exec(`DROP TABLE stockprices`)
	helper.ErrorSuite(err)
}
