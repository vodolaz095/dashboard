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
