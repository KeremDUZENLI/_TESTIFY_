package model

import (
	"database/sql"
	"time"
)

type priceProvider struct {
	databasePostgre *sql.DB
}

type PriceProvider interface {
	Latest() (*TimeAndPrice, error)
	List(date time.Time) ([]*TimeAndPrice, error)
}

func NewPriceProvider(db *sql.DB) PriceProvider {
	return &priceProvider{
		databasePostgre: db,
	}
}
