package lambda

//go:generate mockgen -destination=internal/mocks/aws/handler.gen.go -package=aws_mocks github.com/aws/aws-lambda-go/lambda Handler
