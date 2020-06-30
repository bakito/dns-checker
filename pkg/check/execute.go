package check

import (
	"context"
	"time"
)

// Execute execute the given checks
func Execute(checks ...Check) {
	for _, check := range checks {
		go execute(check)
	}
}

func execute(check Check) {
	start := time.Now()
	ctx := context.Background()
	result, err := check.Execute(ctx)
	elapsed := time.Since(start)
	dur := float64(elapsed.Nanoseconds()) / 1000000.

	check.Report(result, err, dur)
}
