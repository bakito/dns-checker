package main

import (
	"fmt"
	"github.com/bakito/dns-checker/version"
	"net/http"
	"os"
	"strconv"
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
	envLogJSON     = "LOG_JSON"
	envInterval    = "INTERVAL"
	envWorker      = "WORKER"
)

var (
	logLevel    = log.InfoLevel
	metricsPort = "2112"
	interval    = 30 * time.Second
	worker      = 10
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
	if json, exists := os.LookupEnv(envLogJSON); exists {
		if enabled, err := strconv.ParseBool(json); enabled && err == nil {
			log.SetFormatter(&log.JSONFormatter{})
		}
	}

	if p, exists := os.LookupEnv(envMetricsPort); exists {
		metricsPort = p
	}
	if i, exists := os.LookupEnv(envInterval); exists {
		interval, err = time.ParseDuration(i)
		if err != nil {
			panic(fmt.Errorf("env var %s %q can not be parsed as duration", envInterval, i))
		}
	}

	if w, exists := os.LookupEnv(envWorker); exists {
		worker, err = strconv.Atoi(w)
		if err != nil {
			panic(fmt.Errorf("env var %s %q can not be parsed as int", envWorker, w))
		}
	}
	log.WithFields(log.Fields{
		"interval": fmt.Sprintf("%v", interval),
		"workers":  worker,
		"version":  version.Version}).
		Info("Interval")
}

func main() {
	go serveMetrics()

	values := findTargets()
	if len(values) == 0 {
		panic(fmt.Errorf("env var %s is needed", envTarget))
	}
	err := run.Check(values, interval, worker)
	if err != nil {
		panic(err)
	}
}

func serveMetrics() {
	log.WithField("port", metricsPort).Info("Starting metrics")
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
