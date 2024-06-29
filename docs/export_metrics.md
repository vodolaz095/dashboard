Exporting sensor data into Prometheus v4 or InfluxDB via scrapper
========================

Endpoint `/metrics` exposes sensor readings in
[Prometheous v4](https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example)
format.

Data can be filtered by tags. So, if you have configuration as defined here
https://github.com/vodolaz095/dashboard/blob/master/contrib/dashboard.yaml
and dashboard running on http://localhost:3000, you can get all metrics via URL like this:

http://localhost:3000/metrics 

and you can filter readings only to mysql and postgres database related via URL like this:

http://localhost:3000/metrics?dialect=sql&kind=database


Consuming data from Prometheus
=========================
TODO


Consuming data using InfluxDB v2 scrappers
=========================

Process of setting InfluxDB 2 scrappers is explained here
https://docs.influxdata.com/influxdb/v2/write-data/no-code/scrape-data/

