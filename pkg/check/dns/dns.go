package dns

import (
	"net"
	"strings"

	"go.ocp.six-group.net/opspoc/dns-checker/pkg/check"
)

// New create a new dns resolve check
func New(target string) check.Check {
	c := &dnsCheck{target: target}
	c.Setup("Host resolved to %s",
		"Error resolving host: %v",
		"dns_checker_check_dns",
		"Result of DNS check 0 = error, 1 = OK",
		"target")
	return c
}

type dnsCheck struct {
	check.BaseCheck
	target string
}

func (c *dnsCheck) Execute() ([]interface{}, error) {
	ips, err := net.LookupIP(c.target)
	return toString(ips), err
}

func (c *dnsCheck) Report(result []interface{}, err error, duration float64) {
	c.ReportResults(result, err, duration, c.target)
}

func toString(ips []net.IP) []interface{} {
	if ips == nil {
		return []interface{}{""}
	}
	var s []string
	for _, ip := range ips {
		s = append(s, ip.String())
	}
	return []interface{}{strings.Join(s, ",")}
}
