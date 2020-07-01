package check

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
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
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, labels)
	c.HistogramMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: fmt.Sprintf("%s_histogram", metricName),
		Help: fmt.Sprintf("The duration of resolver lookups %s in ms and percentiles", metricName),
	}, labels)
	c.name = metricName
	c.labels = labels
	c.MessageOK = ok
	c.MessageNOK = nok
}

// ReportResults report the check results
func (c *BaseCheck) Report(target string, port *int, result Result) {

	fields := log.Fields{}
	fields["name"] = c.name
	fields["duration"] = fmt.Sprintf("%dms", result.Duration)

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
