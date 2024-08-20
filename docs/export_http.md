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

![dashboard_example.png](..%2Fcontrib%2Fdashboard_example.png)
![elinks.png](..%2Fcontrib%2Felinks.png)
![mobile.jpg](..%2Fcontrib%2Fmobile.jpg)

Data can be filtered by tags.
If you want TLS encryption or even basic authorization, you can do it by reverse proxy with config like this one
https://github.com/vodolaz095/dashboard/blob/master/contrib/nginx/dashboard.conf

It is worth notice data is updated in realtime using [Server Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
technology


Method 2. Via JSON endpoint
===================================
You can receive all dashboard data in JSON format by calling appropriate endpoint
http://localhost:3000/json

Filtering sensors by tags is supported, so this should work with default config
http://localhost:3000/json?kind=database&dialect=sql

Method 3: Via HTTP REST API
=====================================
You can receive all dashboard sensors' reading in JSON format by calling appropriate endpoint
http://localhost:3000/api/v1/sensor

Filtering sensors by tags is supported, so this should work with default config
http://localhost:3000/api/v1/sensor?a=b

You can load particular sensors readings via this call
http://localhost:3000/api/v1/sensor/load1


Method 4. Via Prometheus v4 Metrics scrapper endpoint
===================================
Endpoint `/metrics` exposes sensor readings in
[Prometheous v4](https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example)
format.

Filtering sensors by tags is supported, so this should work with default config
http://localhost:3000/metrics?kind=database&dialect=sql

Method 4. Via plain text (console+curl friendly metrics)
===================================
Endpoint `/text` exposes sensor readings in plain text format.
Filtering sensors by tags is supported, so this should work with default config
http://localhost:3000/text?kind=database&dialect=sql

Method 5. Via CSV metrics
===================================
Endpoint `/csv` exposes sensor readings in plain text format with CSV encoding, something like this:

| Name             | Minimum  | Value       | Maximum   | Error      | UpdatedAt | Description                                                                                                                             |
|------------------|----------|-------------|-----------|------------|-----------|-----------------------------------------------------------------------------------------------------------------------------------------|
| load1            | 0.0000   | 0.4700      | 0.0000    |            | 00:43:04  | Get system load average during last minute https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md              |
| load5            | 0.0000   | 0.4700      | 0.0000    |            | 00:43:04  | Get system load average during last 5 minutes https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md           |
| load15           | 0.0000   | 0.4700      | 0.0000    |            | 00:43:05  | Get system load average during last 15 minutes https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md          |
| process          | 0.0000   | 1204.0000   | 0.0000    |            | 00:43:05  | Number of currently running processes https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md                   |
| free_ram         | 500.0000 | 3462.0469   | 8000.0000 |            | 00:43:05  | Current free Random Access Memory volume in megabytes https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md   |
| free_home        | 0.0000   | 10620.6289  | 0.0000    |            | 00:43:06  | Size of free space in megabytes in /home directory https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md      |
| used_home        | 0.0000   | 109758.9336 | 0.0000    |            | 00:43:06  | Size of used disk space in megabytes in /home directory https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md |
| ratio_home       | 0.0000   | 91.1774     | 0.0000    |            | 00:43:06  | Free space ratio for /home directory in percents https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md        |
| echo             | 0.0000   | 6.0000      | 60.0000   |            | 00:43:06  | Get current second executing `date` command https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_shell.md                    |
| endpoint         | 0.0000   | 0.0000      | 0.0000    |            | 00:00:00  | Update value by incoming POST request https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_endpoint.md                       |
| redis            | 0.0000   | 0.0000      | 0.0000    | redis: nil | 00:43:06  | Get value of redis key a https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_redis.md                                       |
| mysql            | 1.0000   | 87.4750     | 100.0000  |            | 00:43:06  | Select random number from range https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_sql.md                                  |
| postgres         | 1.0000   | 5.6523      | 100.0000  |            | 00:43:07  | Select random number from range https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_sql.md                                  |
| thermal0         | 1.0000   | 44.0000     | 100.0000  |            | 00:43:07  | Get thermal sensor status from area 0 https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_file.md                           |
| redis subscriber | 1.0000   | 0.0000      | 100.0000  |            | 00:00:00  | Subscribe to redis channel and get values from it https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_redis.md              |


Filtering sensors by tags is supported, so this should work with default config
http://localhost:3000/csv?kind=database&dialect=sql
