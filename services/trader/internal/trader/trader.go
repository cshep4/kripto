package trader

import (
	"fmt"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
)

const (
	Buy  TradeType = "buy"
	Sell TradeType = "sell"
)

type (
	TradeType string

	TradeResponse struct {
		Id            string    `json:"id"`
		Side          string    `json:"side"`
		ProductId     string    `json:"productId"`
		Settled       bool      `json:"settled"`
		CreatedAt     time.Time `json:"createdAt,string,omitempty"`
		Funds         string    `json:"funds,omitempty"`         // Spent Funds in GBP.
		FillFees      string    `json:"fillFees,omitempty"`      // Fees in GBP.
		FilledSize    string    `json:"filledSize,omitempty"`    // Value in BTC.
		ExecutedValue string    `json:"executedValue,omitempty"` // Value in GBP.
	}

	Coinbase interface {
		CreateOrder(order *coinbasepro.Order) (coinbasepro.Order, error)
		GetOrder(id string) (coinbasepro.Order, error)
	}

	trader struct {
		coinbase Coinbase
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func New(coinbase Coinbase) (*trader, error) {
	if coinbase == nil {
		return nil, InvalidParameterError{Parameter: "coinbase"}
	}

	return &trader{
		coinbase: coinbase,
	}, nil
}

func (t *trader) Trade(tradeType TradeType, amount string) (*TradeResponse, error) {
	o := &coinbasepro.Order{
		Funds:     amount,
		Side:      string(tradeType),
		ProductID: "BTC-GBP",
		Type:      "market",
	}
	order, err := t.coinbase.CreateOrder(o)
	if err != nil {
		return nil, fmt.Errorf("create_order: %w", err)
	}

	// Must get the order after creating as it will not be settled in previous response
	order, err = t.coinbase.GetOrder(order.ID)
	if err != nil {
		return nil, fmt.Errorf("get_order: %w", err)
	}

	return &TradeResponse{
		Side:          order.Side,
		ProductId:     order.ProductID,
		Funds:         order.Funds,
		Id:            order.ID,
		Settled:       order.Settled,
		CreatedAt:     time.Time(order.CreatedAt),
		FillFees:      order.FillFees,
		FilledSize:    order.FilledSize,
		ExecutedValue: order.ExecutedValue,
	}, nil
}
