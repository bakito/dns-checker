package port

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
)

// New create a new port probe check
func New(interval time.Duration) check.Check {
	c := &probeCheck{}
	c.Setup(interval,
		"Probe was successful",
		"Error probing",
		"dns_checker_probe_port",
		"target", "port")
	return c
}

type probeCheck struct {
	check.BaseCheck
}

func (c *probeCheck) Run(ctx context.Context, address check.Address) *check.Result {
	if address.Port == nil {
		return nil
	}
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%v:%v", address.Host, *address.Port))
	if conn != nil {
		_ = conn.Close()
	}
	return &check.Result{Values: []string{address.Host, fmt.Sprintf("%d", *address.Port)}, Err: err}
}
