package check

import (
	"context"
	"time"
)

// Check interface for checks
type Check interface {
	Run(ctx context.Context, address Address) *Result
	Report(result Result)
	Name() string
}

// Result check result
type Result struct {
	Values   []string
	Duration *time.Duration
	Err      error
	TimedOut bool
	WorkerId int
}

// Address address with host and port
type Address struct {
	Host string
	Port *int
}
