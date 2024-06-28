MySQL/PostgreSQL sensor
===================================

This sensor queries database with user defined SQL request periodically.
SQL request **SHOULD** return **SINGLE** value parsable as **FLOAT64**.
Background process periodically calls query and cache result in memory, and webserver
process serves data from memory, so database cannot be DDoSed.


Shared sensor parameters are explained in
[sensor_shared.md](https://github.com/vodolaz095/dashboard/blob/master/docs/sensor_shared.md)
file.

All config parameters for sensors are depicted in this file
[sensor.go](https://github.com/vodolaz095/dashboard/blob/master/config/sensor.go)
with comments explaining things.


Connecting to database
=======================================

Database connections should be defined in [Connection Pool](https://github.com/vodolaz095/dashboard/blob/master/docs/connection_pool.md)
part of configuration.


MySQL database connection strings should satisfy this data source name syntax:
```
username:password@tcp(hostname:3306)/database_name
```
See https://pkg.go.dev/github.com/go-sql-driver/mysql#readme-dsn-data-source-name

PostgreSQL database connection strings should be like this:
```
postgres://username:password@hostname:5432/database_name
```
See https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#ParseConfig



Configuration examples
=======================================

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
