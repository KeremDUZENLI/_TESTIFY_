package model

import (
	"database/sql"
	"testify/helper"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

type Data struct {
	Timestamp time.Time
	Price     float64
}

// ----------------------------------------------------------------

type PriceProvider interface {
	Latest() (*Data, error)
	List(date time.Time) ([]*Data, error)
}

type priceProvider struct {
	db *sql.DB
}

func NewPriceProvider(db *sql.DB) PriceProvider {
	return &priceProvider{
		db: db,
	}
}

func (p *priceProvider) Latest() (*Data, error) {
	var data Data
	err := p.db.
		QueryRow("SELECT * FROM stockprices ORDER BY timestamp DESC limit 1").
		Scan(&data.Timestamp, &data.Price)

	helper.ErrorLog(err)

	return &data, nil
}

func (p *priceProvider) List(date time.Time) ([]*Data, error) {
	var listeData []*Data = make([]*Data, 0)
	var timestamp time.Time
	var price float64

	rows, _ := p.db.Query("SELECT * FROM stockprices where timestamp::date = $1 ORDER BY timestamp DESC",
		date.Format(dateFormat))

	for rows.Next() {
		if err := rows.Scan(&timestamp, &price); err != nil {
			helper.ErrorLog(err)
		}

		listeData = append(listeData, &Data{
			Timestamp: timestamp,
			Price:     price,
		})
	}

	return listeData, nil
}
