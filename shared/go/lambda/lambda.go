package lambda

import (
	"sort"

	"github.com/aws/aws-lambda-go/lambda"
)

type runner struct {
	lambda.Handler
	terminated  bool
	baseHandler lambda.Handler
	opts        []option
}

// New creates a new AWS lambda runner
func New(handler interface{}, opts ...option) runner {
	h := lambda.NewHandler(handler)
	r := runner{
		baseHandler: h,
	}
	r.Apply(opts...)

	return r
}

// Start beings execution of the AWS lambda function
func (r *runner) Start() {
	lambda.StartHandler(r.Handler)
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
}
