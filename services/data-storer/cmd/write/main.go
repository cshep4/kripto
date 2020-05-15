package main

import (
	"context"

	"github.com/cshep4/kripto/services/data-storer/internal/handler/aws"
	"github.com/cshep4/kripto/shared/go/lambda"
)

var (
	cfg = lambda.FunctionConfig{
		LogLevel:     "info",
		ServiceName:  "data-storer",
		FunctionName: "write",
		Setup:        setup,
		Initialised:  func() bool { return handler != nil },
	}

	handler *aws.Handler

	runner = lambda.New(
		handler.Store,
		lambda.WithPreExecute(lambda.LogMiddleware(cfg.LogLevel, cfg.ServiceName, cfg.FunctionName)),
	)
)

func setup(ctx context.Context) error {
	//mongoClient, err := mongodb.New(ctx)
	//if err != nil {
	//	return err
	//}
	//
	//store, err := mongo.New(ctx, mongoClient)
	//if err != nil {
	//	return err
	//}
	//
	//svc, err := service.New(store)
	//if err != nil {
	//	return err
	//}
	//
	//handler, err = aws.New(svc)
	//return err
	return nil
}

func main() {
	lambda.Init(cfg)
	runner.Start()
}
