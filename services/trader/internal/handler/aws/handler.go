package aws

import (
	"context"
	"fmt"
	"github.com/cshep4/kripto/services/trader/internal/model"
	"strconv"

	"github.com/cshep4/kripto/shared/go/log"
)

type (
	Servicer interface {
		Trade(ctx context.Context, tradeType, amount string) error
		GetWallet(ctx context.Context) (*model.Wallet, error)
	}

	Handler struct {
		Service Servicer
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	BadRequestError struct {
		Parameter string
		Err       string
	}

	TradeRequest struct {
		IdempotencyKey string `json:"idempotencyKey"`
		TradeType      string `json:"tradeType"`
		Amount         string `json:"amount"`
	}
)

func (i BadRequestError) Error() string {
	return fmt.Sprintf("bad request - param: %s, error: %s", i.Parameter, i.Err)
}

func (h *Handler) Trade(ctx context.Context, req TradeRequest) error {
	switch {
	case req.TradeType == "":
		return BadRequestError{Parameter: "tradeType", Err: "empty"}
	case req.TradeType != "buy" && req.TradeType != "sell":
		return BadRequestError{Parameter: "tradeType", Err: "invalid value - should be either buy/sell"}
	case req.Amount == "":
		return BadRequestError{Parameter: "amount", Err: "empty"}
	}

	a, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return BadRequestError{Parameter: "amount", Err: "invalid value - should be numeric"}
	}

	err = h.Service.Trade(ctx, req.TradeType, fmt.Sprintf("%.2f", a))
	if err != nil {
		log.Error(ctx, "error_trading",
			log.ErrorParam(err),
			log.SafeParam("tradeType", req.TradeType),
			log.SafeParam("amount", req.Amount),
		)
		return fmt.Errorf("trade: %w", err)
	}

	return nil
}

func (h *Handler) GetWallet(ctx context.Context) (*model.Wallet, error) {
	wallet, err := h.Service.GetWallet(ctx)
	if err != nil {
		log.Error(ctx, "error_getting_wallet", log.ErrorParam(err))
		return nil, fmt.Errorf("get_wallet: %w", err)
	}

	return wallet, nil
}
