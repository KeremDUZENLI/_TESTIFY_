package model

import (
	"testify/common/helper"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

type TimeAndPrice struct {
	Timestamp time.Time
	Price     float64
}

func (p *priceProvider) Latest() (*TimeAndPrice, error) {
	var data TimeAndPrice
	err := p.databasePostgre.
		QueryRow("SELECT * FROM stockprices ORDER BY timestamp DESC limit 1").
		Scan(&data.Timestamp, &data.Price)

	helper.ErrorLog(err)

	return &data, nil
}

func (p *priceProvider) List(date time.Time) ([]*TimeAndPrice, error) {
	var listeData []*TimeAndPrice = make([]*TimeAndPrice, 0)
	var timestamp time.Time
	var price float64

	rows, err := p.databasePostgre.Query("SELECT * FROM stockprices where timestamp::date = $1 ORDER BY timestamp DESC",
		date.Format(dateFormat))

	helper.ErrorLog(err)

	for rows.Next() {
		err := rows.Scan(&timestamp, &price)
		helper.ErrorLog(err)

		listeData = append(listeData, &TimeAndPrice{
			Timestamp: timestamp,
			Price:     price,
		})
	}

	return listeData, nil
}
