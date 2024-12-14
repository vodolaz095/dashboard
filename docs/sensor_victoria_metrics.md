Victoria Metrics Sensor
=========================================

Sensor sends periodical queries to Victoria Metrics using PromQL/MetricsQL instant query
https://docs.victoriametrics.com/url-examples/#apiv1query

Shared sensor parameters are explained in
[sensor_shared.md](https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_shared.md)
file.

All config parameters for sensors are depicted in this file
[sensor.go](https://github.com/vodolaz095/dashboard/blob/master/config/sensor.go)
with comments explaining things.

```yaml

  - name: victoria metrics1
    type: victoria_metrics
    description: "Fetch timeseries data from Victoria Metrics"
    link: "https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_victoria_metrics.md"
    endpoint: "http://localhost:8428/"
    query: '{instance="steel"}'
    filter:
      __name__: "file1"
    refresh_rate: 5s
    tags:
      dialect: promql
      kind: database

  - name: victoria metrics2
    type: victoria_metrics
    description: "Fetch timeseries data from Victoria Metrics"
    link: "https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_victoria_metrics.md"
    endpoint: "http://localhost:8428/"
    query: '{instance="steel", __name__="file1"}'
    refresh_rate: 5s
    tags:
      dialect: promql
      kind: database
```

Important parameters are `query` used to list timeseries we need, and `filter` - used to pick metrics we use to extract last data.
See https://docs.victoriametrics.com/metricsql/ query syntax.
For example, this query `{instance="steel"}`, returns few series like this ones:

```
go_gc_pauses_seconds_bucket{instance="steel",le="0.00016384"}
go_gc_duration_seconds_count{instance="steel"}
go_cpu_count{instance="steel"}
file1{instance="steel",kind="file",origin="feeder"}
```

and only the first one with `name` (as shown in filter by name `__name__`) will be used for getting value.

Or, actually, you can unset filter and use query like this - `{instance="steel", __name__="file1"}`.
If query returns few metrics, filter allows to pick first matching - it can make sense in some situations.
