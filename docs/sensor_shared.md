Shared sensor parameters
===============================

All config parameters for sensors are depicted in this file 
[sensor.go](https://github.com/vodolaz095/dashboard/blob/master/config/sensor.go)
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
`5h, 1m, 10s` are OK. Even `5m 2s` is ok too!

Parameters `minimum` and `maximum` does nothing, but they can give a hint for operator to understand save range of values.
