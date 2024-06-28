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
