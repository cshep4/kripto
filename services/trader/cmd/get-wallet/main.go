package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/cshep4/kripto/services/trader/internal/handler/aws"
	"github.com/cshep4/kripto/services/trader/internal/secrets"
	"github.com/cshep4/kripto/services/trader/internal/service"
	"github.com/cshep4/kripto/services/trader/internal/trader"
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/cshep4/kripto/shared/go/log"
	"github.com/preichenberger/go-coinbasepro/v2"
)

const (
	logLevel     = "info"
	serviceName  = "trader"
	functionName = "get-wallet"
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
		handler.GetWallet,
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

	trader, err := trader.New(coinbaseClient)
	if err != nil {
		return fmt.Errorf("initialise_trader: %w", err)
	}

	handler.Service, err = service.New("topic", &sns.SNS{}, trader)
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
			Key:        s.CoinbasePro.Sandbox.Key,
			Passphrase: s.CoinbasePro.Sandbox.Passphrase,
			Secret:     s.CoinbasePro.Sandbox.Secret,
		})
	}
	return coinbaseClient
}
