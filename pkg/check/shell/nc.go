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
	// NameNC the name of the nc check
	NameNC = "nc"
)

var (
	receivedPattern = regexp.MustCompile(`.*bytes received in (\d+[.]\d*) seconds.*`)
)

// NewNc create a new nc command check
func NewNc() check.Check {
	c := &ncCheck{}
	c.Setup(
		"Netcat succeeded",
		"Error executing nc",
		NameNC)
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
	res := &check.Result{Err: err}
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
