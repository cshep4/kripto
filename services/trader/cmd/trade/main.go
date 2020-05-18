package main

import (
	"context"
	"fmt"
	"github.com/Netflix/go-env"
	"github.com/cshep4/kripto/services/trader/internal/handler/aws"
	"github.com/cshep4/kripto/services/trader/internal/service"
	"github.com/cshep4/kripto/services/trader/internal/trader"
	"github.com/cshep4/kripto/shared/go/lambda"
	"github.com/preichenberger/go-coinbasepro/v2"
)

type environment struct {
	CoinbasePro struct {
		Key        string `env:"COINBASE_PRO_KEY"`
		Passphrase string `env:"COINBASE_PRO_PASSPHRASE"`
		Secret     string `env:"COINBASE_PRO_SECRET"`
	}
	CoinbaseProSandbox struct {
		Key        string `env:"COINBASE_PRO_SANDBOX_KEY"`
		Passphrase string `env:"COINBASE_PRO_SANDBOX_PASSPHRASE"`
		Secret     string `env:"COINBASE_PRO_SANDBOX_SECRET"`
	}
	MockTrade   bool   `env:"MOCK_TRADE"`
	TradeAmount string `env:"TRADE_AMOUNT"`
}

var (
	cfg = lambda.FunctionConfig{
		LogLevel:     "info",
		ServiceName:  "trader",
		FunctionName: "read",
		Setup:        setup,
		Initialised:  func() bool { return handler.Service != nil },
	}

	handler = &aws.Handler{}

	runner = lambda.New(
		handler.Trade,
		lambda.WithPreExecute(lambda.LogMiddleware(cfg.LogLevel, cfg.ServiceName, cfg.FunctionName)),
	)
)

func main() {
	lambda.Init(cfg)
	runner.Start()
}

func setup(ctx context.Context) error {
	env, err := getEnv()
	if err != nil {
		return err
	}

	coinbaseClient := initCoinbaseProClient(env)

	trader, err := trader.New(coinbaseClient)
	if err != nil {
		return fmt.Errorf("initialise_trader: %w", err)
	}

	handler.Service, err = service.New(env.TradeAmount, trader)
	if err != nil {
		return fmt.Errorf("initialise_service: %w", err)
	}

	return nil
}

func getEnv() (environment, error) {
	var e environment
	_, err := env.UnmarshalFromEnviron(&e)
	if err != nil {
		return environment{}, fmt.Errorf("unmarshal_environment_variables: %w", err)
	}
	switch {
	case !e.MockTrade && e.CoinbasePro.Key == "":
		return environment{}, fmt.Errorf("missing_environment_variable: COINBASE_PRO_KEY")
	case !e.MockTrade && e.CoinbasePro.Passphrase == "":
		return environment{}, fmt.Errorf("missing_environment_variable: COINBASE_PRO_PASSPHRASE")
	case !e.MockTrade && e.CoinbasePro.Secret == "":
		return environment{}, fmt.Errorf("missing_environment_variable: COINBASE_PRO_SECRET")
	case e.MockTrade && e.CoinbaseProSandbox.Key == "":
		return environment{}, fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_KEY")
	case e.MockTrade && e.CoinbaseProSandbox.Passphrase == "":
		return environment{}, fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_PASSPHRASE")
	case e.MockTrade && e.CoinbaseProSandbox.Secret == "":
		return environment{}, fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_SECRET")
	case e.TradeAmount == "":
		return environment{}, fmt.Errorf("missing_environment_variable: TRADE_AMOUNT")
	}
	return e, err
}

func initCoinbaseProClient(env environment) *coinbasepro.Client {
	coinbaseClient := coinbasepro.NewClient()
	if env.MockTrade {
		coinbaseClient.UpdateConfig(&coinbasepro.ClientConfig{
			BaseURL:    "https://api-public.sandbox.pro.coinbase.com",
			Key:        env.CoinbaseProSandbox.Key,
			Passphrase: env.CoinbaseProSandbox.Passphrase,
			Secret:     env.CoinbaseProSandbox.Secret,
		})
	}
	return coinbaseClient
}
