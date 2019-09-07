package check

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

// BaseCheck basic check functionality
type BaseCheck struct {
	MessageOK      string
	MessageNOK     string
	StateMetric    *prometheus.GaugeVec
	DurationMetric *prometheus.GaugeVec
	name           string
	labels         []string
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
	c.name = metricName
	c.labels = labels
	c.MessageOK = ok
	c.MessageNOK = nok
}

// ReportResults report the check results
func (c *BaseCheck) ReportResults(result []interface{}, err error, duration float64, values ...string) {

	fields := log.Fields{}
	fields["name"] = c.name
	fields["duration"] = fmt.Sprintf("%vms", duration)
	for i, v := range values {
		fields[c.labels[i]] = v
	}

	l := log.WithFields(fields)
	if err != nil {
		l.Warnf(c.MessageNOK, err)
		c.StateMetric.WithLabelValues(values...).Set(0)
	} else {
		l.Debugf(c.MessageOK, result...)
		c.StateMetric.WithLabelValues(values...).Set(1)
	}
	c.DurationMetric.WithLabelValues(values...).Set(duration)
}

// ToResult maps to interface array
func (c *BaseCheck) ToResult(values ...interface{}) []interface{} {
	return values
}
