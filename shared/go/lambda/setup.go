package lambda

import (
	"context"
	"github.com/cshep4/kripto/shared/go/log"
)

type (
	SetupFunc       func(context.Context) error
	InitialisedFunc func() bool

	FunctionConfig struct {
		Setup        SetupFunc
		Initialised  InitialisedFunc
		LogLevel     string
		ServiceName  string
		FunctionName string
	}
)

func Init(cfg FunctionConfig) {
	if cfg.Initialised() {
		return
	}

	ctx := log.WithFunctionName(context.Background(),
		log.New(cfg.LogLevel),
		cfg.ServiceName,
		cfg.FunctionName,
	)

	log.Info(ctx, "initialisation")

	if err := cfg.Setup(ctx); err != nil {
		log.Fatal(ctx, "initialisation_error", log.ErrorParam(err))
	}
}
