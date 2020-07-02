package main

import (
	"context"
	"fmt"

	"github.com/cshep4/kripto/services/data-storer/internal/handler/aws"
	"github.com/cshep4/kripto/services/data-storer/internal/service"
	rate "github.com/cshep4/kripto/services/data-storer/internal/store/rate/mongo"
	trade "github.com/cshep4/kripto/services/data-storer/internal/store/trade/mongo"
	"github.com/cshep4/kripto/shared/go/idempotency"
	idempotent "github.com/cshep4/kripto/shared/go/idempotency/middleware"
	"github.com/cshep4/kripto/shared/go/idempotency/middleware/sqs"
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/cshep4/kripto/shared/go/log"
	"github.com/cshep4/kripto/shared/go/mongodb"
)

const (
	logLevel     = "info"
	serviceName  = "data-storer"
	functionName = "trade-writer"
)

var (
	cfg = lambda.FunctionConfig{
		LogLevel:     logLevel,
		ServiceName:  serviceName,
		FunctionName: functionName,
		Setup:        setup,
		Initialised:  func() bool { return handler.Service != nil && middleware != nil },
	}

	handler    aws.Handler
	middleware idempotent.Middleware

	runner = lambda.New(
		handler.StoreTrade,
		lambda.WithPreExecute(log.Middleware(cfg.LogLevel, cfg.ServiceName, cfg.FunctionName)),
	)
)

func main() {
	lambda.Init(cfg)
	runner.Start()
}

func setup(ctx context.Context) error {
	mongoClient, err := mongodb.New(ctx)
	if err != nil {
		return fmt.Errorf("initialise_mongo_client: %w", err)
	}

	rateStore, err := rate.New(ctx, mongoClient)
	if err != nil {
		return fmt.Errorf("initialise_rate_store: %w", err)
	}

	tradeStore, err := trade.New(ctx, mongoClient)
	if err != nil {
		return fmt.Errorf("initialise_trade_store: %w", err)
	}

	handler.Service, err = service.New(rateStore, tradeStore)
	if err != nil {
		return fmt.Errorf("initialise_service: %w", err)
	}

	idempotencer, err := idempotency.New(ctx, "trade", mongoClient)
	if err != nil {
		return fmt.Errorf("initialise_idempotencer: %w", err)
	}

	middleware, err = sqs.NewMiddleware(idempotencer)
	if err != nil {
		return fmt.Errorf("initialise_sqs_idempotency_middleware: %w", err)
	}

	runner.Apply(
		lambda.WithPreExecute(middleware.PreExecute),
		lambda.WithPostExecute(middleware.PostExecute),
		lambda.WithErrorHandler(middleware.HandleError),
	)

	return nil
}
