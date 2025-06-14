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
part of configuration. Each SQL database connection can be reused by multiple sensors.

MySQL compatible database connection strings should satisfy this data source name syntax:
```
username:password@tcp(hostname:3306)/database_name?charset=utf8&parseTime=True&loc=Local
username:password@unix(/var/lib/mysql/mysql.sock)/database_name?charset=utf8&parseTime=True&loc=Local

```
See https://pkg.go.dev/github.com/go-sql-driver/mysql#readme-dsn-data-source-name

PostgresSQL compatible database connection strings should be like this:
```

# Example Keyword/Value
user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10

# Example URL
postgres://username:password@hostname:5432/database_name
postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10

```
See https://pkg.go.dev/github.com/jackc/pgx/v5@v5.6.0/pgconn#ParseConfig



Configuration examples
=======================================

```yaml

# https://github.com/vodolaz095/dashboard/blob/master/docs/connection_pool.md
database_connections:
  - name: mysql@container
    type: mysql
    connection_string: "root:dashboard@tcp(127.0.0.1:3306)/dashboard"
    max_open_cons: 3
    max_idle_cons: 1
    
  - name: postgres@container
    type: postgres
    connection_string: "postgres://dashboard:dashboard@127.0.0.1:5432/dashboard"
    max_open_cons: 3
    max_idle_cons: 1
    
sensors:
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
    connection_name: "mysql@container"
    # multiline queries are supported!
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
