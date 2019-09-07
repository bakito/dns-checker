package check

import "time"

// Execute execute the given checks
func Execute(checks ...Check) {
	for _, check := range checks {
		go execute(check)
	}
}

func execute(check Check) {
	start := time.Now()
	result, err := check.Execute()
	elapsed := time.Since(start)
	dur := float64(elapsed.Nanoseconds()) / 1000000.

	check.Report(result, err, dur)
}
