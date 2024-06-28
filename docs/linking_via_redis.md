Linking few dashboard instances via redis Pub/Sub
==============================
Let's consider case when you have few servers running dashboards, and you want to have
one dashboard showing values from all other dashboards.
It can be done via `broadcaster` feature, utilizing shared redis server pub/sub channels.

For example, let the `server1` has his configuration and very important script
`/usr/bin/do_something.sh` which will be run every minute to extract readings.

```yaml

web_ui:
  listen: "0.0.0.0:3000"
  domain: "server1.example.org"
  title: "dashboard"
  do_index: false

log:
  level: trace
  to_journald: false

database_connections:
  - name: redis@container
    type: redis
    connection_string: "redis://username@password:redis.example.org:6379"

sensors:
  - name: sense_something
    type: shell
    command: '/usr/bin/do_something.sh'
    description: "Execute script"
    refresh_rate: 1m

broadcasters:
  - connection_name: redis@container
    subject: "vodolaz095/dashboard/sensor/%s"
    value_only: false
  - connection_name: redis@container
    subject: "vodolaz095/dashboard/value/%s"
    value_only: true
```

By using broadcaster feature, we can send sensor readings as redis publishing into 2 channels.
Notice `subject` field - it is template with name of sensor in it.
So, sensor `sense_something` will publish it readings into  redis with topics
- vodolaz095/dashboard/sensor/**sense_something** with format



```yaml

web_ui:
  listen: "0.0.0.0:3000"
  domain: "server2.example.org"
  title: "dashboard"
  do_index: false

log:
  level: trace
  to_journald: false

database_connections:
  - name: redis@container
    type: redis
    connection_string: "redis://username@password:redis.example.org:6379"

sensors:
  - name: sense_something_from_server1
    type: subscriber
    description: "Subscribe to redis channel and get values broadcast from server1"
    channel: "vodolaz095/dashboard/sensor/sense_something"
    connection_name: "redis@container"
    value_only: false
    minimum: 1
    maximum: 100
    tags:
       a: b
       c: d
       kind: database


```
