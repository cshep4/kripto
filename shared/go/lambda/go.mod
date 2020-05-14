module github.com/cshep4/kripto/shared/go/lambda

go 1.14

replace github.com/cshep4/kripto/shared/go/log => ../log

require (
	github.com/aws/aws-lambda-go v1.15.0
	github.com/cshep4/kripto/shared/go/log v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.4.0
	github.com/stretchr/testify v1.4.0
)
