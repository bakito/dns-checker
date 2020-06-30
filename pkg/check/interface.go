package check

import "context"

// Check interface for checks
type Check interface {
	Execute(ctx context.Context) ([]interface{}, error)
	Report(result []interface{}, err error, duration float64)
}
