[![Build Status](https://travis-ci.com/bakito/dns-checker.svg?branch=master)](https://travis-ci.com/bakito/dns-checker) [![Docker Repository on Quay](https://quay.io/repository/bakito/dns-checker/status "Docker Repository on Quay")](https://quay.io/repository/bakito/dns-checker) [![Go Report Card](https://goreportcard.com/badge/github.com/bakito/dns-checker)](https://goreportcard.com/report/github.com/bakito/dns-checker)

# DNS Checker

Check and reports host names and port.


## Run

```bash
docker run -p 2112:2112 -e TARGET=<target-host> quay.io/bakito/dns-checker
```

## Env Variables
| Name | Description | Required | Default 
| :---: | --- | :---: | :---: |
| TARGET | The DNS target hosts to check. ',' separated host(:port) list. Env variables can be used here with notation '${ENV_VAR_NAME}' | X |  |
| INTERVAL | The check interval as duration | O | 30s |
| TIMEOUT | The check timeout as duration | O | 10s |
| WORKER | The number of workers to be used for the checks | O | 10 |
| METRICS_PORT | The port for the metrics service | O | 2112 |
| LOG_LEVEL | The log level (panic, fatal, error, warn, info, debug, trace)| O | info |
| LOG_JSON | Enables json log format if set to tue | O | false |
| METRICS_HISTOGRAM_BUCKETS | Custom histogram metric buckets  | O | "0.002,0.005,0.01,0.025,0.05,0.1,0.25,0.5,1,2.5,5,10,20" |
| METRICS_SUMMARY_OBJECTIVES | Custom summary metric objectives | O | "0.5:0.05,0.9:0.01,0.99:0.001" |


## Metrics

Exposes metrics under localhost:2112/metrics

| Key | Description  
| :---: | --- |
| dns_checker_check | The result of the checks 1 OK / 0 error |
| dns_checker_check_error | check resulted in an error 1 = error /  0 = OK |
| dns_checker_check_duration | The duration result of the check in milliseconds|
| dns_checker_check_summary | The summary metric of the duration|
| dns_checker_check_histogram | The histogram metric of the duration |

### Metrics Labels

Each metric has the following labels
| Name | Description  
| :---: | --- |
| target | The target of the checks |
| port | The port of the checks (may be empty) |
| check_name | The name of the check |
| version | The application version  |