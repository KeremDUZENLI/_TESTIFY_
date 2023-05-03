package model

import (
	"testify/common/helper"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

type timeAndPrice struct {
	Timestamp time.Time
	Price     float64
}

func (p *priceProvider) Latest() (*timeAndPrice, error) {
	var data timeAndPrice
	err := p.db.
		QueryRow("SELECT * FROM stockprices ORDER BY timestamp DESC limit 1").
		Scan(&data.Timestamp, &data.Price)

	helper.ErrorLog(err)

	return &data, nil
}

func (p *priceProvider) List(date time.Time) ([]*timeAndPrice, error) {
	var listeData []*timeAndPrice = make([]*timeAndPrice, 0)
	var timestamp time.Time
	var price float64

	rows, _ := p.db.Query("SELECT * FROM stockprices where timestamp::date = $1 ORDER BY timestamp DESC",
		date.Format(dateFormat))

	for rows.Next() {
		if err := rows.Scan(&timestamp, &price); err != nil {
			helper.ErrorLog(err)
		}

		listeData = append(listeData, &timeAndPrice{
			Timestamp: timestamp,
			Price:     price,
		})
	}

	return listeData, nil
}
