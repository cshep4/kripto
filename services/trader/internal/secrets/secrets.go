package secrets

import (
	"fmt"

	"github.com/Netflix/go-env"
)

type Secrets struct {
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
	SNS struct {
		Topic  string `env:"TOPIC"`
		Region string `env:"REGION"`
	}
	MockTrade   bool   `env:"MOCK_TRADE"`
	TradeAmount string `env:"TRADE_AMOUNT"`
}

func (s *Secrets) Fetch() error {
	_, err := env.UnmarshalFromEnviron(&s)
	if err != nil {
		return fmt.Errorf("unmarshal_environment_variables: %w", err)
	}
	switch {
	case !s.MockTrade && s.CoinbasePro.Key == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_KEY")
	case !s.MockTrade && s.CoinbasePro.Passphrase == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_PASSPHRASE")
	case !s.MockTrade && s.CoinbasePro.Secret == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_SECRET")
	case s.MockTrade && s.CoinbaseProSandbox.Key == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_KEY")
	case s.MockTrade && s.CoinbaseProSandbox.Passphrase == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_PASSPHRASE")
	case s.MockTrade && s.CoinbaseProSandbox.Secret == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_SECRET")
	case s.TradeAmount == "":
		return fmt.Errorf("missing_environment_variable: TRADE_AMOUNT")
	}
	return err
}
