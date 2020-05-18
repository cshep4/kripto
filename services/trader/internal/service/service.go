package service

import (
	"context"
	"fmt"
	"github.com/cshep4/kripto/services/trader/internal/trader"
)

type (
	Trader interface {
		Trade(tradeType trader.TradeType, amount string) (*trader.TradeResponse, error)
	}

	service struct {
		amount string
		trader Trader
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func New(amount string, trader Trader) (*service, error) {
	switch {
	case amount == "":
		return nil, InvalidParameterError{Parameter: "amount"}
	case trader == nil:
		return nil, InvalidParameterError{Parameter: "trader"}
	}

	return &service{
		amount: amount,
		trader: trader,
	}, nil
}

func (s *service) Trade(ctx context.Context, tradeType string) error {
	_, err := s.trader.Trade(trader.TradeType(tradeType), s.amount)
	if err != nil {
		return fmt.Errorf("trade: %w", err)
	}

	return nil
}
