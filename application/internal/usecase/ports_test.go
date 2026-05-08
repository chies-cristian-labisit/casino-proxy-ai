package usecase

// Compile-time assertions: both interfaces are declared in this package.
var (
	_ CustomerRepository
	_ IdempotencyStore
)
