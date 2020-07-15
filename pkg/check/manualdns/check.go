package manualdns

import (
	"context"
	"fmt"

	"github.com/bakito/dns-checker/pkg/check"
)

const (
	Name = "manual_dns"
)

// New create a new dns resolve check
func New(dnsHost string) check.Check {
	c := &dnsCheck{}
	c.Setup(
		fmt.Sprintf("Host resolved with dns server %s", dnsHost),
		fmt.Sprintf("Error resolving host with dns server %s", dnsHost),
		Name)
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

func (c *dnsCheck) Run(ctx context.Context, address check.Address) *check.Result {
	result, err := resolve(ctx, c.query(address.Host), c.dnsHost)
	if err != nil {
		return &check.Result{Err: err}
	}

	_, err = responseCode(result)
	return &check.Result{Err: err}
}
