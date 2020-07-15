package dns

import (
	"context"
	"net"

	"github.com/bakito/dns-checker/pkg/check"
)

const (
	// Name the name of this check
	Name = "dns"
)

// New create a new dns resolve check
func New() check.Check {
	c := &dnsCheck{}
	c.Setup(
		"Host resolved",
		"Error resolving host",
		Name)
	return c
}

type dnsCheck struct {
	check.BaseCheck
}

func (c *dnsCheck) Run(ctx context.Context, address check.Address) *check.Result {
	_, err := net.DefaultResolver.LookupHost(ctx, address.Host)
	return &check.Result{Err: err}
}
