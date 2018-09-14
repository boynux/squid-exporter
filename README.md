[![Build Status](https://travis-ci.org/boynux/squid-exporter.svg?branch=master)](https://travis-ci.org/boynux/squid-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/boynux/squid-exporter)](https://goreportcard.com/report/github.com/boynux/squid-exporter)

Squid Prometheus exporter
--------------------------

Exports squid metrics in Prometheus format

**NOTE**: From release 1.0 metric names and some parameters has changed. Make sure you check the docs and update your deployments accordingly!

Usage:
------
Simple usage:

    squid-exporter -squid-hostname "localhost" -squid-port 3128

Configure Prometheus to scrape metrics from `localhost:9301/metrics`

    - job_name: squid
      # squid-exporter is installed, grab stats about the local
      # squid instance.
      target_groups:
        - targets: ['localhost:9301']

To get all the parameteres

    squid-exprter -help


Usage with docker:
------
Basic setup assuming Squid is running on the same machine:

    docker run --net=host -d boynux/squid-exporter

Setup with Squid running on a different host

    docker run -p 9301:9301 -d boynux/squid-exporter -squid-hostname "192.168.0.2" -squid-port 3128 -listen ":9301"


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
- [ ] Squid Authentication

Contribution:
-------------

Pull request and issues are very welcome.

Copyright:
----------

[MIT License](https://opensource.org/licenses/MIT)


