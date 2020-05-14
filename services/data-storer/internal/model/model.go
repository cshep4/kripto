package model

import "time"

const (
	Buy  TradeType = "buy"
	Sell TradeType = "sell"
)

type (
	GetResponse struct {
		Rates  []Rate          `json:"rates"`
		Trades []TradeResponse `json:"trades"`
	}

	TradeResponse struct {
		Value    Value     `json:"value"`
		Type     TradeType `json:"type"`
		Rate     float64   `json:"rate"`
		DateTime time.Time `json:"dateTime"`
	}

	StoreRequest struct {
		Rate Rate         `json:"rate"`
		Buy  TradeRequest `json:"buy"`
		Sell TradeRequest `json:"sell"`
	}

	TradeRequest struct {
		Traded   bool      `json:"traded"`
		GBP      float64   `json:"gbp"`
		BTC      float64   `json:"btc"`
		Rate     float64   `json:"rate"`
		DateTime time.Time `json:"dateTime"`
	}

	Value struct {
		GBP float64 `json:"gbp"`
		BTC float64 `json:"btc"`
	}

	Rate struct {
		Id       string    `json:"id"`
		Rate     float64   `json:"rate"`
		DateTime time.Time `json:"dateTime"`
	}

	TradeType string

	Trade struct {
		Id       string    `json:"id"`
		GBP      float64   `json:"gbp"`
		BTC      float64   `json:"btc"`
		Rate     float64   `json:"rate"`
		Type     TradeType `json:"type"`
		DateTime time.Time `json:"dateTime"`
	}
)
