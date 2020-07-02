package run

import (
	"os"
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

const (
	testEnv = "___TEST___"
)

func Test_toTarget(t *testing.T) {
	target := toTarget("host.name")
	assert.Assert(t, is.Equal(target.Host, "host.name"))
	assert.Assert(t, is.Nil(target.Port))

	target = toTarget("host.name:1234")
	assert.Assert(t, is.Equal(target.Host, "host.name"))
	assert.Assert(t, is.Equal(*target.Port, 1234))
}

func Test_fromEnv(t *testing.T) {
	os.Setenv(testEnv, "foo")

	assert.Assert(t, is.Equal(fromEnv("bar"), "bar"))
	assert.Assert(t, is.Equal(fromEnv("${"+testEnv+"}"), "foo"))
}
