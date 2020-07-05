module github.com/cshep4/kripto/shared/go/idempotency

go 1.14

require (
	github.com/aws/aws-lambda-go v1.17.0
	github.com/cshep4/kripto/shared/go/log v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.4.3
	github.com/stretchr/testify v1.4.0
	go.mongodb.org/mongo-driver v1.3.3
)

replace github.com/cshep4/kripto/shared/go/log => ../log
