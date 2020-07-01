package check

import (
	"context"
)

// Check interface for checks
type Check interface {
	Run(ctx context.Context, target string, port *int) (bool, []string, error)
	Report(target string, port *int, result Result)
}

func NewExecution(check Check, target string, port *int) Execution {
	return Execution{
		Target: target,
		Port:   port,
		Check:  check,
		Result: Result{},
	}
}

type Execution struct {
	Check  Check
	Target string
	Port   *int
	Result Result
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
