Redis Synchronous Query Sensor
=========================================
This sensor periodically executes redis command and records value returned.
It is possible to run LUA stored procedures in redis and get their values as sensor readings

Shared sensor parameters are explained in
[sensor_shared.md](https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_shared.md)
file.

All config parameters for sensors are depicted in this file
[sensor.go](https://github.com/vodolaz095/dashboard/blob/master/config/sensor.go)
with comments explaining things.


For example, this configuration makes sensor to read value of key `a` in redis database, which is expected to have 
meaningful value which can be parsed by [strconv#ParseFloat](https://pkg.go.dev/strconv#ParseFloat)

```yaml
# https://github.com/vodolaz095/dashboard/blob/master/docs/connection_pool.md
database_connections:
  - name: redis@container
    type: redis
    connection_string: "redis://127.0.0.1:6379"
    max_open_cons: 3
    max_idle_cons: 1
    
sensor:
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

Redis Subscriber Sensor
=========================================

This sensor subscribes to redis database channels and reads updates provided as float numbers.
Due to redis limitations, connection can only in one of modes - sync commands or subscription.

So, **sensors of type `redis` and `subscriber` cannot share redis connections.** 


```yaml
database_connections:
  - name: redis@container
    type: redis
    connection_string: "redis://default:secret@127.0.0.1:6379"
    max_open_cons: 3
    max_idle_cons: 1
  - name: subscribe2redis@container
    type: redis
    connection_string: "redis://default:secret@127.0.0.1:6379"
    max_open_cons: 2
    max_idle_cons: 1
    
sensor:
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
for successful updates
```json

{
   "name": "redis subscriber",
   "value": 47.1,
   "error": "",
   "timestamp": "2024-06-16T11:21:56.238Z"
}
```
and this one for erroneous updates
```json
{
"name": "redis subscriber",
"value": 0,
"error": "something is broken",
"timestamp": "2024-06-16T11:21:56.238Z"
}

```
