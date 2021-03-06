package shell

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	log "github.com/sirupsen/logrus"
)

const (
	digCommand = "dig %s"
	// NameDig the name of the dig check
	NameDig = "dig"
)

var (
	queryTimePattern = regexp.MustCompile(`.*;; Query time: (\d+) msec.*`)
	noErrorPattern   = regexp.MustCompile(`.*status: NOERROR.*`)
)

// NewDig create a new dig command check
func NewDig() check.Check {
	c := &digCheck{}
	c.Setup(
		"Dig succeeded",
		"Error executing dig",
		NameDig)
	return c
}

type digCheck struct {
	check.BaseCheck
}

func (c *digCheck) Run(ctx context.Context, address check.Address) *check.Result {
	command := fmt.Sprintf(digCommand, address.Host)
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	out, err := cmd.Output()
	res := &check.Result{Err: err}
	if err != nil {
		return res
	}

	// fail if 'NOERROR' is not found
	if !noErrorPattern.Match(out) {
		res.Err = errors.New(string(out))
		return res
	}

	log.WithField("command", "dig").Debugf("%s\n", out)

	if queryTimePattern.Match(out) {
		m := queryTimePattern.FindStringSubmatch(string(out))
		dur, _ := time.ParseDuration(fmt.Sprintf("%sms", m[1]))
		res.Duration = &dur
	} else {
		log.WithField("command", "dig").Debugf("error parsing query time")
	}

	return res
}
