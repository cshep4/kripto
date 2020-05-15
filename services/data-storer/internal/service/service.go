package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cshep4/kripto/services/data-storer/internal/model"
)

type (
	RateStore interface {
		Store(ctx context.Context, rate float64, dateTime time.Time) error
		GetPreviousWeeks(ctx context.Context) ([]model.Rate, error)
	}

	TradeStore interface {
		Store(ctx context.Context, trade model.Trade) error
		GetPreviousWeeks(ctx context.Context) ([]model.Trade, error)
	}

	service struct {
		rateStore  RateStore
		tradeStore TradeStore
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func New(rateStore RateStore, tradeStore TradeStore) (*service, error) {
	if rateStore == nil {
		return nil, InvalidParameterError{Parameter: "rateStore"}
	}
	if tradeStore == nil {
		return nil, InvalidParameterError{Parameter: "tradeStore"}
	}

	return &service{
		rateStore:  rateStore,
		tradeStore: tradeStore,
	}, nil
}

func (s *service) Get(ctx context.Context) (*model.GetResponse, error) {
	panic("implement me")
}

func (s *service) Store(ctx context.Context, req model.StoreRequest) error {
	panic("implement me")
}

func (s *service) StoreRate(ctx context.Context, rate float64, dateTime time.Time) error {
	err := s.rateStore.Store(ctx, rate, dateTime)
	if err != nil {
		return fmt.Errorf("store_rate: %w", err)
	}

	return nil
}
