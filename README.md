[![Build Status](https://travis-ci.org/boynux/squid-exporter.svg?branch=master)](https://travis-ci.org/boynux/squid-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/boynux/squid-exporter)](https://goreportcard.com/report/github.com/boynux/squid-exporter)
[![Maintainability](https://api.codeclimate.com/v1/badges/a99a88d28ad37a79dbf6/maintainability)](https://codeclimate.com/github/boynux/squid-exporter)

Squid Prometheus exporter
--------------------------

Exports squid metrics in Prometheus format

**NOTE**: From release 1.0 metric names and some parameters has changed. Make sure you check the docs and update your deployments accordingly!

Usage:
------
Simple usage:

    squid-exporter -squid-hostname "localhost" -squid-port 3128

[Configure Prometheus](https://github.com/boynux/squid-exporter/blob/master/prometheus/prometheus.yml) to scrape metrics from `localhost:9301/metrics`

    - job_name: squid
      # squid-exporter is installed, grab stats about the local
      # squid instance.
      target_groups:
        - targets: ['localhost:9301']

To get all the parameteres, command line arguments always override default and environment variables configs:

    squid-exporter -help

The following environmnet variables can be used to override default parameters:

```
SQUID_EXPORTER_LISTEN
SQUID_EXPORTER_METRICS_PATH
SQUID_HOSTNAME
SQUID_PORT
SQUID_LOGIN
SQUID_PASSWORD
```

Usage with docker:
------
Basic setup assuming Squid is running on the same machine:

    docker run --net=host -d boynux/squid-exporter

Setup with Squid running on a different host

    docker run -p 9301:9301 -d boynux/squid-exporter -squid-hostname "192.168.0.2" -squid-port 3128 -listen ":9301"

With environment variables

    docker run -p 9301:9301 -d -e SQUID_PORT="3128" -e SQUID_HOSTNAME="192.168.0.2" -e SQUID_EXPORTER_LISTEN=":9301" boynux/squid-exporter


Build:
--------

This project is written in Go, so all the usual methods for building (or cross compiling) a Go application would work.

If you are not very familiar with Go you can download the binary from [releases](https://github.com/boynux/squid-exporter/releases).

Or build it for your OS:

`go install https://github.com/boynux/squid-exporter`

then you can find the binary in: `$GOPATH/bin/squid-exporter`

Features:
---------

- [ ] Expose Squid counters
  -  [x] Client HTTP
  -  [x] Server HTTP
  -  [x] Server ALL
  -  [x] Server FTP
  -  [x] Server Other
  -  [ ] ICP
  -  [ ] CD
  -  [x] Swap
  -  [ ] Page Faults
  -  [ ] Others
- [ ] Histograms
- [ ] Other metrics
- [x] Squid Authentication (Basic Auth)

Contribution:
-------------

Pull request and issues are very welcome.

Copyright:
----------

[MIT License](https://opensource.org/licenses/MIT)


