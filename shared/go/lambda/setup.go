package lambda

import (
	"context"
	"reflect"

	"github.com/cshep4/kripto/shared/go/log"
)

type (
	Handler interface {
		ServiceName() string
	}
	SetupFunc func(context.Context) error

	FunctionConfig struct {
		Setup        SetupFunc
		LogLevel     string
		ServiceName  string
		FunctionName string
	}
)

func Init(handler Handler, cfg FunctionConfig) {
	if !reflect.ValueOf(handler).IsNil() {
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
