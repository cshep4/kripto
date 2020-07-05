package idempotency

//go:generate mockgen -destination=internal/mocks/idempotency/idempotencer.gen.go -package=idempotency_mocks github.com/cshep4/kripto/shared/go/idempotency Idempotencer
