Squid Prometheus exporter
--------------------------

Exports squid metrics in Prometheus format

Usage:
------
Simple usage:

    squid-exporter -squid-hostname "localhost" -squid-port 3128

Configure Prometheus to scrape metrics from `localhost:9301/metrics`

    - job_name: squid
      # If prometheus-node-exporter is installed, grab stats about the local
      # machine by default.
      target_groups:
        - targets: ['localhost:9301']

To get all the paramteres

    squid-exprter -help


Features:
---------

- [ ] Expose Squid counters
  -  [x] Client HTTP
  -  [x] Server HTTP
  -  [ ]  Server ICP
  -  [ ]  Other
- [ ] Histograms
- [ ] Other metics
- [ ] Squid Authentication

Contribution:
-------------

Pull request and issues are very welcome.

Copyright:
----------

[MIT License](https://opensource.org/licenses/MIT)


