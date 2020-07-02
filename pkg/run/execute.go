package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
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

var (
	targetEnvVarPattern = regexp.MustCompile(`^\${(.*)}$`)
)

// Check run the checks
func Check(values []string, interval time.Duration) error {
	targetsAddresses, err := toTargets(values)
	if err != nil {
		return err
	}

	checks := []check.Check{dns.New(), port.New()}

	if dnsHost, exists := os.LookupEnv("MANUAL_DNS_HOST"); exists {
		checks = append(checks, manualdns.New(dnsHost))
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
			runChecks(ctx, execChan, targetsAddresses, checks)

		case <-sigChan:
			cancel()
			return nil
		}
	}
}

func toTargets(values []string) ([]check.Address, error) {
	var targetsAddresses []check.Address
	for _, value := range values {
		targets := strings.Split(value, ",")
		for _, t := range targets {
			target, err := toTarget(t)
			if err != nil {
				return nil, err
			}
			if target.Port != nil {
				log.Infof("Setup check for %s on port %d", target.Host, *target.Port)
			} else {
				log.Infof("Setup check for %s", target.Host)
			}
			targetsAddresses = append(targetsAddresses, target)
		}
	}
	return targetsAddresses, nil
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

func toTarget(in string) (check.Address, error) {
	hp := strings.Split(strings.TrimSpace(in), ":")

	host := fromEnv(strings.TrimSpace(hp[0]))

	addr := check.Address{Host: host}
	if len(hp) == 1 {
		return addr, nil
	}

	port := fromEnv(strings.TrimSpace(hp[1]))

	p, err := strconv.Atoi(port)
	if err != nil {
		return addr, fmt.Errorf("port %q of host %q can not be parsed as int", port, host)
	}
	addr.Port = &p
	return addr, nil
}

func fromEnv(in string) string {
	if targetEnvVarPattern.MatchString(in) {
		match := targetEnvVarPattern.FindStringSubmatch(in)
		return os.Getenv(match[1])
	}
	return in
}
