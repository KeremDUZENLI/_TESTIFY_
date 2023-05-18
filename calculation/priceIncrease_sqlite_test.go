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

type SqlTestSuite struct {
	suite.Suite
	databaseSqlite *sql.DB
	priceIncrease  PriceIncrease
}

func TestSqlTestSuite(t *testing.T) {
	suite.Run(t, &SqlTestSuite{})
}

func (sts *SqlTestSuite) SetupSuite() {
	setDatabaseSql(sts)

	priceProvider := model.NewPriceProvider(sts.databaseSqlite)
	sts.priceIncrease = NewPriceIncrease(priceProvider)
}

func (sts *SqlTestSuite) TearDownSuite() {
	deleteDatabaseSql()
}

func (sts *SqlTestSuite) Test_PriceIncrease() {
	percentage, err := sts.priceIncrease.PriceIncrease(sts.databaseSqlite)

	fmt.Println("percentage, err: ", percentage, err)

	sts.Nil(err)
	sts.Equal(25.0, percentage)
}

// ----------------------------------------------------------------
func setDatabaseSql(sts *SqlTestSuite) {
	createDatabaseSql(sts)

	// createTable(sts.databaseSqlite)
	database.DbCreateTable(sts.databaseSqlite)

	// insertData(sts.databaseSqlite)
	database.DbSeedTable(sts.databaseSqlite)

	retrieveData(sts.databaseSqlite)
}

func createDatabaseSql(sts *SqlTestSuite) {
	var err error
	sts.databaseSqlite, err = sql.Open("sqlite3", "./database_sqlite.db")

	helper.ErrorLog(err)
	helper.ErrorLog(sts.databaseSqlite.Ping())
}

func deleteDatabaseSql() {
	err := os.Remove("./database_sqlite.db")
	helper.ErrorPrint(err)
}

func createTable(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp timestamp,
		price float64
	)`)

	helper.ErrorLog(err)
}

func insertData(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO stockprices(timestamp, price) VALUES(?, ?)")
	helper.ErrorLog(err)
	defer stmt.Close()

	time := time.Now()
	_, err = stmt.Exec(time, 25)
	helper.ErrorLog(err)

	_, err = stmt.Exec(time, 30)
	helper.ErrorLog(err)
}

func retrieveData(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM stockprices")
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
