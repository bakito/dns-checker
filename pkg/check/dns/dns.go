package dns

import (
	"context"
	"net"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
)

// New create a new dns resolve check
func New(interval time.Duration) check.Check {
	c := &dnsCheck{}
	c.Setup(interval,
		"Host resolved",
		"Error resolving host",
		"dns")
	return c
}

type dnsCheck struct {
	check.BaseCheck
}

func (c *dnsCheck) Run(ctx context.Context, address check.Address) *check.Result {
	_, err := net.DefaultResolver.LookupHost(ctx, address.Host)
	return &check.Result{Err: err}
}
