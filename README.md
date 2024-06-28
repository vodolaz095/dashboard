Vodolaz095's Dashboard
======================
Minimalistic and DDOS-proof Golang powered dashboard

Usage example
======================
Consider your business depends on MySQL database of CRM, PostgreSQL database for shipping,
few 3rd party APIs (like get balance of bank account), redis database with real time machinery state and few scripts
you are running on servers on site to see its working. 
So, important readings can be, for example, 

- number of active orders in MySQL database of CRM extracted by query like 
```sql
SELECT COALESCE(count(orders.id),0) as "orders_pending"
FROM orders
WHERE DATE(orders.created_at) = CURDATE() and orders.status=1;
```

- number of completed orders extracted this way
```sql
SELECT COALESCE(count(orders.id),0) as "orders_completed"
FROM orders
WHERE DATE(orders.created_at) = CURDATE() and orders.status=2;
```

- query like this (with stored procedure) is used to count active deliveries in PostgreSQL
```sql
SELECT doCountActiveDeliveries(CURDATE());
```
 
- real time machinery readings are extracted by redis commands like this one
```
127.0.0.1:6379> hget reactor1 power_output
```

- bank account can be checked by sending HTTP POST request to, for example, https://example.org/api/v1/rpc

Checking every parameter separately can be automated by scripts, but making it all easy and in one place
can be complicated. It can be wise idea to combine all these readings in single dashboard available for all stakeholders 
and important employees, so they can have eagle's eye perspective on what is happening. It can be wise to conceal 
some technical data (like database connection strings) but, in general, all important data should be available 
on single page in a way it can be understood by general audience without technical skills.

Example dashboards
=====================
![dashboard_example.png](contrib%2Fdashboard_example.png)

Architecture
=====================


Main features
======================
1. Manifold of very hackable sensors - MySQL/PostgreSQL queryes, Redis sync and subscription, file, remote HTTP endpoint, 
   periodical shell command execution, local HTTP POST endpoint updated by remote script with secret token.
2. Single cross-platform binary with simple `yaml` encoded config
3. Light-weight (dashboard has ~1 kb [style.css](assets%2Fstyle.css), ~1 kb [feed.js](assets%2Ffeed.js) and ~ 5kb 
   main page)- works ok even on IPhone 6 and 2013 year Android Smartphones
4. Real time updates using [SSE](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
5. JSON and [Prometheous v4](https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example)
   endpoints to read sensors readings
6. DDOS (distributed denial of service attacks) proof - sensors readings are updated in memory by background process 
   and served by HTTP server from memory. No matter how many clients opens dashboard - they receive values from memory,
   no extra calls to database and other resources are issued. 
7. Database access credentials, tokens, passwords and other sensitive data is concealed from visitors.


Quickstart
======================
Big configuration example with comments:
[dashboard.yaml](contrib%2Fdashboard.yaml)

List of contents
=======================
- [Defining database connections pool in config](docs%2Fconnection_pool.md)
- [Shared sensor parameters](docs%2Fsensor_shared.md)
- [Shell sensor](docs%2Fsensor_shell.md)
- [MySQL/PostgreSQL sensor](docs%2Fsensor_sql.md)
- [Redis sensor (synchronous and subscriber)](docs%2Fsensor_redis.md)
- [Read data from file sensor](docs%2Fsensor_redis.md)
- [Incoming HTTP POST Endpoint (webhook) sensor](docs%2Fsensor_endpoint.md)
- [HTTP Request (CURL) sensor](docs%2Fsensor_curl.md)
- [Creating your own sensor](docs%2Fsensor_your_own.md)
- [Logging](docs%2Flogging.md)
- [Dashboard customization](docs%2Fui_customization.md)
- [Exporting data via HTTP transport](docs%2Fexport_http.md)
- [Exporting data via redis publishers](docs%2Fexport_redis_publisher.md)
- [Exporting sensor data into InfluxDB via wire protocol](docs%2Fexport_influx.md)
- [Exporting sensor data into Prometheus/InfluxDB via scrapper](docs%2Fexport_metrics.md)
- [Linking few dashboards via redis pub/sub](docs%2Flinking_via_redis.md)


Security
=============================
1. All sensor readings are available to all dashboard users, while database access credentials and database queries are concealed
2. Dashboard WebUI access can be restricted either by reverse proxy, or it can be served only in local network - so
   if somebody can view this dashboard - he/she is allowed to do to.
3. Updating dashboard is performed automatically
4. Configuring dashboard is done by system administrators, allowed to work with data required.


Deployment
=============================
NGINX as reverse proxy, encryption and authorization is done by NGINX.
Configuration example - [dashboard.conf](contrib%2Fnginx%2Fdashboard.conf)
Good read:
- https://nginx.org/ru/docs/http/ngx_http_proxy_module.html
- https://stackoverflow.com/questions/23844761/upstream-sent-too-big-header-while-reading-response-header-from-upstream
- https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-http-basic-authentication/


License
===================================
MIT License

Copyright (c) 2024 Остроумов Анатолий

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
