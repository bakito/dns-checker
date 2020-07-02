package check_test

import (
	"testing"

	"github.com/bakito/dns-checker/pkg/check"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func Test_Setup_Report(t *testing.T) {
	bc := check.BaseCheck{}
	bc.Setup("ok", "nok", "metricName", "metricHelp", "labels1", "label2")

	assert.Assert(t, is.Equal(bc.MessageOK, "ok"))
	assert.Assert(t, is.Equal(bc.MessageNOK, "nok"))
	bc.Report(check.Result{
		Values: []string{"a", "b"},
	})
}
