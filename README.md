Vodolaz095's Dashboard
======================
Goland powered dashboard

Usage example
======================
Consider your business depends on MySQL database of CRM, PostgreSQL database for shipping,
few 3rd party APIs (like get balance of bank account), redis database with real time machinery state and few scripts
you are running on servers on site to see its working. It can be wise idea to combine all these readings in
single dashboard available for all stakeholders and important employees, so they can have eagle's eye perspective 
on what is happening. It can be wise to conceal some technical data (like database connection strings) but, in general,
all important data should be available on single page in a way it can be understood by general audience without technical 
skills.

Example dashboards
=====================



Main features
======================
1. Manifold of very hackable sensors
2. Single cross-platform binary with simple `yaml` powered config
3. Light-weight (dashboard has ~1 kb [style.css](assets%2Fstyle.css), ~1 kb [feed.js](assets%2Ffeed.js) and ~ 5kb 
   main page)- works ok even on IPhone 6 and 2013 year Android Smartphones
4. Real time updates using [SSE](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
5. JSON and [Prometheous v4](https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example)
   endpoints to read sensors readings
6. DDOS proof - sensors readings are updated in memory by background process and served by HTTP server from memory


Sensors and their configuration examples
=========================

All config parameters for sensors are depicted in this file [sensor.go](config%2Fsensor.go)
with comments explaining things.
Important parameter is `Tags` - they can be used to generate dashboard with partial data provided.
For example, you have 3 databases monitored with 10+ sensors and they all have tag `kind=database`, 
while few of them has tag `unit=sales` and other - `unit=shipping`. Consider your dashboard is working on
`dashboard.example.org`. So, dashboard with all databases will be available on
`dashboard.example.org?kind=database`, sales unit database dashboards - on
`dashboard.example.org?kind=database&unit=sales` and shipping ones - on
`dashboard.example.org?kind=database&unit=shipping`. So, applying query parameter you can filter sensors to ones
having **ALL** tags matching query.

Also important parameter is `link` - it can be, for example, hyperlink to Wiki/Confluence page with explanation which 
procedures required to perform when, for example, reactor temperature is too high.

There are two parameters present - `A` and `B` - which can be used to apply linear transformation to sensor reading.
```
f(x) = A*x+B
```
For example, `A`=0.5555 `B`=-17.7778 allows converting sensor readings from Fahrenheit to Celsius degrees. 

Parameter `refresh_rate` defines how often sensor should be refreshed, it accepts strings that
[time#ParseDuration](https://pkg.go.dev/time#ParseDuration) understands, so
`5h, 1m, 10s` are OK.

***MySQL/PostgreSQL query sensor***

This sensor queries database with user defined SQL request periodically.
SQL request **SHOULD** return **SINGLE** value parsable as **FLOAT64**.
Background process periodically calls query and cache result in memory, and webserver 
process serves data from memory, so database cannot be DDoSed.
Configuration examples:

```yaml

  - name: mysql
    type: mysql
    description: "Select random number from range"
    link: "https://dev.mysql.com/doc/refman/8.0/en/mathematical-functions.html#function_rand"
    connection_name: "mysql@container"
    query: "SELECT rand()*99+1 as random"
    minimum: 1
    maximum: 100
    refresh_rate: 5s
    tags:
      dialect: sql
      kind: database

  - name: postgres
    type: postgres
    description: "Select random number from range"
    link: "https://www.postgresql.org/docs/current/functions-math.html"
    connection_name: "postgres@container"
    query: "SELECT random()*99+1 as random"
    minimum: 1
    maximum: 100
    refresh_rate: 5s
    tags:
      dialect: sql
      kind: database


  - name: AnatolijCaloriesLeft
    type: mysql
    description: "Сколько калорий осталось для Анатолия"
    link: "https://eda.example.org"
    connection_name: "eda"
    query: >
       SELECT COALESCE(metadata.value-SUM(calories),0) as "calories_left"
       FROM meals                         
       LEFT JOIN metadata on meals.username = metadata.username
       WHERE DATE(created_at) = CURDATE() and metadata.name = "calories" and metadata.username='vodolaz095';
    minimum: 1
    maximum: 2100
    refresh_rate: 30s
    tags:
       kind: eda
       user: vodolaz095


```

***Redis Synchronous Query Sensor***

This sensor periodically executes query `get a` to read value of a key `a`.
It is possible to run LUA stored procedures in redis and get their values as sensor readings

```yaml

- name: redis
  type: redis
  description: "Get value of redis key a"
  link: "https://example.org"
  connection_name: "redis@container"
  query: "get a"
  refresh_rate: 5s
  tags:
    kind: database

```

***Redis Subscriber Sensor***

This sensor subscribes to redis database channels and reads updates provided as float numbers.

```yaml

  - name: redis subscriber
    type: subscriber
    description: "Subscribe to redis channel and get values from it"
    channel: "vodolaz095/dashboard/subscriber/value"
    connection_name: "subscribe2redis@container"
    value_only: true
    minimum: 1
    maximum: 100
    tags:
      a: b
      c: d
      kind: database
```

can be updated by executing this redis command

```
  PUBLISH vodolaz095/dashboard/subscriber 47.1
```

```shell

	$ redis-cli publish vodolaz095/dashboard/subscriber `date "+%S"`

```

If we want to provide more data (not only value), we can parse messages as jsons


```yaml
  - name: redis subscriber
    type: subscriber
    description: "Subscribe to redis channel and get values from it"
    channel: "vodolaz095/dashboard/subscriber/all"
    connection_name: "subscribe2redis@container"
    value_only: false # <-------
    minimum: 1
    maximum: 100
    tags:
       a: b
       c: d
       kind: database

```
Sensor expects messages in this format (timestamp is in ISO 8601)
```json

{
   "name": "redis subscriber",
   "value": 47.1,
   "error": "",
   "timestamp": "2024-06-16T11:21:56.238Z"
}
```

```json
{
"name": "redis subscriber",
"value": 0,
"error": "something is broken",
"timestamp": "2024-06-16T11:21:56.238Z"
}

```



***Shell command sensor***

This sensor periodically executes shell command, expecting it to return measurable value
in STDOUT.


```yaml

  - name: doSomethingScriptSensor
    type: shell
    command: '/usr/bin/do_something.sh'
    description: "Execute script every 3th minute to measure something important"
    refresh_rate: 3m
    link: "http://example.org"
    environment:
       token: 2NpRTOwsEzseYUjVUVVfw
    tags:
      a: b
      c: d
```

If script outputs JSON to stdout, it can be parsed using [JSON Path](https://jsonpath.com/).
For example, script returns
```json
{
   "a": 5.3,
   "b": "something",
   "d": [
      10, 11, 24
   ]
}

```

This `$.a` JSONPath query will provide `5.3` - value of `a` key, and this one
`$.d[1]` will provide 11 - 2nd element of array under `d` key.
Parameters `a:10` and `b: 1` will make linear transformation of reading by
multiplying it by 10 and adding 1.

```yaml

  - name: doSomethingScriptSensor
    type: shell
    command: '/usr/bin/do_something_json.sh'
    description: "Execute script every 3th minute to measure something important"
    refresh_rate: 3m
    link: "http://example.org"
    environment:
      token: 2NpRTOwsEzseYUjVUVVfw
    json_path: $.a     
    a: 10
    b: 1
    tags:
       a: b
       c: d

```

***File sensor***

Sensor reads values from file, applying JSONPath query extraction if required

```yaml
  - name: thermal0
    type: file
    description: "Get thermal sensor status from area 0"
    link: "https://example.org"
    path_to_reading: /sys/class/thermal/thermal_zone0/temp
    a: 0.001
    b: 0
    minimum: 1
    maximum: 100
    refresh_rate: 5s
    tags:
      kind: thermal
```

***Endpoint sensor***

Waits for incoming HTTP POST request from external scripts/applications to update value.
Incoming HTTP request should have `Token: ....` with value equal to the one in config. 

Consider dashboard application is running on `localhost:3000`. For config like this

```yaml

- name: endpoint1
  type: endpoint
  description: "Update value by incoming POST request"
  token: "test321"
  
- name: endpoint2
  type: endpoint
  description: "Update value by incoming POST request"
  token: "test321forEndpoint2"
  

```

This curl command updates sensor `endpoint1` with value 21

```shell

curl -v -H "Host: localhost" \
  -H "Token: test321" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST \
  -d "name=endpoint1&value=21" \
  http://localhost:3000/update

```

and this one:

```shell

curl -v -H "Host: localhost" \
  -H "Token: test321forEndpoint2" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST \
  -d "name=endpoint2&value=53.5" \
  http://localhost:3000/update

```

updates `endpoint2` sensor with value 53.5

***CURL sensor***

This sensor sends periodical HTTP requests to external endpoint providing sensor readings in form of raw string or JSON data.

```yaml

- name: curl1
  type: curl
  description: "Sensor sends simple HTTP GET request expecting float string in response with latitude of IP address origin"
  http_method: "GET"
  link: "https://ip-api.com/"
  endpoint: "http://ip-api.com/line/193.41.76.51?fields=lat"
  
  
- name: curl2
  type: curl
  description: "Sensor sends simple HTTP GET request expecting JSON response"
  link: "https://ip-api.com/"
  http_method: "GET"
  endpoint: "http://ip-api.com/json/193.41.76.51"
  headers:
     User-Agent: "Vodolaz095's Dashboard"
  json_path: "@.lat"

- name: curl3
  type: curl
  description: "Sensor sends POST request expecting JSON response"
  http_method: "POST"
  endpoint: "https://example.org/api/v1/rpc"
  headers:
     User-Agent: "Vodolaz095's Dashboard"
     Authorization: "Bearer: EFLXCXxv7QCU7GyDvE36Azl8e8gIc0kG0BvGHNEnxAYA"
     Content-Type: "application/x-www-form-urlencoded"
  json_path: "@.balance"
  body: "entity=portfolio&action=get"


```

***Creating your own sensor***

Sensor should implement **ISensor** interface as provided [sensor.go](sensors%2Fsensor.go).
Sensor can be build on top of **UnimplementedSensor** with methods required impemented.
See examples in [sensors](sensors) directory.


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


