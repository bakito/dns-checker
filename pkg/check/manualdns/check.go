package manualdns

import (
	"github.com/bakito/dns-checker/pkg/check"
)

// New create a new dns resolve check
func New(target, dnsHost string) check.Check {
	c := &dnsCheck{target: target}
	c.Setup("%s",
		"Error resolving host: %v",
		"dns_checker_check_manual_dns",
		"Result of DNS check 0 = error, 1 = OK",
		"target")
	c.dnsHost = dnsHost

	c.query = dnsQuery{
		ID: 0xAAAA,
		RD: true,
		Questions: []dnsQuestion{{
			Domain: target,
			Type:   0x1, // A record
			Class:  0x1, // Internet
		}},
	}.encode()

	return c
}

type dnsCheck struct {
	check.BaseCheck
	target  string
	dnsHost string
	query   []byte
}

func (c *dnsCheck) Execute() ([]interface{}, error) {
	result, err := resolve(c.query, c.dnsHost)
	return c.ToResult(result), err
}

func (c *dnsCheck) Report(result []interface{}, err error, duration float64) {
	r, e := responseCode(result[0].(byte))
	c.ReportResults(c.ToResult(r), e, duration, c.target)
}
