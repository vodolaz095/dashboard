File sensor
=========================
Sensor periodically reads values from file, applying [JSONPath](https://jsonpath.com/) query extraction if required


Shared sensor parameters are explained in
[sensor_shared.md](https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_shared.md)
file.

All config parameters for sensors are depicted in this file
[sensor.go](https://github.com/vodolaz095/dashboard/blob/master/config/sensor.go)
with comments explaining things.

```yaml
sensors:
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

If file has JSON in it, it can be parsed using [JSON Path](https://jsonpath.com/).
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

sensor:
  - name: fileWithJSON
    type: file
    description: "Get something from big json file updated periodically"
    path_to_reading: /var/run/something.json
    json_path: $.something.value # extract value from field of JSON object...     
    refresh_rate: 5s

```
