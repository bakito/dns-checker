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
| TARGET | The DNS target host to check | X |  |
| TARGET_PORT | The target port. If defined the port will be probed | O | |
| INTERVAL | The check intercal in seconds | O | 30 |
| METRICS_PORT | The port for the metrics service | O | 2112 |
| LOG_LEVEL | The log level (panic, fatal, error, warn, info, debug, trace)| O | info |

## Metrics

Exposes metrics under localhost:2112/metrics

| Key | Description  
| :---: | --- |
| dns_checker_check_dns | The result of the DNS check 1 OK / 0 error |
| dns_checker_check_dns_duration | The duration result of the DNS check in milliseconds|
| dns_checker_probe_port | The result of the port probe check 1 OK / 0 error |
| dns_checker_probe_port_duration | The duration result of the port probe check in milliseconds |

