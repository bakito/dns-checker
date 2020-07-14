package shell

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	log "github.com/sirupsen/logrus"
)

const (
	digCommand = "dig %s +noall +stats 1>&2 stderr"
)

var (
	queryTimePattern = regexp.MustCompile(`.*;; Query time: (\d+) msec.*`)
)

// New create a new dig command check
func NewDig(interval time.Duration) check.Check {
	c := &shellCheck{}
	c.Setup(interval,
		"Dig succeeded",
		"Error executing dig",
		"dns_checker_check_dig",
		"Result of dig command check 0 = error, 1 = OK",
		"target")
	return c
}

type shellCheck struct {
	check.BaseCheck
}

func (c *shellCheck) Run(ctx context.Context, address check.Address) *check.Result {
	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf(digCommand, address.Host))
	stdoutStderr, err := cmd.CombinedOutput()
	if err == nil {
		log.WithField("command", "dig").Debugf("%s\n", stdoutStderr)
	}

	res := &check.Result{Values: []string{address.Host}, Err: err}
	if queryTimePattern.MatchString(string(stdoutStderr)) {
		m := queryTimePattern.FindStringSubmatch(string(stdoutStderr))
		dur, _ := strconv.ParseFloat(m[1], 64)
		res.Duration = dur
	} else {
		log.WithField("command", "dig").Debugf("error parsing query time")
	}

	return res
}
