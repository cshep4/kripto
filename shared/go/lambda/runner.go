package lambda

import (
	"context"
	"sort"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cshep4/kripto/shared/go/log"
)

type (
	FunctionConfig struct {
		Setup        func(context.Context) error
		Initialised  func() bool
		LogLevel     string
		ServiceName  string
		FunctionName string
	}

	runner struct {
		lambda.Handler
		terminatedRequests map[string]struct{}
		lock               sync.Mutex
		baseHandler        lambda.Handler
		opts               []option
	}
)

// New creates a new AWS lambda runner
func New(handler interface{}, opts ...option) runner {
	h := lambda.NewHandler(handler)
	r := runner{
		terminatedRequests: make(map[string]struct{}),
		baseHandler:        h,
	}
	r.Apply(opts...)

	return r
}

// Apply will apply all runner options from scratch along
// with specified opts to ensure correct ordering
func (r *runner) Apply(opts ...option) {
	r.opts = append(opts, r.opts...)
	r.Handler = r.baseHandler

	sort.Slice(r.opts, func(i, j int) bool {
		return r.opts[i].optionType < r.opts[j].optionType
	})

	for _, opt := range r.opts {
		opt.apply(r)
	}

	r.Handler = &initialiser{
		handler: r.Handler,
		runner:  r,
	}
}

// Start beings execution of the AWS lambda function
func (r *runner) Start(cfg FunctionConfig) {
	if !cfg.Initialised() {
		r.init(cfg)
	}
	lambda.StartHandler(r)
}

func (r *runner) init(cfg FunctionConfig) {
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

func (r *runner) finish(ctx context.Context) {
	r.lock.Lock()
	defer r.lock.Unlock()

	i := ctx.Value(requestId)
	if i == nil {
		return
	}
	id, ok := i.(string)
	if !ok {
		return
	}

	r.terminatedRequests[id] = struct{}{}
}

func (r *runner) done(ctx context.Context) bool {
	i := ctx.Value(requestId)
	if i == nil {
		return false
	}
	id, ok := i.(string)
	if !ok {
		return false
	}
	_, ok = r.terminatedRequests[id]

	return ok
}
