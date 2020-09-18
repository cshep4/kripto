package secrets

import (
	"fmt"

	"github.com/Netflix/go-env"
)

type Secrets struct {
	CoinbasePro struct {
		Live struct {
			Key        string `env:"COINBASE_PRO_KEY"`
			Passphrase string `env:"COINBASE_PRO_PASSPHRASE"`
			Secret     string `env:"COINBASE_PRO_SECRET"`
		}
		Sandbox struct {
			Key        string `env:"COINBASE_PRO_SANDBOX_KEY"`
			Passphrase string `env:"COINBASE_PRO_SANDBOX_PASSPHRASE"`
			Secret     string `env:"COINBASE_PRO_SANDBOX_SECRET"`
		}
	}
	SNS struct {
		Topic  string `env:"TOPIC"`
		Region string `env:"REGION"`
	}
	MockTrade bool `env:"MOCK_TRADE"`
}

func (s *Secrets) Fetch() error {
	_, err := env.UnmarshalFromEnviron(s)
	if err != nil {
		return fmt.Errorf("unmarshal_environment_variables: %w", err)
	}
	switch {
	case !s.MockTrade && s.CoinbasePro.Live.Key == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_KEY")
	case !s.MockTrade && s.CoinbasePro.Live.Passphrase == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_PASSPHRASE")
	case !s.MockTrade && s.CoinbasePro.Live.Secret == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_SECRET")
	case s.MockTrade && s.CoinbasePro.Sandbox.Key == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_KEY")
	case s.MockTrade && s.CoinbasePro.Sandbox.Passphrase == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_PASSPHRASE")
	case s.MockTrade && s.CoinbasePro.Sandbox.Secret == "":
		return fmt.Errorf("missing_environment_variable: COINBASE_PRO_SANDBOX_SECRET")
	}
	return err
}
