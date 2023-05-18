package calculation

import (
	"database/sql"
	"errors"
	"testify/model"
	"time"
)

type PriceIncrease interface {
	PriceIncrease(args ...*sql.DB) (float64, error)
}

type priceIncrease struct {
	PriceProvider model.PriceProvider
}

func NewPriceIncrease(pp model.PriceProvider) PriceIncrease {
	return &priceIncrease{
		PriceProvider: pp,
	}
}

func (pi *priceIncrease) PriceIncrease(args ...*sql.DB) (float64, error) {
	prices, err := pi.PriceProvider.List(time.Now(), args...)
	if err != nil {
		return 0.0, err
	}

	if len(prices) < 2 {
		return 0.0, errors.New("not enough data")
	}

	return (prices[0].Price/prices[1].Price - 1.0) * 100.0, nil
}
