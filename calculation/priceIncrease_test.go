package calculation

import (
	"errors"
	"testify/common/env"
	"testify/common/helper"
	"testify/database"
	"testify/mocks"
	"testify/model"
	"time"

	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/mock"
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

type UnitTestSuite struct {
	suite.Suite
	priceIncrease     PriceIncrease
	priceProviderMock *mocks.PriceProvider
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) SetupTest() {
	priceProviderMock := mocks.PriceProvider{}

	uts.priceIncrease = NewPriceIncrease(&priceProviderMock)
	uts.priceProviderMock = &priceProviderMock
}

func (uts *UnitTestSuite) TestCalculate() {
	uts.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{
		{
			Timestamp: time.Now(),
			Price:     2.0,
		},
		{
			Timestamp: time.Now().Add(time.Duration(-1 * time.Minute)),
			Price:     1.0,
		},
	}, nil)

	actual, err := uts.priceIncrease.PriceIncrease()

	uts.Equal(100.0, actual)
	uts.Nil(err)
}

func (uts *UnitTestSuite) TestCalculate_Error() {
	uts.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{}, nil)

	actual, err := uts.priceIncrease.PriceIncrease()

	uts.Equal(0.0, actual)
	uts.EqualError(err, "not enough data")
}

func (uts *UnitTestSuite) TestCalculate_ErrorFromPriceProvider() {
	expectedError := errors.New("oh my deuss")

	uts.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{}, expectedError)

	actual, err := uts.priceIncrease.PriceIncrease()

	uts.Equal(0.0, actual)
	uts.Equal(expectedError, err)
}
