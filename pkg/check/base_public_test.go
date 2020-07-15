package check_test

import (
	"testing"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func Test_Setup_Report(t *testing.T) {
	bc := check.BaseCheck{}
	bc.Setup(5*time.Second, "ok", "nok", "metricName", "labels1", "label2")

	assert.Assert(t, is.Equal(bc.MessageOK, "ok"))
	assert.Assert(t, is.Equal(bc.MessageNOK, "nok"))
	duration := 1 * time.Second
	bc.Report(check.Result{
		Duration: &duration,
		Values:   []string{"a", "b"},
	})
}
