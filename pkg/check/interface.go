package check

import (
	"context"
)

// Check interface for checks
type Check interface {
	Run(ctx context.Context, address Address) *Result
	Report(result Result)
}

// Result check result
type Result struct {
	Values   []string
	Duration *float64
	Err      error
	TimedOut bool
}

// Address address with host and port
type Address struct {
	Host string
	Port *int
}
