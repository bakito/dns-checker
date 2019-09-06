package port

import (
	"fmt"
	"net"

	"github.com/bakito/dns-checker/pkg/check"
)

// New create a new port probe check
func New(target, port string) check.Check {
	c := &probeCheck{target: target, port: port}
	c.Setup("Probe was successful",
		"Error probing: %v",
		"dns_checker_probe_port",
		"Result of port probe 0 = error, 1 = OK",
		"target", "port")
	return c
}

type probeCheck struct {
	check.BaseCheck
	target string
	port   string
}

func (c *probeCheck) Execute() ([]interface{}, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", c.target, c.port))
	if conn != nil {
		_ = conn.Close()
	}
	return nil, err
}

func (c *probeCheck) Report(result []interface{}, err error, duration float64) {
	c.ReportResults(result, err, duration, c.target, c.port)
}
