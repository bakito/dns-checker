package run

import (
	"github.com/bakito/dns-checker/pkg/check"
)

func newExecution(check check.Check, address check.Address) execution {
	ex := execution{
		check: check,
	}
	ex.Address = address
	return ex
}

type execution struct {
	check.Result
	check.Address
	check check.Check
}
