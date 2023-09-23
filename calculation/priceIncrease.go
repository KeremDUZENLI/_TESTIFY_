package calculation

import (
	"errors"
	"testify/model"
	"time"
)

type priceIncrease struct {
	PriceProvider model.PriceProvider
}

type PriceIncrease interface {
	PriceIncrease() (float64, error)
}

func NewPriceIncrease(pp model.PriceProvider) PriceIncrease {
	return &priceIncrease{
		PriceProvider: pp,
	}
}

func (pi *priceIncrease) PriceIncrease() (float64, error) {
	prices, err := pi.PriceProvider.List(time.Now())
	if err != nil {
		return 0.0, err
	}

	if len(prices) < 2 {
		return 0.0, errors.New("not enough data")
	}

	return (prices[0].Price/prices[1].Price - 1.0) * 100.0, nil
}
