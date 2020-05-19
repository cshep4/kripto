package trader

//go:generate mockgen -destination=internal/mocks/service/servicer.gen.go -package=service_mocks github.com/cshep4/kripto/services/trader/internal/handler/aws Servicer
//go:generate mockgen -destination=internal/mocks/trader/trader.gen.go -package=trader_mocks github.com/cshep4/kripto/services/trader/internal/service Trader
//go:generate mockgen -destination=internal/mocks/aws/sqs.gen.go -package=aws_mocks github.com/cshep4/kripto/services/trader/internal/service SQS
//go:generate mockgen -destination=internal/mocks/coinbase/coinbase.gen.go -package=coinbase_mocks github.com/cshep4/kripto/services/trader/internal/trader Coinbase
