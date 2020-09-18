package main

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cshep4/kripto/services/trader/internal/handler/aws"
	"github.com/cshep4/kripto/services/trader/internal/secrets"
	"github.com/cshep4/kripto/services/trader/internal/service"
	"github.com/cshep4/kripto/services/trader/internal/trader"
	"github.com/cshep4/kripto/shared/go/idempotency"
	idempotent "github.com/cshep4/kripto/shared/go/idempotency/middleware"
	"github.com/cshep4/kripto/shared/go/idempotency/middleware/invoke"
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/cshep4/kripto/shared/go/log"
	"github.com/cshep4/kripto/shared/go/mongodb"
	"github.com/preichenberger/go-coinbasepro/v2"
)

const (
	logLevel     = "info"
	serviceName  = "trader"
	functionName = "trade"
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
		handler.Trade,
		lambda.WithPreExecute(log.Middleware(logLevel, serviceName, functionName)),
	)
)

func main() {
	runner.Start(cfg)
}

func setup(ctx context.Context) error {
	var s secrets.Secrets
	if err := s.Fetch(); err != nil {
		return err
	}

	coinbaseClient := initCoinbaseProClient(s)

	sess, err := session.NewSession(&awsconfig.Config{
		Region: &s.SNS.Region,
	})
	if err != nil {
		return fmt.Errorf("new_session: %w", err)
	}

	publisher := sns.New(sess)

	trader, err := trader.New(coinbaseClient)
	if err != nil {
		return fmt.Errorf("initialise_trader: %w", err)
	}

	handler.Service, err = service.New(s.SNS.Topic, publisher, trader)
	if err != nil {
		return fmt.Errorf("initialise_service: %w", err)
	}

	mongoClient, err := mongodb.New(ctx)
	if err != nil {
		return fmt.Errorf("initialise_mongo_client: %w", err)
	}

	idempotencer, err := idempotency.New(ctx, "trade", mongoClient)
	if err != nil {
		return fmt.Errorf("initialise_idempotencer: %w", err)
	}

	middleware, err = invoke.NewMiddleware(idempotencer)
	if err != nil {
		return fmt.Errorf("initialise_idempotency_middleware: %w", err)
	}

	runner.Apply(
		lambda.WithPreExecute(middleware.PreExecute),
		lambda.WithPostExecute(middleware.PostExecute),
		lambda.WithErrorHandler(middleware.HandleError),
	)

	return nil
}

func initCoinbaseProClient(s secrets.Secrets) *coinbasepro.Client {
	coinbaseClient := coinbasepro.NewClient()
	if s.MockTrade {
		coinbaseClient.UpdateConfig(&coinbasepro.ClientConfig{
			BaseURL:    "https://api-public.sandbox.pro.coinbase.com",
			Key:        s.CoinbasePro.Sandbox.Key,
			Passphrase: s.CoinbasePro.Sandbox.Passphrase,
			Secret:     s.CoinbasePro.Sandbox.Secret,
		})
	}
	return coinbaseClient
}
