package check

type Check interface {
	Execute() ([]interface{}, error)
	Report(result []interface{}, err error, duration float64)
}
