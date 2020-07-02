package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bakito/dns-checker/pkg/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	envTarget      = "TARGET"
	envMetricsPort = "METRICS_PORT"
	envLogLevel    = "LOG_LEVEL"
	envInterval    = "INTERVAL"
)

var (
	logLevel    = log.InfoLevel
	metricsPort = "2112"
	interval    = 30 * time.Second

	targetEnvVarPattern = regexp.MustCompile(`^\${(.*)}$`)
)

func init() {
	var err error
	if ll, exists := os.LookupEnv(envLogLevel); exists {
		logLevel, err = log.ParseLevel(ll)
		if err != nil {
			panic(fmt.Errorf("error parsing log level"))
		}
	}
	log.SetLevel(logLevel)

	if p, exists := os.LookupEnv(envMetricsPort); exists {
		metricsPort = p
	}
	if i, exists := os.LookupEnv(envInterval); exists {
		interval, err = time.ParseDuration(i)
		if err != nil {
			panic(fmt.Errorf("env var %s %q can not be parsed as duration", envInterval, i))
		}
	}
	log.Infof("Interval is %v", interval)
}

func main() {
	go serveMetrics()

	values := findTargets()
	if len(values) == 0 {
		panic(fmt.Errorf("env var %s is needed", envTarget))
	}
	run.Check(values, interval)
}

func serveMetrics() {
	log.Infof("Starting metrics on port %s", metricsPort)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), nil))
}

func findTargets() []string {
	var targets []string
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], envTarget) {
			targets = append(targets, pair[1])
		}
	}
	return targets
}
