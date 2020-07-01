package check

import (
	"context"
)

// Check interface for checks
type Check interface {
	Run(ctx context.Context, target string, port *int) (bool, []string, error)
	Report(target string, port *int, result Result)
}

type Result struct {
	Values   []string
	Duration float64
	Err      error
	TimedOut bool
}

func (r *Result) valuesAsInterface() []interface{} {
	var out []interface{}
	for _, r := range r.Values {
		out = append(out, r)
	}
	return out
}
