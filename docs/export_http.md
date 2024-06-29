Exporting data via HTTP transport
===================================

All methods allow data to be can be filtered by tags. So, if you have configuration as defined here
https://github.com/vodolaz095/dashboard/blob/master/contrib/dashboard.yaml
and dashboard running on http://localhost:3000, you can get all metrics via URL like this:

http://localhost:3000/metrics 

and you can filter readings only to mysql and postgres database related via URL like this:

http://localhost:3000/metrics?dialect=sql&kind=database


Method 1. Via Browser
===================================

Just open http://localhost:3000 via browser to see something like this

![dashboard_example.png](contrib%2Fdashboard_example.png)

Data can be filtered by tags.
If you want TLS encryption or even basic authorization, you can do it by reverse proxy with config like this one
https://github.com/vodolaz095/dashboard/blob/master/contrib/nginx/dashboard.conf

It is worth notice data is updated in realtime using [Server Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
technology


Method 2. Via JSON endpoint
===================================
You can recieve all dashboard data in JSON format by calling appropriate endpoint

http://localhost:3000/json

Filtering sensors by tags is supported, so this should work with default config
http://localhost:3000/json?kind=database&dialect=sql

Method 3. Via Prometheus v4 Metrics scrapper endpoint
===================================
Endpoint `/metrics` exposes sensor readings in
[Prometheous v4](https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example)
format.

Filtering sensors by tags is supported, so this should work with default config
http://localhost:3000/metrics?kind=database&dialect=sql
