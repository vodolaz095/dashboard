Shell sensor
==============================

Sensor periodically calls user's defined command.
Command should provide reading in stdout either in format
of parsable by [strconv#ParseFloat](https://pkg.go.dev/strconv#ParseFloat) or as
a field in JSON object extractable by [JSON Path](https://jsonpath.com/) query.

Shared sensor parameters are explained in
[sensor_shared.md](https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_shared.md)
file.

All config parameters for sensors are depicted in this file
[sensor.go](https://github.com/vodolaz095/dashboard/blob/master/config/sensor.go)
with comments explaining things.



Configuration examples
==============================

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

  - name: doSomethingScriptSensorJson
    type: shell
    command: '/usr/bin/do_something_json.sh'
    description: "Execute script every 3th minute to measure something important"
    refresh_rate: 3m
    link: "http://example.org"
    environment:
      token: 2NpRTOwsEzseYUjVUVVfw
    json_path: $.a # jsonquery allows to extract parameters from JSON from script STDOUT    
    a: 10
    b: 1
    tags:
       a: b
       c: d

```

