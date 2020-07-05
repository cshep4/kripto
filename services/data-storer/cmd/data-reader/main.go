package main

import (
	"context"
	"fmt"

	"github.com/cshep4/kripto/services/data-storer/internal/handler/aws"
	"github.com/cshep4/kripto/services/data-storer/internal/service"
	rate "github.com/cshep4/kripto/services/data-storer/internal/store/rate/mongo"
	trade "github.com/cshep4/kripto/services/data-storer/internal/store/trade/mongo"
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/cshep4/kripto/shared/go/log"
	"github.com/cshep4/kripto/shared/go/mongodb"
)

const (
	logLevel     = "info"
	serviceName  = "data-storer"
	functionName = "data-reader"
)

var (
	cfg = lambda.FunctionConfig{
		LogLevel:     logLevel,
		ServiceName:  serviceName,
		FunctionName: functionName,
		Setup:        setup,
		Initialised:  func() bool { return handler.Service != nil },
	}

	handler aws.Handler

	runner = lambda.New(
		handler.Get,
		lambda.WithPreExecute(log.Middleware(logLevel, serviceName, functionName)),
	)
)

func main() {
	runner.Start(cfg)
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
