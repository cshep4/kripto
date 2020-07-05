package model

type (
	Wallet struct {
		GBP Account `json:"gbp"`
		BTC Account `json:"btc"`
	}

	Account struct {
		ID        string  `json:"id"`
		Balance   float32 `json:"balance"`
		Hold      float32 `json:"hold"`
		Available float32 `json:"available"`
	}
)
