package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.ocp.six-group.net/opspoc/dns-checker/pkg/check"

	"go.ocp.six-group.net/opspoc/dns-checker/pkg/check/dns"
	"go.ocp.six-group.net/opspoc/dns-checker/pkg/check/port"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	logLevel    = log.InfoLevel
	metricsPort = "2112"
	targetPort  string
	target      string
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
	if p, exists := os.LookupEnv("METRICS_PORT"); exists {
		metricsPort = p
	}
	if t, exists := os.LookupEnv("TARGET"); exists {
		target = t
	} else {
		panic(fmt.Errorf("env var TARGET is needed"))
	}
	if tp, exists := os.LookupEnv("TARGET_PORT"); exists {
		targetPort = tp
	} else {
		panic(fmt.Errorf("env var TARGET_PORT is needed"))
	}
	if i, exists := os.LookupEnv("INTERVAL"); exists {
		ii, err := strconv.Atoi(i)
		if err != nil {
			panic(fmt.Errorf("env var TARGET_PORT is needed"))
		}
		interval = time.Duration(ii)
	}
	log.SetLevel(logLevel)
}

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Checking %s on port %s", target, targetPort)
	log.Infof("Interval is %d seconds", interval)
	log.Infof("Starting on port %s", metricsPort)
	_ = http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), nil)
}

func recordMetrics() {
	checks := []check.Check{dns.New(target), port.New(target, targetPort)}
	go func() {
		for {
			log.Info("checking...")

			check.Execute(checks...)
			time.Sleep(interval * time.Second)
		}
	}()
}
