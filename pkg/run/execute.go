package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	"github.com/bakito/dns-checker/pkg/check/dns"
	"github.com/bakito/dns-checker/pkg/check/manualdns"
	"github.com/bakito/dns-checker/pkg/check/port"
	"github.com/bakito/dns-checker/pkg/check/shell"
	log "github.com/sirupsen/logrus"
)

const (
	envManualDnsHost  = "MANUAL_DNS_HOST"
	envRunShellChecks = "RUN_SHELL_CHECKS"
	envLogDuration    = "LOG_DURATION"
)

var (
	targetEnvVarPattern = regexp.MustCompile(`^\${(.*)}$`)
)

// Check run the checks
func Check(values []string, interval time.Duration, worker int) error {
	targetsAddresses, err := toTargets(values)
	if err != nil {
		return err
	}

	checks := []check.Check{dns.New(interval), port.New(interval)}

	if dnsHost, exists := os.LookupEnv(envManualDnsHost); exists {
		checks = append(checks, manualdns.New(dnsHost, interval))
	}

	if boolEnv(envRunShellChecks) {
		checks = append(checks, shell.NewDig(interval))
		checks = append(checks, shell.NewNc(interval))
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	execChan := make(chan execution)
	go handleResults(ctx, execChan)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	collector := startDispatcher(worker) // start up worker pool

	for {
		select {
		case <-ticker.C:

			for _, t := range targetsAddresses {
				for chk := range checks {
					collector.work <- work{ctx, interval, execChan, t, checks[chk]}
				}
			}

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
			l := log.WithField("host", target.Host)
			if target.Port != nil {
				l = l.WithField("port", *target.Port)
			}

			l.Info("Setup check")

			targetsAddresses = append(targetsAddresses, target)
		}
	}
	return targetsAddresses, nil
}

func handleResults(ctx context.Context, ex chan execution) {
	for {
		select {
		case e := <-ex:
			e.check.Report(e.Address, e.Result)

		case <-ctx.Done():
			return
		}
	}
}

func runCheck(w work, workerID int) {
	ctx, cancel := context.WithTimeout(w.ctx, w.interval)
	defer cancel()

	start := time.Now()
	result := w.chk.Run(ctx, w.target)
	result.WorkerId = workerID
	duration := time.Since(start)

	if log.GetLevel() > log.InfoLevel || boolEnv(envLogDuration) {
		logDuration(w.chk, w.target, result, duration)
	}
	if result != nil {
		ex := newExecution(w.chk, w.target)
		if result.Duration == nil {
			ex.Duration = &duration
		} else {
			ex.Duration = result.Duration
		}
		ex.Err = result.Err
		ex.TimedOut = result.Err == context.Canceled
		w.resultsChan <- ex
	}
}

func logDuration(chk check.Check, target check.Address, result *check.Result, duration time.Duration) {
	l := log.WithFields(log.Fields{
		"name":     chk.Name(),
		"host":     target.Host,
		"worker":   result.WorkerId,
		"duration": float64(duration) / float64(time.Millisecond)})
	if result.Duration != nil {
		l = l.WithField("check-duration", float64(*result.Duration)/float64(time.Millisecond))
	}
	if target.Port != nil {
		l = l.WithField("port", *target.Port)
	}
	l.Info("check executed")
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

func boolEnv(name string) bool {
	if val, exists := os.LookupEnv(name); exists {
		run, _ := strconv.ParseBool(val)
		return run
	}
	return false
}
