module github.com/cshep4/kripto/services/data-storer

go 1.14

require (
	github.com/aws/aws-lambda-go v1.15.0
	github.com/cshep4/kripto/shared/go/lambda v0.0.0-00010101000000-000000000000
	github.com/cshep4/kripto/shared/go/mongodb v0.0.0-20200227002205-4db5d3107521
	go.mongodb.org/mongo-driver v1.3.1
)

replace github.com/cshep4/kripto/shared/go/mongodb => ../../shared/go/mongodb

replace github.com/cshep4/kripto/shared/go/log => ../../shared/go/log

replace github.com/cshep4/kripto/shared/go/lambda => ../../shared/go/lambda
