package calculation

import (
	"database/sql"
	"fmt"
	"os"
	"testify/common/helper"
	"testify/database"
	"testify/model"
	"testing"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type SqliteTestSuite struct {
	suite.Suite
	databaseSqlite *sql.DB
	priceIncrease  PriceIncrease
}

func TestSqliteTestSuite(t *testing.T) {
	suite.Run(t, &SqliteTestSuite{})
}

func (sTS *SqliteTestSuite) SetupSuite() {
	databaseSqlite_Set(sTS)

	priceProvider := model.NewPriceProvider(sTS.databaseSqlite)
	sTS.priceIncrease = NewPriceIncrease(priceProvider)
}

func (sTS *SqliteTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "Test_PriceIncrease" {
		database.DbSeedTable(sTS.databaseSqlite)
		// sTS.table_Insert()
	}
}

func (sTS *SqliteTestSuite) TearDownTest() {
	sTS.table_Clean()
}

func (sTS *SqliteTestSuite) TearDownSuite() {
	sTS.table_Drop()
	databaseSqlite_Delete()
}

func (sTS *SqliteTestSuite) Test_PriceIncrease() {
	sTS.table_Retrieve()

	percentage, err := sTS.priceIncrease.PriceIncrease()

	sTS.Nil(err)
	sTS.Equal(25.0, percentage)
}

func (sTS *SqliteTestSuite) Test_PriceIncrease_Error() {
	percentage, err := sTS.priceIncrease.PriceIncrease()

	sTS.EqualError(err, "not enough data")
	sTS.Equal(0.0, percentage)
}

// ----------------------------------------------------------------
func databaseSqlite_Set(sTS *SqliteTestSuite) {
	databaseSqlite_Create(sTS)
	database.DbCreateTable(sTS.databaseSqlite)
	// sTS.table_Create()
}

func databaseSqlite_Create(sTS *SqliteTestSuite) {
	var err error
	sTS.databaseSqlite, err = sql.Open("sqlite3", "./database_sqlite.db")

	helper.ErrorSuite(err)
	helper.ErrorSuite(sTS.databaseSqlite.Ping())
}

func databaseSqlite_Delete() {
	err := os.Remove("./database_sqlite.db")
	helper.ErrorPrint(err)
}

// ----------------------------------------------------------------
func (sTS *SqliteTestSuite) table_Create() {
	_, err := sTS.databaseSqlite.Exec(
		`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp timestamp,
		price float64
	)`)

	helper.ErrorSuite(err)
}

func (sTS *SqliteTestSuite) table_Insert() {
	for i := 1; i <= 5; i++ {
		price := float64((6 - i) * 5)
		timestamp := time.Now().Add(time.Duration(-i) * time.Minute)

		_, err := sTS.databaseSqlite.Exec("INSERT INTO stockprices (timestamp, price) VALUES (?, ?)", timestamp, price)
		helper.ErrorSuite(err)
	}
}

func (sTS *SqliteTestSuite) table_Retrieve() {
	rows, err := sTS.databaseSqlite.Query("SELECT * FROM stockprices")
	helper.ErrorSuite(err)
	defer rows.Close()

	printFromSqlite(rows)
	helper.ErrorSuite(rows.Err())
}

func (sTS *SqliteTestSuite) table_Clean() {
	_, err := sTS.databaseSqlite.Exec("DELETE FROM stockprices")
	helper.ErrorSuite(err)
}

func (sTS *SqliteTestSuite) table_Drop() {
	_, err := sTS.databaseSqlite.Exec("DROP TABLE IF EXISTS stockprices")
	helper.ErrorSuite(err)
}

// ----------------------------------------------------------------
func printFromSqlite(lines *sql.Rows) {
	for lines.Next() {
		var timestamp any
		var price decimal.Decimal

		helper.ErrorSuite(lines.Scan(&timestamp, &price))
		fmt.Printf("timestamp: %v, price: %v\n", timestamp, price)
	}
}
