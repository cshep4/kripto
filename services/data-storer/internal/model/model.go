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
		Settled       bool      `json:"settled"`
		CreatedAt     time.Time `json:"createdAt"`
		Funds         string    `json:"funds"`         // Spent Funds in GBP.
		FillFees      string    `json:"fillFees"`      // Fees in GBP.
		FilledSize    string    `json:"filledSize"`    // Value in BTC.
		ExecutedValue string    `json:"executedValue"` // Value in GBP.
	}

	InvalidPropertyError struct {
		Parameter string
		Err       string
	}
)

func (i InvalidPropertyError) Error() string {
	return fmt.Sprintf("invalid parameter %s: %s", i.Parameter, i.Err)
}

func (t *TradeRequest) ToTrade() (Trade, error) {
	switch {
	case t.Id == "":
		return Trade{}, InvalidPropertyError{Parameter: "id", Err: "value is empty"}
	case t.Side == "":
		return Trade{}, InvalidPropertyError{Parameter: "side", Err: "value is empty"}
	case t.ProductId == "":
		return Trade{}, InvalidPropertyError{Parameter: "productId", Err: "value is empty"}
	}

	funds, err := strconv.ParseFloat(strings.TrimSpace(t.Funds), 64)
	if err != nil {
		return Trade{}, InvalidPropertyError{Parameter: "funds", Err: err.Error()}
	}
	fees, err := strconv.ParseFloat(strings.TrimSpace(t.FillFees), 64)
	if err != nil {
		return Trade{}, InvalidPropertyError{Parameter: "fillFees", Err: err.Error()}
	}
	btc, err := strconv.ParseFloat(strings.TrimSpace(t.FilledSize), 64)
	if err != nil {
		return Trade{}, InvalidPropertyError{Parameter: "filledSize", Err: err.Error()}
	}
	gbp, err := strconv.ParseFloat(strings.TrimSpace(t.ExecutedValue), 64)
	if err != nil {
		return Trade{}, InvalidPropertyError{Parameter: "executedValue", Err: err.Error()}
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
