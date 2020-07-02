package check

import (
	"os"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func Test_buckets(t *testing.T) {
	bucketTestData := []struct {
		env      string
		expected []float64
	}{
		{"",
			defaultBuckets,
		},
		{"foo",
			defaultBuckets,
		},
		{"0.05,0.1,0.25  ,  0.5,1,2.5",
			[]float64{0.05, 0.1, 0.25, 0.5, 1, 2.5},
		},
	}

	for _, data := range bucketTestData {
		if data.env == "" {
			os.Unsetenv(envMetricHistogramBuckets)
		} else {
			os.Setenv(envMetricHistogramBuckets, data.env)
		}

		currBuckets = nil
		b := buckets()

		assert.Assert(t, is.DeepEqual(b, data.expected))
		assert.Assert(t, is.DeepEqual(b, currBuckets))
	}
}

func Test_objectives(t *testing.T) {
	bucketTestData := []struct {
		env      string
		expected map[float64]float64
	}{
		{"",
			defaultObjectives,
		},
		{"foo",
			defaultObjectives,
		},
		{"0.5:0.05,0.9:0.01,  0.99 : 0.001  ",
			map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
	}

	for _, data := range bucketTestData {
		if data.env == "" {
			os.Unsetenv(envMetricSummaryObjectives)
		} else {
			os.Setenv(envMetricSummaryObjectives, data.env)
		}

		currObjectives = nil
		o := objectives()

		assert.Assert(t, is.DeepEqual(o, data.expected))
		assert.Assert(t, is.DeepEqual(o, currObjectives))
	}
}
