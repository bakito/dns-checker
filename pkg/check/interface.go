package check

// Check interface for checks
type Check interface {
	Execute() ([]interface{}, error)
	Report(result []interface{}, err error, duration float64)
}
