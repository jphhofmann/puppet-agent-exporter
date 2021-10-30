# Puppet Agent Exporter

`Puppet Agent Exporter` is puppet agent report statistics exporter for Prometheus. This is an fork, as the original project is no longer maintained and the previous code base needed some work for future usage.

The `Puppet Agent Exporter` reads last run summary report and exposes them for Prometheus consumption.

## Building and running

```
go build
```

### Flags

Name                    | Default                                     | Description
----------------------- | ------------------------------------------- | -----------
web.listen-address      | localhost:9001                              | Address on which to expose metrics and web interface.
web.telemetry-path      | /metrics                                    | Path under which to expose metrics.
namespace               | puppet                                      | The namespace of metrics.
puppet.last-run-summary | /var/lib/puppet/state/last_run_summary.yaml | Path to the puppet's last run summary report.
puppet.last-run-report  |                                             | Path to the puppet's last run full report - need only if you want to get info from full report.
puppet.disabled-lock    |                                             | Path to the puppet's agent disabled lock-file.

## What's exported?
It exports statistics from standard Puppet report (https://puppet.com/blog/puppet-monitoring-how-to-monitor-success-or-failure-of-puppet-runs).
