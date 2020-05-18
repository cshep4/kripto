package aws

import (
	"context"
	"fmt"

	"github.com/cshep4/kripto/shared/go/log"
)

type (
	Servicer interface {
		Trade(ctx context.Context, tradeType string) error
	}

	Handler struct {
		Service Servicer
	}

	// InvalidParameterError is returned when a required parameter passed to New is invalid.
	InvalidParameterError struct {
		Parameter string
	}

	TradeRequest struct {
		TradeType string `json:"tradeType"`
	}
)

func (i InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid parameter %s", i.Parameter)
}

func (h *Handler) Trade(ctx context.Context, req TradeRequest) error {
	err := h.Service.Trade(ctx, req.TradeType)
	if err != nil {
		log.Error(ctx, "error_trading", log.ErrorParam(err))
		return fmt.Errorf("trade: %w", err)
	}

	return nil
}