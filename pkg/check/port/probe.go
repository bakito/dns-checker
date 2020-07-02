package port

import (
	"context"
	"fmt"
	"net"

	"github.com/bakito/dns-checker/pkg/check"
)

// New create a new port probe check
func New() check.Check {
	c := &probeCheck{}
	c.Setup("Probe was successful",
		"Error probing",
		"dns_checker_probe_port",
		"Result of port probe 0 = error, 1 = OK",
		"target", "port")
	return c
}

type probeCheck struct {
	check.BaseCheck
}

func (c *probeCheck) Run(ctx context.Context, address check.Address) (bool, []string, error) {
	if address.Port == nil {
		return false, nil, nil
	}
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%v:%v", address.Host, *address.Port))
	if conn != nil {
		_ = conn.Close()
	}
	return true, []string{address.Host, fmt.Sprintf("%d", *address.Port)}, err
}
