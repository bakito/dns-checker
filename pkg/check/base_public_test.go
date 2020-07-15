package check_test

import (
	"testing"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func Test_Setup_Report(t *testing.T) {
	check.Init(time.Second)
	bc := check.BaseCheck{}
	bc.Setup("ok", "nok", "metricName")

	assert.Assert(t, is.Equal(bc.MessageOK, "ok"))
	assert.Assert(t, is.Equal(bc.MessageNOK, "nok"))
	duration := 1 * time.Second
	bc.Report(check.Address{}, check.Result{
		Duration: &duration,
	})
}
