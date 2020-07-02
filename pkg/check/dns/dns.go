package dns

import (
	"context"
	"net"

	"github.com/bakito/dns-checker/pkg/check"
)

// New create a new dns resolve check
func New() check.Check {
	c := &dnsCheck{}
	c.Setup("Host resolved",
		"Error resolving host",
		"dns_checker_check_dns",
		"Result of DNS check 0 = error, 1 = OK",
		"target")
	return c
}

type dnsCheck struct {
	check.BaseCheck
}

func (c *dnsCheck) Run(ctx context.Context, address check.Address) (bool, []string, error) {
	_, err := net.DefaultResolver.LookupHost(ctx, address.Host)
	return true, []string{address.Host}, err
}
