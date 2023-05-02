package database

import (
	"database/sql"
	"fmt"
	"testify/helper"
	"time"
)

const (
	dbHost     = "192.168.3.204"
	dbPort     = 5402
	dbUser     = "user_testify"
	dbPassword = "password_testify"
	dbName     = "postgres"
)

func DbConnect(args ...string) *sql.DB {
	var databaseName string
	if args == nil {
		databaseName = dbName
	} else {
		databaseName = args[0]
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, databaseName)

	db, err := sql.Open("postgres", psqlInfo)
	helper.ErrorLog(err)

	helper.ErrorLog(db.Ping())

	return db
}

func DbCreateTable(db *sql.DB) {
	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS stockprices (
		timestamp TIMESTAMPTZ PRIMARY KEY,
		price DECIMAL NOT NULL
	)`)
	helper.ErrorLog(err)

	_, err = stmt.Exec()
	helper.ErrorLog(err)
}

func DbSeedTable(db *sql.DB) {
	var rowCount int
	db.QueryRow("SELECT COUNT(*) FROM stockprices").Scan(&rowCount)

	if rowCount != 5 {
		for i := 1; i <= 5; i++ {
			db.Exec("INSERT INTO stockprices (timestamp, price) VALUES ($1,$2)",
				time.Now().Add(time.Duration(-i)*time.Minute), float64((6-i)*5))
		}
	}
}

func DbCreateExtra(db *sql.DB) {
	_, err := db.Exec(`CREATE DATABASE postgres_test`)
	if err != nil {
		println(err.Error())
	}
}
