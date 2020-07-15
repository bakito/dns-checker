package check

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bakito/dns-checker/version"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

const (
	// envMetricSummaryObjectives a ',' separated list of ':' separated float64 objective key-value pairs.
	//  E.g: "0.5:0.05,0.9:0.01,0.99:0.001"
	envMetricSummaryObjectives = "METRICS_SUMMARY_OBJECTIVES"

	// envMetricHistogramBuckets a ',' separated list of float64 histogram buckets
	//  E.g: "0.002,0.005,0.01,0.025,0.05,0.1,0.25,0.5,1,2.5,5,10,20"
	envMetricHistogramBuckets = "METRICS_HISTOGRAM_BUCKETS"

	Separator = ","

	metricName          = "dns_checker_check"
	metricErrorName     = metricName + "_error"
	metricDurationName  = metricName + "_duration"
	metricSummaryName   = metricName + "_summary"
	metricHistogramName = metricName + "_histogram"
)

var (
	defaultObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
	defaultBuckets    = []float64{0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 20}

	currObjectives map[float64]float64
	currBuckets    []float64

	successMetric   *prometheus.GaugeVec
	errorMetric     *prometheus.GaugeVec
	durationMetric  *prometheus.GaugeVec
	summaryMetric   *prometheus.SummaryVec
	histogramMetric *prometheus.HistogramVec
)

func Init(timeout time.Duration) {
	labels := []string{"target", "port", "check_name", "version"}
	successMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName,
		Help: "Result of the check 0 = error, 1 = OK",
	}, labels)
	errorMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricErrorName,
		Help: "Check resulted in an error; 1 = error, 0 = OK",
	}, labels)
	durationMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricDurationName,
		Help: fmt.Sprintf("The duration of %s in ms", metricName),
	}, labels)
	summaryMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       metricSummaryName,
		Help:       "The duration of resolver lookups  in ms and percentiles",
		Objectives: objectives(),
	}, labels)
	histogramMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metricHistogramName,
		Help:    "The duration of resolver lookups in ms and buckets",
		Buckets: buckets(timeout),
	}, labels)
}

// BaseCheck basic check functionality
type BaseCheck struct {
	MessageOK  string
	MessageNOK string
	name       string
}

// Name get the name of the check
func (c *BaseCheck) Name() string {
	return c.name
}

// Setup setup the check
func (c *BaseCheck) Setup(ok string, nok string, name string) {
	c.name = name
	c.MessageOK = ok
	c.MessageNOK = nok

	log.WithField("name", metricName).Info("Setup check")
}

// Report report the check results
func (c *BaseCheck) Report(address Address, result Result) {

	duration := float64(*result.Duration) / float64(time.Millisecond)
	fields := log.Fields{}
	fields["name"] = c.name
	fields["duration"] = duration
	fields["worker"] = result.WorkerId
	fields["target"] = address.Host
	values := []string{address.Host}
	if address.Port != nil {
		fields["port"] = *address.Port
		values = append(values, fmt.Sprintf("%d", *address.Port))
	} else {
		values = append(values, "")
	}
	values = append(values, c.name, version.Version)

	l := log.WithFields(fields)
	if result.Err != nil {
		l.Warnf("%s : %v", c.MessageNOK, result.Err)
		successMetric.WithLabelValues(values...).Set(0)
		errorMetric.WithLabelValues(values...).Set(1)
	} else {
		l.Debug(c.MessageOK)
		successMetric.WithLabelValues(values...).Set(1)
		errorMetric.WithLabelValues(values...).Set(0)
	}
	durationMetric.WithLabelValues(values...).Set(duration)
	summaryMetric.WithLabelValues(values...).Observe(duration)
	histogramMetric.WithLabelValues(values...).Observe(duration)
}

func objectives() map[float64]float64 {
	if currObjectives != nil {
		return currObjectives
	}

	currObjectives = defaultObjectives

	if value, exists := os.LookupEnv(envMetricSummaryObjectives); exists {
		custom := make(map[float64]float64)
		objectives := strings.Split(value, Separator)
		for _, o := range objectives {
			objective := strings.Split(o, ":")
			if len(objective) == 2 {
				a, err := strconv.ParseFloat(strings.TrimSpace(objective[0]), 64)
				if err != nil {
					log.WithFields(log.Fields{"env": envMetricSummaryObjectives, "value": value, "default": defaultObjectives}).
						Warn("could not parse the objectives, using the default")
					return currObjectives
				}
				b, err := strconv.ParseFloat(strings.TrimSpace(objective[1]), 64)
				if err != nil {
					log.WithFields(log.Fields{"env": envMetricSummaryObjectives, "value": value, "default": defaultObjectives}).
						Warn("could not parse the objectives, using the default")
					return currObjectives
				}
				custom[a] = b
			} else {
				log.WithFields(log.Fields{"env": envMetricSummaryObjectives, "value": value, "default": defaultObjectives}).
					Warn("could not parse the objectives, using the default")
				return currObjectives
			}
		}
		currObjectives = custom
	}
	return currObjectives
}

func buckets(timeout time.Duration) []float64 {
	if currBuckets != nil {
		return currBuckets
	}
	currBuckets = filter(defaultBuckets, timeout)

	if value, exists := os.LookupEnv(envMetricHistogramBuckets); exists {
		var custom []float64
		objectives := strings.Split(value, Separator)
		for _, o := range objectives {
			a, err := strconv.ParseFloat(strings.TrimSpace(o), 64)
			if err != nil {
				log.WithFields(log.Fields{"env": envMetricHistogramBuckets, "value": value, "default": defaultBuckets}).
					Warn("could not parse the buckets, using the default")
				return currBuckets
			}
			custom = append(custom, a)
		}
		currBuckets = filter(custom, timeout)
		return currBuckets
	}
	return currBuckets
}

func filter(bucket []float64, timeout time.Duration) []float64 {
	var filtered []float64
	sec := timeout.Seconds()
	for _, b := range bucket {
		if b <= sec {
			filtered = append(filtered, b)
		}
	}
	return filtered
}
