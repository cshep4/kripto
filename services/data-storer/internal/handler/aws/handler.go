package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"github.com/cshep4/kripto/shared/go/idempotency"
	"github.com/cshep4/kripto/shared/go/log"
)

type (
	Servicer interface {
		Get(ctx context.Context) (*model.GetResponse, error)
		StoreTrade(ctx context.Context, trade model.Trade) error
		StoreRate(ctx context.Context, rate float64, dateTime time.Time) error
	}

	Handler struct {
		Service      Servicer
		Idempotencer idempotency.Idempotencer
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func (h *Handler) Get(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func (h *Handler) StoreTrade(ctx context.Context, sqsEvent events.SQSEvent) error {
	if len(sqsEvent.Records) == 0 {
		return errors.New("no sqs message passed to function")
	}

	for _, msg := range sqsEvent.Records {
		var req model.TradeRequest
		err := json.Unmarshal([]byte(msg.Body), &req)
		if err != nil {
			log.Error(ctx, "invalid_msg_body", log.ErrorParam(err))
			continue
		}

		trade, err := req.ToTrade()
		if err != nil {
			log.Error(ctx, "invalid_msg_body",
				log.SafeParam("id", req.Id),
				log.SafeParam("funds", req.Funds),
				log.SafeParam("btc", req.FilledSize),
				log.SafeParam("gbp", req.ExecutedValue),
				log.ErrorParam(err),
			)
			continue
		}

		ok, err := h.Idempotencer.Check(ctx, req.Id)
		if err != nil {
			log.Error(ctx, "error_checking_idempotency", log.ErrorParam(err))
			continue
		}
		if ok {
			log.Info(ctx, "msg_already_processed",
				log.SafeParam("id", req.Id),
				log.SafeParam("createdAt", req.CreatedAt),
				log.SafeParam("btc", req.FilledSize),
				log.SafeParam("gbp", req.ExecutedValue),
			)
			continue
		}

		err = h.Service.StoreTrade(ctx, trade)
		if err != nil {
			log.Error(ctx, "error_storing_trade",
				log.SafeParam("id", req.Id),
				log.SafeParam("createdAt", req.CreatedAt),
				log.SafeParam("btc", req.FilledSize),
				log.SafeParam("gbp", req.ExecutedValue),
				log.ErrorParam(err),
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
			log.Error(ctx, "invalid_msg_body", log.ErrorParam(err))
			continue
		}

		ok, err := h.Idempotencer.Check(ctx, req.IdempotencyKey)
		if err != nil {
			log.Error(ctx, "error_checking_idempotency", log.ErrorParam(err))
			continue
		}
		if ok {
			log.Info(ctx, "msg_already_processed",
				log.SafeParam("rate", req.Rate),
				log.SafeParam("dateTime", req.DateTime),
				log.SafeParam("key", req.IdempotencyKey),
			)
			continue
		}

		err = h.Service.StoreRate(ctx, req.Rate, req.DateTime)
		if err != nil {
			log.Error(ctx, "error_storing_rate",
				log.SafeParam("rate", req.Rate),
				log.SafeParam("dateTime", req.DateTime),
				log.SafeParam("key", req.IdempotencyKey),
				log.ErrorParam(err),
			)
			return err
		}
	}

	return nil
}
