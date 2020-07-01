package dns

import (
	"context"
	"net"

	"github.com/bakito/dns-checker/pkg/check"
)

// New create a new dns resolve check
func New() check.Check {
	c := &dnsCheck{}
	c.Setup("Host resolved to %s",
		"Error resolving host: %v",
		"dns_checker_check_dns",
		"Result of DNS check 0 = error, 1 = OK",
		"target")
	return c
}

type dnsCheck struct {
	check.BaseCheck
}

func (c *dnsCheck) Run(ctx context.Context, target string, port *int) (bool, []string, error) {
	_, err := net.DefaultResolver.LookupHost(ctx, target)
	return true, []string{target}, err
}
