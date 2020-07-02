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
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/cshep4/kripto/shared/go/log"
	"github.com/preichenberger/go-coinbasepro/v2"
)

var (
	cfg = lambda.FunctionConfig{
		LogLevel:     "info",
		ServiceName:  "trader",
		FunctionName: "trade",
		Setup:        setup,
		Initialised:  func() bool { return handler.Service != nil },
	}

	handler = &aws.Handler{}

	runner = lambda.New(
		handler.Trade,
		lambda.WithPreExecute(log.Middleware(cfg.LogLevel, cfg.ServiceName, cfg.FunctionName)),
	)
)

func main() {
	lambda.Init(cfg)
	runner.Start()
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

	publisher := sns.New(sess)

	trader, err := trader.New(coinbaseClient)
	if err != nil {
		return fmt.Errorf("initialise_trader: %w", err)
	}

	handler.Service, err = service.New(s.TradeAmount, s.SNS.Topic, publisher, trader)
	if err != nil {
		return fmt.Errorf("initialise_service: %w", err)
	}

	return nil
}

func initCoinbaseProClient(s secrets.Secrets) *coinbasepro.Client {
	coinbaseClient := coinbasepro.NewClient()
	if s.MockTrade {
		coinbaseClient.UpdateConfig(&coinbasepro.ClientConfig{
			BaseURL:    "https://api-public.sandbox.pro.coinbase.com",
			Key:        s.CoinbaseProSandbox.Key,
			Passphrase: s.CoinbaseProSandbox.Passphrase,
			Secret:     s.CoinbaseProSandbox.Secret,
		})
	}
	return coinbaseClient
}
