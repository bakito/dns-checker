package run

import (
	"github.com/bakito/dns-checker/pkg/check"
)

func newExecution(check check.Check, target string, port *int) execution {
	ex := execution{
		check: check,
	}
	ex.target = target
	ex.port = port
	return ex
}

type execution struct {
	check.Result
	tp
	check check.Check
}

type tp struct {
	target string
	port   *int
}
