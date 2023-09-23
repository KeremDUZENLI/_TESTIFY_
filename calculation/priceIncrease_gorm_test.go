package calculation

import (
	"database/sql"
	"fmt"
	"os"
	"testify/common/helper"
	"testify/model"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormTestSuite struct {
	suite.Suite
	databaseGorm  *gorm.DB
	priceIncrease PriceIncrease
}

func TestGormTestSuite(t *testing.T) {
	suite.Run(t, &GormTestSuite{})
}

func (gTS *GormTestSuite) SetupSuite() {
	databaseGorm_Set(gTS)

	priceProvider := model.NewPriceProvider(database_GormToSql(gTS.databaseGorm))
	gTS.priceIncrease = NewPriceIncrease(priceProvider)
}

func (gTS *GormTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "Test_PriceIncrease" {
		// database.DbSeedTable(database_GormToSql(gTS.databaseGorm))
		gTS.table_Insert()
	}
}

func (gTS *GormTestSuite) TearDownTest() {
	gTS.table_Clean()
}

func (gTS *GormTestSuite) TearDownSuite() {
	gTS.table_Drop()
	databaseGorm_Delete()
}

func (gTS *GormTestSuite) Test_PriceIncrease() {
	gTS.table_Retrieve()

	percentage, err := gTS.priceIncrease.PriceIncrease()

	gTS.Nil(err)
	gTS.Equal(25.0, percentage)
}

func (gTS *GormTestSuite) Test_PriceIncrease_Error() {
	percentage, err := gTS.priceIncrease.PriceIncrease()

	gTS.EqualError(err, "not enough data")
	gTS.Equal(0.0, percentage)
}

// ----------------------------------------------------------------
func databaseGorm_Set(gTS *GormTestSuite) {
	databaseGorm_Create(gTS)
	// database.DbCreateTable(database_GormToSql(gTS.databaseGorm))
	gTS.table_Create()
}

func databaseGorm_Create(gTS *GormTestSuite) {
	var err error
	gTS.databaseGorm, err = gorm.Open(sqlite.Open("./database_sqlite.db"), &gorm.Config{})

	helper.ErrorSuite(err)
	helper.ErrorSuite(database_GormToSql(gTS.databaseGorm).Ping())
}

func databaseGorm_Delete() {
	err := os.Remove("./database_sqlite.db")
	helper.ErrorLog(err)
}

// ----------------------------------------------------------------
type stockprices struct {
	Timestamp time.Time
	Price     float64
}

func (gTS *GormTestSuite) table_Create() {
	gTS.databaseGorm.AutoMigrate(&stockprices{})
}

func (gTS *GormTestSuite) table_Insert() {
	for i := 1; i <= 5; i++ {
		price := float64((6 - i) * 5)
		timestamp := time.Now().Add(time.Duration(-i) * time.Minute)

		stockPrice := stockprices{
			Timestamp: timestamp,
			Price:     price,
		}

		err := gTS.databaseGorm.Create(&stockPrice).Error
		helper.ErrorSuite(err)
	}
}

func (gTS *GormTestSuite) table_Retrieve() {
	var retrieveParser []stockprices
	err := gTS.databaseGorm.Find(&retrieveParser).Error

	helper.ErrorSuite(err)
	printFromGorm(retrieveParser)
}

func (gTS *GormTestSuite) table_Clean() {
	gTS.databaseGorm.Exec("DELETE FROM stockprices")
}

func (gTS *GormTestSuite) table_Drop() {
	gTS.databaseGorm.Migrator().DropTable(&stockprices{})
}

// ----------------------------------------------------------------
func printFromGorm(aStruct []stockprices) {
	for _, v := range aStruct {
		fmt.Printf("Timestamp: %v, Price: %v\n",
			v.Timestamp, v.Price)
	}
}

func database_GormToSql(database *gorm.DB) *sql.DB {
	sqlDB, err := database.DB()
	helper.ErrorLog(err)

	return sqlDB
}
