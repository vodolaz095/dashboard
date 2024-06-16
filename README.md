dashboard
======================
Goland powered dashboard

Main features
======================
1. Few very hackable sensors
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

  - name: redis subscriber
    type: subscriber
    description: "Subscribe to redis channel and get values from it"
    channel: "vodolaz095/dashboard/subscriber/all"
    connection_name: "subscribe2redis@container"
    value_only: false
    minimum: 1
    maximum: 100
    tags:
       a: b
       c: d
       kind: database


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

And this `$.a` JSONPath query will provide 5.3 - value of `a` key, and this one
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

***Endpoint sensor***

***CURL sensor***


