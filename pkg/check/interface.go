package check

import (
	"context"
	"time"
)

// Check interface for checks
type Check interface {
	Run(ctx context.Context, address Address) *Result
	Report(address Address, result Result)
	Name() string
}

// Result check result
type Result struct {
	Duration *time.Duration
	Err      error
	TimedOut bool
	WorkerID int
}

// Address address with host and port
type Address struct {
	Host string
	Port *int
}
