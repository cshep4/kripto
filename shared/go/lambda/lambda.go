package lambda

import (
	"github.com/aws/aws-lambda-go/lambda"
)

type Runner struct {
	lambda.Handler
}

func New(handler interface{}, opts ...option) Runner {
	r := Runner{
		Handler: lambda.NewHandler(handler),
	}

	for _, opt := range opts {
		opt(&r)
	}

	return r
}

func (r *Runner) Start() {
	lambda.StartHandler(r.Handler)
}
