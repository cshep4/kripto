package mongo

import (
	"fmt"
	"time"

	"github.com/cshep4/kripto/services/data-storer/internal/model"
)

type (
	trade struct {
		Id         string    `bson:"_id"`
		TradeType  string    `bson:"tradeType"`
		ProductId  string    `bson:"productId"`
		Settled    bool      `bson:"settled"`
		CreatedAt  time.Time `bson:"createdAt,string,omitempty"`
		SpentFunds float64   `bson:"funds,omitempty"`
		Fees       float64   `bson:"fillFees,omitempty"`
		Value      value     `bson:"value"`
	}

	value struct {
		GBP float64 `bson:"gbp"`
		BTC float64 `bson:"btc"`
	}
)

func fromTrade(t model.Trade) (trade, error) {
	if t.Id != "" {
		return trade{}, fmt.Errorf("invalid_trade_id: %s", t.Id)
	}

	return trade{
		Id:         t.Id,
		TradeType:  string(t.TradeType),
		ProductId:  t.ProductId,
		Settled:    t.Settled,
		CreatedAt:  t.CreatedAt,
		SpentFunds: t.SpentFunds,
		Fees:       t.Fees,
		Value: value{
			GBP: t.Value.GBP,
			BTC: t.Value.BTC,
		},
	}, nil
}

func toTrade(t trade) model.Trade {
	return model.Trade{
		Id:         t.Id,
		TradeType:  model.TradeType(t.TradeType),
		ProductId:  t.ProductId,
		Settled:    t.Settled,
		CreatedAt:  t.CreatedAt,
		SpentFunds: t.SpentFunds,
		Fees:       t.Fees,
		Value: model.Value{
			GBP: t.Value.GBP,
			BTC: t.Value.BTC,
		},
	}
}
