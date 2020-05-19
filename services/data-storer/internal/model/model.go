package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	Buy  TradeType = "buy"
	Sell TradeType = "sell"
)

type (
	GetResponse struct {
		Rates  []Rate `json:"rates"`
		Wallet Wallet `json:"wallet"`
	}

	Wallet struct {
		Trades []Trade `json:"trades"`
	}

	StoreRateRequest struct {
		IdempotencyKey string    `json:"idempotencyKey"`
		Rate           float64   `json:"rate"`
		DateTime       time.Time `json:"dateTime"`
	}

	Rate struct {
		Id       string    `json:"id"`
		Rate     float64   `json:"rate"`
		DateTime time.Time `json:"dateTime"`
	}

	Trade struct {
		Id         string    `json:"id"`
		TradeType  TradeType `json:"tradeType"`
		ProductId  string    `json:"productId"`
		Settled    bool      `json:"settled"`
		CreatedAt  time.Time `json:"createdAt,string,omitempty"`
		SpentFunds float64   `json:"funds,omitempty"`
		Fees       float64   `json:"fillFees,omitempty"`
		Value      Value     `json:"value"`
	}

	Value struct {
		GBP float64 `json:"gbp"`
		BTC float64 `json:"btc"`
	}

	TradeType string

	TradeRequest struct {
		Id            string    `json:"id"`
		Side          TradeType `json:"side"`
		ProductId     string    `json:"productId"`
		Funds         string    `json:"funds,omitempty"` // Spent Funds in GBP.
		Settled       bool      `json:"settled"`
		CreatedAt     time.Time `json:"createdAt,string,omitempty"`
		FillFees      string    `json:"fillFees,omitempty"`      // Fees in GBP.
		FilledSize    string    `json:"filledSize,omitempty"`    // Value in BTC.
		ExecutedValue string    `json:"executedValue,omitempty"` // Value in GBP.
	}
)

func (t *TradeRequest) ToTrade() (Trade, error) {
	funds, err := strconv.ParseFloat(strings.TrimSpace(t.Funds), 64)
	if err != nil {
		return Trade{}, fmt.Errorf("invalid_funds: %w", err)
	}
	fees, err := strconv.ParseFloat(strings.TrimSpace(t.FillFees), 64)
	if err != nil {
		return Trade{}, fmt.Errorf("invalid_fees: %w", err)
	}
	btc, err := strconv.ParseFloat(strings.TrimSpace(t.FilledSize), 64)
	if err != nil {
		return Trade{}, fmt.Errorf("invalid_btc_filled_size: %w", err)
	}
	gbp, err := strconv.ParseFloat(strings.TrimSpace(t.ExecutedValue), 64)
	if err != nil {
		return Trade{}, fmt.Errorf("invalid_gbp_executed_value: %w", err)
	}

	return Trade{
		Id:         t.Id,
		TradeType:  t.Side,
		ProductId:  t.ProductId,
		Settled:    t.Settled,
		CreatedAt:  t.CreatedAt,
		SpentFunds: funds,
		Fees:       fees,
		Value: Value{
			GBP: gbp,
			BTC: btc,
		},
	}, nil
}
