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

func Test_toTargets(t *testing.T) {
	targets, err := toTargets([]string{"a", "b:1234","c,d:5678 , e:9999     "})
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, is.Len(targets, 5))
}

func Test_toTarget(t *testing.T) {
	target, err := toTarget("host.name")
	assert.Assert(t, is.Equal(target.Host, "host.name"))
	assert.Assert(t, is.Nil(target.Port))
	assert.Assert(t, is.Nil(err))

	target, err = toTarget("host.name:1234")
	assert.Assert(t, is.Equal(target.Host, "host.name"))
	assert.Assert(t, target.Port != nil)
	assert.Assert(t, is.Equal(*target.Port, 1234))
	assert.Assert(t, is.Nil(err))

	_, err = toTarget("host.name:not-a-port")
	assert.Assert(t, is.Error(err, `port "not-a-port" of host "host.name" can not be parsed as int`))
}

func Test_fromEnv(t *testing.T) {
	os.Setenv(testEnv, "foo")

	assert.Assert(t, is.Equal(fromEnv("bar"), "bar"))
	assert.Assert(t, is.Equal(fromEnv("${"+testEnv+"}"), "foo"))
}
