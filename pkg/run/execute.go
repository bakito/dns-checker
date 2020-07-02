package run

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	"github.com/bakito/dns-checker/pkg/check/dns"
	"github.com/bakito/dns-checker/pkg/check/manualdns"
	"github.com/bakito/dns-checker/pkg/check/port"
	log "github.com/sirupsen/logrus"
)

const (
	dnsCallTimeout = 5 * time.Second
)

// Check run the checks
func Check(targets []check.Address, interval time.Duration) {
	checks := []check.Check{dns.New(), port.New()}

	if dnsHost, exists := os.LookupEnv("MANUAL_DNS_HOST"); exists {
		checks = append(checks, manualdns.New(dnsHost))
	}

	for _, target := range targets {

		if target.Port != nil {
			log.Infof("Setup check for %s on port %d", target.Host, *target.Port)
		} else {
			log.Infof("Setup check for %s", target.Host)
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	execChan := make(chan execution)
	go handleResults(ctx, execChan)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runChecks(ctx, execChan, targets, checks)

		case <-sigChan:
			cancel()
			return
		}
	}
}

func handleResults(ctx context.Context, ex chan execution) {
	for {
		select {
		case e := <-ex:
			e.check.Report(e.Result)

		case <-ctx.Done():
			return
		}
	}
}

func runChecks(ctx context.Context, resultsChan chan execution, targets []check.Address, checks []check.Check) {
	var wg sync.WaitGroup

	for _, t := range targets {
		for i := range checks {
			chk := checks[i]
			wg.Add(1)
			go func(target check.Address) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(ctx, dnsCallTimeout)
				defer cancel()

				start := time.Now()

				executed, values, err := chk.Run(ctx, target)
				elapsed := time.Since(start)

				if executed {
					ex := newExecution(chk, target)
					ex.Values = values
					ex.Duration = float64(elapsed.Nanoseconds()) / 1000000.
					ex.Err = err
					ex.TimedOut = err == context.Canceled
					resultsChan <- ex
				}
			}(t)
		}
	}
	wg.Wait()
}
