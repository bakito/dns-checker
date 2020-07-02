[![Docker Repository on Quay](https://quay.io/repository/bakito/dns-checker/status "Docker Repository on Quay")](https://quay.io/repository/bakito/dns-checker) [![Go Report Card](https://goreportcard.com/badge/github.com/bakito/dns-checker)](https://goreportcard.com/report/github.com/bakito/dns-checker)

# DNS Checker

Check and reports host names and port.


## Run

```bash
docker run -p 2112:2112 -e TARGET=<target-host> quay.io/bakito/dns-checker
```

## Env Variables
| Name | Description | Required | Default 
| :---: | --- | :---: | :---: |
| TARGET | The DNS target hosts to check. ',' separated host(:port) list | X |  |
| INTERVAL | The check intercal in seconds | O | 30 |
| METRICS_PORT | The port for the metrics service | O | 2112 |
| LOG_LEVEL | The log level (panic, fatal, error, warn, info, debug, trace)| O | info |
| METRICS_HISTOGRAM_BUCKETS | Custom histogram metric buckets  | O | "0.002,0.005,0.01,0.025,0.05,0.1,0.25,0.5,1,2.5,5,10,20" |
| METRICS_SUMMARY_OBJECTIVES | Custom summary metric objectives | O | "0.5:0.05,0.9:0.01,0.99:0.001" |


## Metrics

Exposes metrics under localhost:2112/metrics

| Key | Description  
| :---: | --- |
| dns_checker_check_dns | The result of the DNS check 1 OK / 0 error |
| dns_checker_check_dns_duration | The duration result of the DNS check in milliseconds|
| dns_checker_check_dns_summary | The summary metric of the duration|
| dns_checker_check_dns_histogram | The histogram metric of the duration |
| dns_checker_probe_port | The result of the port probe check 1 OK / 0 error |
| dns_checker_probe_port_duration | The duration result of the port probe check in milliseconds |
| dns_checker_check_port_summary | The summary metric of the duration|
| dns_checker_check_port_histogram | The histogram metric of the duration |

