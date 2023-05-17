package calculation

import (
	"database/sql"
	"fmt"
	"os"
	"testify/common/helper"
	"testing"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type SqlTestSuite struct {
	suite.Suite
	databaseSqlite *sql.DB
}

func TestSqlTestSuite(t *testing.T) {
	suite.Run(t, &SqlTestSuite{})
}

func (sts *SqlTestSuite) SetupSuite() {
	setDatabaseSql(sts)
}

func (sts *SqlTestSuite) TearDownSuite() {
	deleteDatabaseSql()
}

func (sts *SqlTestSuite) Test_PriceIncrease() {
}

// ----------------------------------------------------------------
func setDatabaseSql(sts *SqlTestSuite) {
	createDatabaseSql(sts)

	createTable(sts.databaseSqlite)
	insertData(sts.databaseSqlite)
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

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp timestamp,
		price float64
	)`)

	if err != nil {
		return err
	}

	return nil
}

func insertData(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO stockprices(timestamp, price) VALUES(?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	time := time.Now()
	_, err = stmt.Exec(time, 25)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(time, 30)
	if err != nil {
		return err
	}

	return nil
}

func retrieveData(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM stockprices")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp time.Time
		var price float64

		err := rows.Scan(&timestamp, &price)
		if err != nil {
			return err
		}

		fmt.Printf("timestamp: %v, price: %v\n", timestamp, price)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}
