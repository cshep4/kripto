package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cshep4/lambda-go/lambda"
	"github.com/cshep4/lambda-go/log/v2"
	"github.com/cshep4/lambda-go/mongodb"

	"github.com/cshep4/kripto/services/data-storer/internal/handler/aws"
	"github.com/cshep4/kripto/services/data-storer/internal/service"
	rate "github.com/cshep4/kripto/services/data-storer/internal/store/rate/mongo"
	trade "github.com/cshep4/kripto/services/data-storer/internal/store/trade/mongo"
)

const (
	logLevel    = "info"
	serviceName = "data-storer"
)

var (
	functionName = os.Getenv("FUNCTION_NAME")

	handler = &aws.Handler{}

	runner = lambda.New(
		functionName,
		handler,
		lambda.WithServiceName(serviceName),
		lambda.WithLogLevel(logLevel),
		lambda.WithPreExecute(log.Middleware(logLevel, serviceName, functionName)),
	)
)

func main() {
	runner.Start(setup)
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
