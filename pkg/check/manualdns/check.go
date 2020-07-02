package manualdns

import (
	"context"

	"github.com/bakito/dns-checker/pkg/check"
)

// New create a new dns resolve check
func New(dnsHost string) check.Check {
	c := &dnsCheck{}
	c.Setup("%s",
		"Error resolving host: %v",
		"dns_checker_check_manual_dns",
		"Result of DNS check 0 = error, 1 = OK",
		"target")
	c.dnsHost = dnsHost
	return c
}

type dnsCheck struct {
	check.BaseCheck
	dnsHost string
}

func (c *dnsCheck) query(target string) []byte {
	return dnsQuery{
		ID: 0xAAAA,
		RD: true,
		Questions: []dnsQuestion{{
			Domain: target,
			Type:   0x1, // A record
			Class:  0x1, // Internet
		}},
	}.encode()
}

func (c *dnsCheck) Run(ctx context.Context, target string, port *int) (bool, []string, error) {
	result, err := resolve(ctx, c.query(target), c.dnsHost)
	if err != nil {
		return true, []string{target}, err
	}

	_, err = responseCode(result)
	return true, []string{target}, err
}
