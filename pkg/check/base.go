package check

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

	separator = ","
)

var (
	defaultObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
	defaultBuckets    = []float64{0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 20}

	currObjectives map[float64]float64
	currBuckets    []float64
)

// BaseCheck basic check functionality
type BaseCheck struct {
	MessageOK       string
	MessageNOK      string
	StateMetric     *prometheus.GaugeVec
	DurationMetric  *prometheus.GaugeVec
	SummaryMetric   *prometheus.SummaryVec
	HistogramMetric *prometheus.HistogramVec
	name            string
	labels          []string
}

// Setup setup the check
func (c *BaseCheck) Setup(ok string, nok string, metricName string, metricHelp string, labels ...string) {
	c.StateMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName,
		Help: metricHelp,
	}, labels)
	c.DurationMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_duration", metricName),
		Help: fmt.Sprintf("The duration of %s in ms", metricName),
	}, labels)
	c.SummaryMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       fmt.Sprintf("%s_summary", metricName),
		Help:       fmt.Sprintf("The duration of resolver lookups %s in ms and percentiles", metricName),
		Objectives: objectives(),
	}, labels)
	c.HistogramMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    fmt.Sprintf("%s_histogram", metricName),
		Help:    fmt.Sprintf("The duration of resolver lookups %s in ms and percentiles", metricName),
		Buckets: buckets(),
	}, labels)
	c.name = metricName
	c.labels = labels
	c.MessageOK = ok
	c.MessageNOK = nok
}

// Report report the check results
func (c *BaseCheck) Report(target string, port *int, result Result) {

	fields := log.Fields{}
	fields["name"] = c.name
	fields["duration"] = fmt.Sprintf("%vms", result.Duration)

	for i, v := range result.Values {
		fields[c.labels[i]] = v
	}

	l := log.WithFields(fields)
	if result.Err != nil {
		l.Warnf(c.MessageNOK, result.Err)
		c.StateMetric.WithLabelValues(result.Values...).Set(0)
	} else {
		l.Debugf(c.MessageOK, result.valuesAsInterface()...)
		c.StateMetric.WithLabelValues(result.Values...).Set(1)
	}
	c.DurationMetric.WithLabelValues(result.Values...).Set(result.Duration)
	c.SummaryMetric.WithLabelValues(result.Values...).Observe(result.Duration)
	c.HistogramMetric.WithLabelValues(result.Values...).Observe(result.Duration)
}

// ToResult maps to interface array
func (c *BaseCheck) ToResult(values ...interface{}) []interface{} {
	return values
}

func objectives() map[float64]float64 {
	if currObjectives != nil {
		return currObjectives
	}

	currObjectives = defaultObjectives

	if value, exists := os.LookupEnv(envMetricSummaryObjectives); exists {
		custom := make(map[float64]float64)
		objectives := strings.Split(value, separator)
		for _, o := range objectives {
			objective := strings.Split(o, ":")
			if len(objective) == 2 {
				a, err := strconv.ParseFloat(objective[0], 64)
				if err != nil {
					log.WithFields(log.Fields{"env": envMetricSummaryObjectives, "value": value, "default": defaultObjectives}).
						Warn("could not parse the objectives, using the default")
					return currObjectives
				}
				b, err := strconv.ParseFloat(objective[1], 64)
				if err != nil {
					log.WithFields(log.Fields{"env": envMetricSummaryObjectives, "value": value, "default": defaultObjectives}).
						Warn("could not parse the objectives, using the default")
					return currObjectives
				}
				custom[a] = b
			}
		}
		currObjectives = custom
	}
	return currObjectives
}

func buckets() []float64 {
	if currBuckets != nil {
		return currBuckets
	}
	currBuckets = defaultBuckets

	if value, exists := os.LookupEnv(envMetricHistogramBuckets); exists {
		var custom []float64
		objectives := strings.Split(value, separator)
		for _, o := range objectives {
			a, err := strconv.ParseFloat(o, 64)
			if err != nil {
				log.WithFields(log.Fields{"env": envMetricHistogramBuckets, "value": value, "default": defaultBuckets}).
					Warn("could not parse the buckets, using the default")
				return currBuckets
			}
			custom = append(custom, a)
		}
		currBuckets = custom
		return currBuckets
	}
	return currBuckets
}
