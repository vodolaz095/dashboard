Linux System Sensor
=======================================

Sensor periodically reads linux system parameters like load average, number of currently running processes and free ram.

Shared sensor parameters are explained in
[sensor_shared.md](https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_shared.md)
file.

All config parameters for sensors are depicted in this file
[sensor.go](https://github.com/vodolaz095/dashboard/blob/master/config/sensor.go)
with comments explaining things.


```yaml

sensors:
  - name: load1
    type: load1
    description: "Get system load average during last minute"
    refresh_rate: 5s
    tags:
      kind: load
      
  - name: load5
    type: load5
    description: "Get system load average during last 5 minutes"
    refresh_rate: 5s
    tags:
      kind: load
      
  - name: load15
    type: load15
    description: "Get system load average during last 15 minutes"
    refresh_rate: 5s
    tags:
      kind: load
      
  - name: process
    type: process
    description: "Number of currently running processes"
    refresh_rate: 5s
    tags:
      kind: load
      
  - name: free_ram
    type: free_ram
    description: "Current free Random Access Memory volume in megabytes"
    refresh_rate: 5s
    link: "https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_linux_system.md"
    minimum: 500
    maximum: 8000
    tags:
      kind: ram
```
