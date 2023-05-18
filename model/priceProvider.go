package model

import (
	"database/sql"
	"testify/common/helper"
	"time"

	"github.com/lib/pq"
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

func (p *priceProvider) List(date time.Time, args ...*sql.DB) ([]*TimeAndPrice, error) {
	var databaseName *sql.DB
	var query string
	if args == nil {
		databaseName = p.databasePostgre
		query = "SELECT * FROM stockprices where timestamp::date = $1 ORDER BY timestamp DESC"
	} else {
		databaseName = args[0]
		query = "SELECT * FROM stockprices WHERE strftime('%Y-%m-%d', timestamp) = $1 ORDER BY timestamp DESC"
	}

	var listeData []*TimeAndPrice = make([]*TimeAndPrice, 0)
	var timestamp pq.NullTime
	var price float64

	rows, err := databaseName.Query(query,
		date.Format(dateFormat))

	helper.ErrorLog(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&timestamp, &price)
		helper.ErrorLog(err)

		listeData = append(listeData, &TimeAndPrice{
			Timestamp: timestamp.Time,
			Price:     price,
		})
	}

	return listeData, nil
}
