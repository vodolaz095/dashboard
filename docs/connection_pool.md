Connection pool
=============================

Each database connection can be reused by few sensors, so, connections are defined separately in config.
This is example how to define redis, mysql and postgresql connections

```yaml

database_connections:
# this connection can be used by few sensors and broadcasters, since `publish` command does not lock redis connection
  - name: redis@container
    type: redis
    connection_string: "redis://127.0.0.1:6379"
# since `subscribe/psubscribe` command locks redis connection, it can only be used by redis subscriber sensors.
  - name: subscribe2redis@container
    type: redis
    connection_string: "redis://127.0.0.1:6379"

# single SQL database connection can be shared between few sensors for this database
  - name: mysql@container
    type: mysql
    connection_string: "root:dashboard@tcp(127.0.0.1:3306)/dashboard"

  - name: postgres@container
    type: postgres
    connection_string: "postgres://dashboard:dashboard@127.0.0.1:5432/dashboard"


```

Redis database connections strings should be understood by [ParseURL](https://pkg.go.dev/github.com/redis/go-redis/v9#ParseURL)
```
redis://<user>:<password>@<host>:<port>/<db_number>
unix://<user>:<password>@</path/to/redis.sock>?db=<db_number>
```
Important: if you have `subscriber` type sensor, it should use separate redis connections, because
redis connection can work only in one of 2 modes - accepting commands, or being subscribed to channels.



If you connect to redis databases of old version (5.x and lower), you can omit `user` -
this should work `redis://:passwd@redis.example.org:6379/1`

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

