[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/boynux/squid-exporter)

[![Github Actions](https://github.com/boynux/squid-exporter/actions/workflows/release.yml/badge.svg)](https://github.com/boynux/squid-exporter/actions/workflows/release.yml)
[![Github Docker](https://github.com/boynux/squid-exporter/actions/workflows/docker.yml/badge.svg)](https://github.com/boynux/squid-exporter/actions/workflows/docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/boynux/squid-exporter)](https://goreportcard.com/report/github.com/boynux/squid-exporter)
[![Maintainability](https://api.codeclimate.com/v1/badges/a99a88d28ad37a79dbf6/maintainability)](https://codeclimate.com/github/boynux/squid-exporter)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=3TH7YAMMEC5L4&source=url)



Squid Prometheus exporter
--------------------------

Exports squid metrics in Prometheus format

**NOTE**: From release 1.0 metric names and some parameters has changed. Make sure you check the docs and update your deployments accordingly!

New
-----

* Using environment variables to configure the exporter
* Adding custom labels to metrics

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

The following environment variables can be used to override default parameters:

```
SQUID_EXPORTER_LISTEN
SQUID_EXPORTER_METRICS_PATH
SQUID_HOSTNAME
SQUID_PORT
SQUID_LOGIN
SQUID_PASSWORD
SQUID_EXTRACTSERVICETIMES
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
- [ ] Expose Squid service times
  - [x] HTTP requests
  - [x] Cache misses
  - [x] Cache hits
  - [x] Near hits
  - [ ] Not-Modified replies
  - [x] DNS lookups
  - [ ] ICP queries
- [ ] Histograms
- [ ] Other metrics
- [x] Squid Authentication (Basic Auth)

FAQ:
--------

- Q: Metrics are not reported by exporter
- A: That usually means the exporter cannot reach squid server or the config manager permissions are not set corretly. To debug and mitigate:
  - First make sure the exporter service can reach to squid server IP Address (you can use telnet to test that)
  - Make sure you allow exporter to query the squid server in config you will need something like this (`172.20.0.0/16` is the network for exporter, you can also use a single IP if needed):
  ```
  #http_access allow manager localhost
  acl prometheus src 172.20.0.0/16
  http_access allow manager prometheus
  ```

- Q: Why `process_open_fds` metric is not exported?
- A: This usualy means exporter don't have permission to read `/proc/<squid_proc_id>/fd` folder. You can either

1. _[recommended]_ Set `CAP_DAC_READ_SEARCH` capability for squid exporter process (or docker). (eg. `sudo setcap 'cap_dac_read_search+ep' ./bin/squid-exporter`)
2. _[not recommended]_ Run the exporter as root.

Contribution:
-------------

Pull request and issues are very welcome.

If you found this program useful please consider donations [![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=3TH7YAMMEC5L4&source=url)

Copyright:
----------

[MIT License](https://opensource.org/licenses/MIT)


