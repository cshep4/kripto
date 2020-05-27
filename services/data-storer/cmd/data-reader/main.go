package main

import (
	"context"
	"fmt"

	"github.com/cshep4/kripto/services/data-storer/internal/handler/aws"
	"github.com/cshep4/kripto/services/data-storer/internal/service"
	rate "github.com/cshep4/kripto/services/data-storer/internal/store/rate/mongo"
	trade "github.com/cshep4/kripto/services/data-storer/internal/store/trade/mongo"
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/cshep4/kripto/shared/go/mongodb"
)

var (
	cfg = lambda.FunctionConfig{
		LogLevel:     "info",
		ServiceName:  "data-storer",
		FunctionName: "data-reader",
		Setup:        setup,
		Initialised:  func() bool { return handler.Service != nil },
	}

	handler = &aws.Handler{}

	runner = lambda.New(
		handler.Get,
		lambda.WithPreExecute(lambda.LogMiddleware(cfg.LogLevel, cfg.ServiceName, cfg.FunctionName)),
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

	return nil
}
