[![Build Status](https://travis-ci.org/boynux/squid-exporter.svg?branch=master)](https://travis-ci.org/boynux/squid-exporter)

Squid Prometheus exporter
--------------------------

Exports squid metrics in Prometheus format

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
  -
- [ ] Histograms
- [ ] Other metrics
- [ ] Squid Authentication

Contribution:
-------------

Pull request and issues are very welcome.

Copyright:
----------

[MIT License](https://opensource.org/licenses/MIT)


