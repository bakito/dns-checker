package port

import (
	"context"
	"fmt"
	"net"

	"github.com/bakito/dns-checker/pkg/check"
)

const (
	// Name the name of this check
	Name = "probe-port"
)

// New create a new port probe check
func New() check.Check {
	c := &probeCheck{}
	c.Setup(
		"Probe was successful",
		"Error probing",
		Name)
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
	return &check.Result{Err: err}
}
