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
		"Error probing: %v",
		"dns_checker_probe_port",
		"Result of port probe 0 = error, 1 = OK",
		"target", "port")
	return c
}

type probeCheck struct {
	check.BaseCheck
}

func (c *probeCheck) Run(ctx context.Context, target string, port *int) (bool, []string, error) {
	if port == nil {
		return false, nil, nil
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", target, *port))
	if conn != nil {
		_ = conn.Close()
	}
	return true, []string{target, fmt.Sprintf("%d", *port)}, err
}
