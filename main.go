package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	"github.com/bakito/dns-checker/pkg/check/dns"
	"github.com/bakito/dns-checker/pkg/check/manualdns"
	"github.com/bakito/dns-checker/pkg/check/port"
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
	log.Infof("Interval is %d seconds", interval)
	log.Infof("Starting on port %s", metricsPort)
	_ = http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), nil)
}

func recordMetrics() {
	var checks []check.Check
	for _, t := range targets {
		hostPort := strings.Split(t, ":")
		checks = append(checks, dns.New(hostPort[0]))
		if len(hostPort) == 2 {
			log.Infof("Checking %s on port %s", hostPort[0], hostPort[1])
			checks = append(checks, port.New(hostPort[0], hostPort[1]))
		} else if len(hostPort) == 1 {
			log.Infof("Checking %s", hostPort[0])
		}

		if dnsHost, exists := os.LookupEnv("MANUAL_DNS_HOST"); exists {
			checks = append(checks, manualdns.New(hostPort[0], dnsHost))
		}

	}

	go func() {
		for {
			log.Info("checking...")

			check.Execute(checks...)
			time.Sleep(interval * time.Second)
		}
	}()
}
