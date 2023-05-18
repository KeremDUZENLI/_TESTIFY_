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
	}
}

func (sTS *SqliteTestSuite) TearDownTest() {
	table_Clean(sTS)
}

func (sTS *SqliteTestSuite) TearDownSuite() {
	databaseSqlite_Delete()
}

func (sTS *SqliteTestSuite) Test_PriceIncrease() {
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
}

func databaseSqlite_Create(sTS *SqliteTestSuite) {
	var err error
	sTS.databaseSqlite, err = sql.Open("sqlite3", "./database_sqlite.db")

	helper.ErrorLog(err)
	helper.ErrorLog(sTS.databaseSqlite.Ping())
}

func databaseSqlite_Delete() {
	err := os.Remove("./database_sqlite.db")
	helper.ErrorPrint(err)
}

func table_Clean(sTS *SqliteTestSuite) {
	_, err := sTS.databaseSqlite.Exec(`DELETE FROM stockprices`)
	helper.ErrorPrint(err)
}

func table_Create(sTS *SqliteTestSuite) {
	_, err := sTS.databaseSqlite.Exec(
		`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp timestamp,
		price float64
	)`)

	helper.ErrorLog(err)
}

func table_Insert(sTS *SqliteTestSuite) {
	stmt, err := sTS.databaseSqlite.Prepare("INSERT INTO stockprices(timestamp, price) VALUES(?, ?)")
	helper.ErrorLog(err)
	defer stmt.Close()

	time := time.Now()
	_, err = stmt.Exec(time, 25)
	helper.ErrorLog(err)

	_, err = stmt.Exec(time, 30)
	helper.ErrorLog(err)
}

func table_Retrieve(sTS *SqliteTestSuite) {
	rows, err := sTS.databaseSqlite.Query("SELECT * FROM stockprices")
	helper.ErrorLog(err)
	defer rows.Close()

	for rows.Next() {
		var timestamp any
		var price decimal.Decimal

		helper.ErrorLog(rows.Scan(&timestamp, &price))
		fmt.Printf("timestamp: %v, price: %v\n", timestamp, price)
	}

	helper.ErrorLog(rows.Err())
}
