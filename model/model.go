package model

import (
	"database/sql"
	"time"
)

type priceProvider struct {
	db *sql.DB
}

type PriceProvider interface {
	Latest() (*timeAndPrice, error)
	List(date time.Time) ([]*timeAndPrice, error)
}

func NewPriceProvider(db *sql.DB) PriceProvider {
	return &priceProvider{
		db: db,
	}
}
