module github.com/cshep4/kripto/services/trader

go 1.14

require (
	github.com/Netflix/go-env v0.0.0-20200512170851-5660fe1ab40a
	github.com/aws/aws-sdk-go v1.31.0
	github.com/cshep4/kripto/shared/go/idempotency v0.0.0-00010101000000-000000000000
	github.com/cshep4/kripto/shared/go/lambda v0.0.0-00010101000000-000000000000
	github.com/cshep4/kripto/shared/go/log v0.0.0-00010101000000-000000000000
	github.com/cshep4/kripto/shared/go/mongodb v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.4.3
	github.com/preichenberger/go-coinbasepro/v2 v2.0.5
	github.com/stretchr/testify v1.5.1
)

replace github.com/cshep4/kripto/shared/go/mongodb => ../../shared/go/mongodb

replace github.com/cshep4/kripto/shared/go/log => ../../shared/go/log

replace github.com/cshep4/kripto/shared/go/lambda => ../../shared/go/lambda

replace github.com/cshep4/kripto/shared/go/idempotency => ../../shared/go/idempotency
