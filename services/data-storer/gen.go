package datastorer

//go:generate mockgen -destination=internal/mocks/service/service.gen.go -package=service_mocks github.com/cshep4/kripto/services/data-storer/internal/handler/aws Servicer
//go:generate mockgen -destination=internal/mocks/trade/store.gen.go -package=trade_mocks github.com/cshep4/kripto/services/data-storer/internal/service TradeStore
//go:generate mockgen -destination=internal/mocks/rate/store.gen.go -package=rate_mocks github.com/cshep4/kripto/services/data-storer/internal/service RateStore
