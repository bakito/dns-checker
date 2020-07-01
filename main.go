package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bakito/dns-checker/pkg/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	logLevel    = log.InfoLevel
	metricsPort = "2112"
	targets     []string
	interval    time.Duration = 30
)

func init() {
	var err error
	if ll, exists := os.LookupEnv("LOG_LEVEL"); exists {
		logLevel, err = log.ParseLevel(ll)
		if err != nil {
			panic(fmt.Errorf("error parsing log level"))
		}
	}
	log.SetLevel(logLevel)

	if p, exists := os.LookupEnv("METRICS_PORT"); exists {
		metricsPort = p
	}
	if t, exists := os.LookupEnv("TARGET"); exists {
		inputTargets := strings.Split(t, ";")
		for _, t := range inputTargets {
			targets = append(targets, strings.TrimSpace(t))
		}
	} else {
		panic(fmt.Errorf("env var TARGET is needed"))
	}
	if tp, exists := os.LookupEnv("TARGET_PORT"); exists {
		if len(targets) == 1 {
			old := targets[0]
			if strings.Contains(old, ":") {
				targets[0] = fmt.Sprintf("%s:%s", old, tp)
			}
		}

	}
	if i, exists := os.LookupEnv("INTERVAL"); exists {
		ii, err := time.ParseDuration(i)
		if err != nil {
			panic(fmt.Errorf("env var INTERVAL can not be parsed as duration"))
		}
		interval = time.Duration(ii)
	}
	log.Infof("Interval is %v", interval)
}

func main() {
	go serveMetrics()
	run.Check(targets, interval)
}

func serveMetrics() {
	log.Infof("Starting metrics on port %s", metricsPort)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), nil))
}
