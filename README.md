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

Example dashboard screenshots
=====================
![dashboard_example.png](contrib%2Fdashboard_example.png)
![elinks.png](contrib%2Felinks.png)
![mobile.jpg](contrib%2Fmobile.jpg)

Architecture
=====================
Application has list of in-memory sensors.
HTTP server load sensor values from memory, so databases cannot be DDoSed.
Background process updates sensors' readings using [defered queue](https://github.com/vodolaz095/dqueue), 
separate goroutines keep readings updates via external events (http requests, redis subscription messages, etc.).



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
6. DDOS (distributed denial of service attacks) proof - sensors readings are updated in memory by background goroutines 
   and served by HTTP server from memory. No matter how many clients open dashboard - they receive values from memory,
   no extra calls to database and other resources are issued. 
7. Database access credentials, tokens, passwords and other sensitive data is concealed from visitors.


Possible alternatives
======================
Requirement - simple realtime dashboard with list of sensors readings, containing actual numerical values 
and some technical background (minimum, maximum, link to wiki) available for all team members.
No historical charts required. Numerical values can be extracted by SQL requests, HTTP requests and so on.
Data can be (but not required to) stored to some 3rd party time series database.
`Pros` here means what does alternative have, and my dashboard - do not, and `cons` describes why
alternative was discarded.


Alternative 1. [monit](https://mmonit.com/monit/)
- Pros: Easy to setup, lots of plugins, [powerfull scripting language](https://mmonit.com/monit/documentation/monit.html#MYSQL) to write tests.
- Cons: dashboard requires password based authorization, with misconfigured board user can start/stop services.
  Complicated scripting language, writing sensor extracting metrics from *SQL database was painful since it required 
  to write shell scripts...
- Conclusion: overcomplicated.

Alternative 2. [grafana](https://grafana.com/)
- Pros: popular system with years of production service
- Cons: too complicated, hard to setup, authorization required for users to view charts data.
  3rd party tools are required to extract sensors' readings from observable servers.
  3rd party tools (Influxdb, Prometheus, etc) are required to store data being visualized.
- Conclusion: overcomplicated.

Alternative 3. [zabbix](https://www.zabbix.com/)
- Pros: popular system with years of production service
- Cons: too complicated, hard to setup, authorization required for users to view charts data.
- Conclusion: overcomplicated.


Alternative 4. [netdata](https://www.netdata.cloud/)
- Pros: lot of plugins, fancy UI, quite easy to setup.
- Cons: webui is quite heavy, works slow via 3G mobile connection. I just need table with few actual readings.
- Conclusion: overcomplicated.


Alternative 5. [Influxdb v2](https://docs.influxdata.com/influxdb/v2/) + [Telegraf](https://docs.influxdata.com/telegraf/v1/)
- Pros: telegraf has lots of inputs and outputs, which are quite easy to configure. Easy to make dashboards in Influxdb. Historical data available. 
  Our dashboard can [send data directly to Influxdb via wire protocol](docs%2Fexport_influx.md).
- Cons: Influxdb doesn't render UI without authorization. Loading simple dashboard eats 12+ mb of traffic. Telegraf 
  does not have easy to use web UI to read actual data manually.
- Conclusion: overcomplicated.

Quickstart
======================

1. Obtain suitable binary from https://github.com/vodolaz095/dashboard/releases
2. Copy configuration example with comments [dashboard.yaml](contrib%2Fdashboard.yaml),
   and change parameters to your own and start application in this way:

   ```shell
     
     $ dashboard /path/to/dashboard_config.yaml
     
   ```

3. See [deployment](docs%2Fdeployment.md) how to ran application for production.

Security
=============================
1. All sensor readings are available to all dashboard users, while database access credentials and database queries are concealed
2. Dashboard WebUI access can be restricted either by reverse proxy, or it can be served only in local network - so
   if somebody can view this dashboard - he/she is allowed to do to.
3. Updating dashboard is performed automatically
4. Configuring dashboard is done by system administrators, allowed to work with data required.


Configuration
=======================
- [Defining database connections pool in config](docs%2Fconnection_pool.md)
- [Logging](docs%2Flogging.md)
- [Dashboard customization](docs%2Fui_customization.md)
- [Exporting data via HTTP transport](docs%2Fexport_http.md)
- [Exporting data via redis publishers](docs%2Fexport_redis_.md)
- [Exporting sensor data into InfluxDB via wire protocol](docs%2Fexport_influx.md)
- [Exporting sensor data into Prometheus/InfluxDB via scrapper](docs%2Fexport_metrics.md)
- [Linking few dashboards via redis pub/sub](docs%2Flinking_via_redis.md)
- [Deployment](docs%2Fdeployment.md)

Sensors documentation
==========================
- [Shared sensor parameters](docs%2Fsensor_shared.md)
- [Shell sensor](docs%2Fsensor_shell.md)
- [Linux system sensor](docs%2Fsensor_linux_system.md)
- [MySQL/PostgreSQL sensor](docs%2Fsensor_sql.md)
- [Redis sensor (synchronous and subscriber)](docs%2Fsensor_redis.md)
- [File sensor which reads data from file](docs%2Fsensor_redis.md)
- [Incoming HTTP POST request / endpoint / webhook sensor](docs%2Fsensor_endpoint.md)
- [Outgoing HTTP Request (CURL) sensor](docs%2Fsensor_curl.md)
- [Creating your own sensor](docs%2Fsensor_your_own.md)


Development using golang compiler on host machine
=============================
Application requires [Golang 1.22.0](https://go.dev/dl/) and [GNU Make](https://www.gnu.org/software/make/) installed.

```shell

# ensure development tools in place
$ make tools

# ensure golang modules are installed
$ make deps

# start application for development using configuration from contrib/dashboard.yaml on http://localhost:3000
$ make start

# build production grade binary at `build/dashboard`
$ make build

```

MySQL, PostgreSQL, Redis and Influxdb can be started by docker/podman

```shell

# start development databases (depending on what container engine is available)  
$ make docker/resource
$ make podman/resource

```

Development using docker + docker compose
=============================
[GNU Make](https://www.gnu.org/software/make/), [Docker engine](https://docs.docker.com/engine/install/) with
[compose plugin](https://docs.docker.com/compose/install/linux/) should be installed.
Installing golang toolchain on host machine is not required.

```shell

# start development databases and build and start application on http://localhost:3001 
$ make docker/up

# start development databases  
$ make docker/resource

# stop all
$ make docker/down

# prune all development environment
$ make docker/prune


```


Development using podman + podman-compose
=============================
Installing golang toolchain on host machine is not required.
Tested on Fedora 39, 40 and Centos 9 Stream. 

```shell

# install development environment
$ sudo dnf install make podman podman-compose podman-plugins containernetworking-plugins

# start development databases and build and start application on http://localhost:3001
$ make podman/up

# start development databases  
$ make podman/resource

# stop all
$ make podman/down

# prune all development environment
$ make podman/prune

```

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
