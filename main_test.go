package main

import (
	"os"
	"sort"
	"strings"
	"testing"

	. "gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func Test_findTargets(t *testing.T) {
	// prepare / cleanup
	for _, e := range os.Environ() {
		variable := strings.Split(e, "=")
		if strings.HasSuffix(variable[0], envTarget) {
			_ = os.Unsetenv(variable[0])
		}
	}

	// setup
	_ = os.Setenv(envTarget+"_Test1", "a")
	_ = os.Setenv(envTarget+"_Test2", "b")

	targets := findTargets()
	Assert(t, is.Len(targets, 2))
	sort.Strings(targets)
	Assert(t, is.Equal(targets[0], "a"))
	Assert(t, is.Equal(targets[1], "b"))

}
