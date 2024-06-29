Exporting data via influxdb
===========================================

Via wire protocol
===========================================
Firstly, we need to obtain influxdb credentials with write access to required bucket.
See ![screenshot](https://github.com/vodolaz095/dashboard/blob/master/docs/influx.png) for insights.

We need token with **write** access to bucket required.
Then, we can update [config](https://github.com/vodolaz095/dashboard/blob/master/contrib/dashboard.yaml) adding influxdb
access credentials, so dashboard application will establish HTTP connection with influxdb API to send 
data in [wire protocol](https://docs.influxdata.com/influxdb/v2/reference/syntax/line-protocol/) format

```yaml

# https://github.com/vodolaz095/dashboard/blob/master/docs/export_influx.md
influx:
  endpoint: http://127.0.0.1:8086
  token: "-l3Y5tIHGJAxXv_Rs5kJ4kAfPbgmf3WPmFUTDuKmD3Z9gp29E7e188-dIt5MAKhSTzv1J6v_pkPuVdIbXqdL1w=="
  organization: dashboard
  bucket: dashboard

```

Sensor tags will be attached to data points, so, if you have 2 sensors like these:

```yaml

sensor:
  - name: mysql
    type: mysql
    connection_name: "mysql@container"
    query: "SELECT rand()*99+1 as random"
    refresh_rate: 5s
    tags:
      dialect: sql
      kind: database

  - name: postgres
    type: postgres
    connection_name: "postgres@container"
    query: "SELECT random()*99+1 as random"
    refresh_rate: 5s
    tags:
      dialect: sql
      kind: database


```

wire protocol data will be like this
```
mysql,dialect=sql,kind=database value=12 1556813561098000000
postgres,dialect=sql,kind=database value=53 1556813561098000000
```



Via scrapper
============================================
Endpoint `/metrics` exposes sensor readings in
[Prometheous v4](https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example)
format.

Data can be filtered by tags. So, if you have configuration as defined here
https://github.com/vodolaz095/dashboard/blob/master/contrib/dashboard.yaml
and dashboard running on http://localhost:3000, you can get all metrics via URL like this:

http://localhost:3000/metrics

and you can filter readings only to mysql and postgres database related via URL like this:

http://localhost:3000/metrics?dialect=sql&kind=database



Process of setting InfluxDB 2 scrappers is explained here:
https://docs.influxdata.com/influxdb/v2/write-data/no-code/scrape-data/

