package run

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

type targetPort struct {
	target string
	port   *int
}

func Check(targets []string, interval time.Duration) {
	var targetPorts []targetPort
	checks := []check.Check{dns.New(), port.New()}

	if dnsHost, exists := os.LookupEnv("MANUAL_DNS_HOST"); exists {
		checks = append(checks, manualdns.New(dnsHost))
	}

	for _, t := range targets {
		hostPort := strings.Split(t, ":")

		target := targetPort{target: hostPort[0]}
		if len(hostPort) == 2 {
			p, err := strconv.Atoi(hostPort[1])
			if err == nil {
				target.port = &p
			}
		}

		if target.port != nil {
			log.Infof("Setup check for %s on port %d", target.target, *target.port)
		} else {
			log.Infof("Setup check for %s", target.target)
		}
		targetPorts = append(targetPorts, target)
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	execChan := make(chan check.Execution)
	go handleResults(ctx, execChan)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runChecks(ctx, execChan, targetPorts, checks)

		case <-sigChan:
			cancel()
			return
		}
	}
}

func handleResults(ctx context.Context, ex chan check.Execution) {
	for {
		select {
		case e := <-ex:
			e.Check.Report(e.Target, e.Port, e.Result)

		case <-ctx.Done():
			return
		}
	}
}

func runChecks(ctx context.Context, resultsChan chan check.Execution, targets []targetPort, checks []check.Check) {
	var wg sync.WaitGroup

	for _, target := range targets {
		for i := range checks {
			chk := checks[i]
			wg.Add(1)
			go func(host string) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(ctx, dnsCallTimeout)
				defer cancel()

				start := time.Now()

				executed, values, err := chk.Run(ctx, target.target, target.port)
				elapsed := time.Since(start)

				if executed {
					ex := check.NewExecution(chk, target.target, target.port)
					ex.Result.Values = values
					ex.Result.Duration = float64(elapsed.Nanoseconds()) / 1000000.
					ex.Result.Err = err
					ex.Result.TimedOut = err == context.Canceled
					resultsChan <- ex
				}
			}(target.target)
		}
	}
	wg.Wait()
}
