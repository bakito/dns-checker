package shell

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	log "github.com/sirupsen/logrus"
)

const (
	ncCommand = "nc -zv  %s %d"
)

var (
	receivedPattern = regexp.MustCompile(`.*bytes received in (\d+[.]\d*) seconds.*`)
)

// New create a new nc command check
func NewNc(interval time.Duration) check.Check {
	c := &ncCheck{}
	c.Setup(interval,
		"Netcat succeeded",
		"Error executing nc",
		"dns_checker_check_nc",
		"target", "port")
	return c
}

type ncCheck struct {
	check.BaseCheck
}

func (c *ncCheck) Run(ctx context.Context, address check.Address) *check.Result {
	if address.Port == nil {
		return nil
	}
	command := fmt.Sprintf(ncCommand, address.Host, *address.Port)
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	out, err := cmd.CombinedOutput()
	res := &check.Result{Values: []string{address.Host, fmt.Sprintf("%d", *address.Port)}, Err: err}
	if err != nil {
		return res
	}
	log.WithField("command", "nc").Debugf("%s\n", out)

	if receivedPattern.MatchString(string(out)) {
		m := receivedPattern.FindStringSubmatch(string(out))
		dur, _ := time.ParseDuration(fmt.Sprintf("%ss", m[1]))
		res.Duration = &dur
	} else {
		log.WithField("command", "nc").Debugf("error parsing query time")
	}
	return res
}
