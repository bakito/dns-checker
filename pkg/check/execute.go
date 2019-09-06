package check

import "time"

// Execute execute the given checks
func Execute(checks ...Check) {
	for _, c := range checks {
		start := time.Now()
		result, err := c.Execute()
		elapsed := time.Since(start)
		dur := float64(elapsed.Nanoseconds()) / 1000000.

		c.Report(result, err, dur)
	}
}
