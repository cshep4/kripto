package main

import (
	"context"

	"github.com/cshep4/kripto/services/data-storer/internal/handler/aws"
	"github.com/cshep4/kripto/services/data-storer/internal/service"
	"github.com/cshep4/kripto/services/data-storer/internal/store/rate/mongo"
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/cshep4/kripto/shared/go/mongodb"
)

var (
	cfg = lambda.FunctionConfig{
		LogLevel:     "info",
		ServiceName:  "data-storer",
		FunctionName: "read",
		Setup:        setup,
	}

	handler *aws.Handler

	runner = lambda.New(
		handler.StoreRate,
		lambda.WithPreExecute(lambda.LogMiddleware(cfg.LogLevel, cfg.ServiceName, cfg.FunctionName)),
	)
)

func setup(ctx context.Context) error {
	mongoClient, err := mongodb.New(ctx)
	if err != nil {
		return err
	}

	store, err := mongo.New(ctx, mongoClient)
	if err != nil {
		return err
	}

	svc, err := service.New(store)
	if err != nil {
		return err
	}

	handler, err = aws.New(svc)
	return err
}

func main() {
	lambda.Init(handler, cfg)
	runner.Start()
}
