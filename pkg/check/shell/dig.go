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
	digCommand = "dig %s +noall +stats"
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
		"target")
	return c
}

type shellCheck struct {
	check.BaseCheck
}

func (c *shellCheck) Run(ctx context.Context, address check.Address) *check.Result {
	command := fmt.Sprintf(digCommand, address.Host)
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	out, err := cmd.Output()
	if err == nil {
		log.WithField("command", "dig").Debugf("%s\n", out)
	}

	res := &check.Result{Values: []string{address.Host}, Err: err}
	if queryTimePattern.MatchString(string(out)) {
		m := queryTimePattern.FindStringSubmatch(string(out))
		qt, _ := strconv.ParseFloat(m[1], 64)
		dur := time.Duration(qt) * time.Millisecond
		res.Duration = &dur
	} else {
		log.WithField("command", "dig").Debugf("error parsing query time")
	}

	return res
}
