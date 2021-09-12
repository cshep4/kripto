package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cshep4/go-log"
	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"go.uber.org/zap"
)

type (
	Servicer interface {
		Get(ctx context.Context) ([]model.Rate, error)
		StoreTrade(ctx context.Context, trade model.Trade) error
		StoreRate(ctx context.Context, rate float64, dateTime time.Time) error
	}

	Handler struct {
		Service Servicer
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func (h *Handler) IsInitialised() bool {
	return h.Service != nil
}

func (h *Handler) Functions() map[string]interface{} {
	return map[string]interface{}{
		"data-reader":  h.Get,
		"trade-writer": h.StoreTrade,
		"rate-writer":  h.StoreRate,
	}
}

func (h *Handler) Get(ctx context.Context) ([]model.Rate, error) {
	rates, err := h.Service.Get(ctx)
	if err != nil {
		log.Error(ctx, "error_getting_data", zap.Error(err))
		return nil, err
	}

	return rates, nil
}

func (h *Handler) StoreTrade(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("no sqs message passed to function")
	}

	for _, msg := range sqsEvent.Records {
		var req model.TradeRequest
		err := json.Unmarshal([]byte(msg.Body), &req)
		if err != nil {
			log.Error(ctx, "invalid_msg_body", zap.Error(err))
			continue
		}

		trade, err := req.ToTrade()
		if err != nil {
			log.Error(ctx, "invalid_msg_body",
				zap.String("id", req.Id),
				zap.String("funds", req.Funds),
				zap.String("btc", req.FilledSize),
				zap.String("gbp", req.ExecutedValue),
				zap.Error(err),
			)
			continue
		}

		err = h.Service.StoreTrade(ctx, trade)
		if err != nil {
			log.Error(ctx, "error_storing_trade",
				zap.String("id", trade.Id),
				zap.Time("createdAt", trade.CreatedAt),
				zap.Float64("btc", trade.Value.BTC),
				zap.Float64("gbp", trade.Value.GBP),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}

func (h *Handler) StoreRate(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("no sqs message passed to function")
	}

	for _, msg := range sqsEvent.Records {
		var req model.StoreRateRequest
		err := json.Unmarshal([]byte(msg.Body), &req)
		if err != nil {
			log.Error(ctx, "invalid_msg_body", zap.Error(err))
			continue
		}

		err = h.Service.StoreRate(ctx, req.Rate, req.DateTime)
		if err != nil {
			log.Error(ctx, "error_storing_rate",
				zap.Float64("rate", req.Rate),
				zap.Time("dateTime", req.DateTime),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}
