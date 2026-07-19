# AGENTS.md

## Overview

This document provides essential information for agents working in the Vodolaz095's Dashboard repository. It covers the project's architecture, key components, commands, and non-obvious patterns to help agents work effectively without trial-and-error discovery.

## Project Type and Architecture

Vodolaz095's Dashboard is a minimalistic, DDoS-proof Golang-powered dashboard application designed to display real-time sensor readings from various sources. The architecture follows a producer-consumer pattern with in-memory storage to prevent database DDoS attacks.

### Key Architectural Principles

1. **DDoS Protection**: Sensor values are updated in background goroutines and served from memory by the HTTP server. Multiple clients accessing the dashboard do not generate additional database queries.
2. **Real-time Updates**: Uses Server-Sent Events (SSE) for real-time updates to clients.
3. **Lightweight Design**: The dashboard is designed to work well on older devices (iPhone 6, 2013 Android smartphones) with minimal assets (~1KB CSS, ~1KB JavaScript, ~5KB main page).
4. **Configuration-Driven**: Most functionality is configured via a YAML file rather than code changes.

### Component Structure

The application is organized into several key packages:

- `internal/service`: Core business logic for sensors, including the `SensorsService` that manages sensor lifecycle and updates
- `internal/transport`: HTTP server implementation using Gin, SSE endpoints, and data export formats
- `internal/sensors`: Various sensor implementations that collect data from different sources
- `internal/views`: HTML templates for the dashboard UI
- `config`: Configuration structures and loading
- `model`: Data transfer objects

## Essential Commands

### Development Commands

```bash
# Set up development environment
cd ~/projects/dashboard
make tools  # Verify required tools are installed
make deps  # Download Go modules

make docker/resource  # Start only development databases using Docker
make podman/resource  # Start only development databases using Podman

# Start development server with example config
make start
# or
go run main.go ./contrib/dashboard.yaml

# Build production binary
make build
# or
go build -o build/dashboard main.go

# Run tests
make test
# or
go test -v ./...

# Run linters
make lint

```

### Container-Based Development

```bash
# Using Docker
make docker/up  # Start development containers and dashboard
make docker/resource  # Start only development databases
make docker/down  # Stop all containers

# Using Podman
make podman/up  # Start development containers and dashboard
make podman/resource  # Start only development databases
make podman/down  # Stop all containers
```

### Build Commands

```bash
# Build for current platform
make build

# Build for specific platforms
make build/linux_amd64
make build/linux_arm6
make build/linux_arm7
make build/windows
make build/macos

# Build all platforms
make build/all
```

## Code Organization and Structure

### Sensor Architecture

Sensors are the core data collection components. They follow a consistent pattern:

1. Implement the `ISensor` interface from `internal/sensors/sensor.go`
2. Embed `UnimplementedSensor` to inherit common functionality
3. Implement `Init`, `Ping`, `Close`, and `Update` methods

Available sensor types include:
- `load1`, `load5`, `load15`: System load averages
- `process`: Number of running processes
- `free_ram`: Free RAM in MB
- `free_disk_space`, `used_disk_space`, `free_disk_space_ratio`: Disk space metrics
- `shell`: Execute shell commands/scripts
- `endpoint`: Receive updates via incoming HTTP POST requests
- `redis`: Execute Redis commands
- `subscriber`: Subscribe to Redis channels
- `mysql`, `postgres`: Execute SQL queries against connected databases
- `curl`: Make HTTP requests to external endpoints
- `file`: Read values from files
- `victoria_metrics`: Query Victoria Metrics

### Configuration Structure

The configuration follows a hierarchical structure defined in `config/config.go`:

- `web_ui`: HTTP server configuration (address, domain, title, etc.)
- `log`: Logging configuration (level, journald integration)
- `database_connections`: Reusable database connections
- `sensors`: Sensor definitions
- `broadcasters`: Redis pub/sub configurations for data export
- `influx`: InfluxDB configuration for time-series data export

### Key Non-Obvious Patterns

1. **Sensor Name Uniqueness**: Sensor names must be unique. Duplicate sensor names will cause the application to fail with a fatal error.

2. **Redis Connection Limitation**: Redis connections cannot be shared between regular sensors and subscriber sensors. Subscriber sensors require dedicated connections because Redis connections are locked in subscription mode.

3. **Linear Transformations**: Sensors support linear transformations via `A` and `B` parameters (f(x) = A*x + B), useful for unit conversions like Fahrenheit to Celsius.

4. **Tag-Based Filtering**: The dashboard supports URL parameter filtering based on sensor tags. For example, `?kind=database&unit=sales` will show only sensors with both tags.

5. **Endpoint Sensor Security**: Endpoint sensors use token-based authentication. The token is specified in the sensor configuration and must be provided in the `Token` header of update requests.


## Testing Approach

The project uses standard Go testing with the `testing` package and `testify` for assertions. Key testing patterns:

1. **Sensor Testing**: Most sensors have corresponding test files that use the `DoTestSensor` helper function from `internal/sensors/sensor.go`.
2. **Integration Testing**: The `.github/workflows/test.yml` workflow sets up Redis, MariaDB, and PostgreSQL containers for integration testing.
3. **Test Structure**: Tests typically follow the pattern of initializing the sensor, pinging it, updating it, and verifying the value.

Example test command:
```bash
make test
# or
go test -v ./...
```

## Important Gotchas and Non-Obvious Patterns

### Redis Connection Gotchas

1. **Separate Connections Required**: Regular Redis sensors (type `redis`) and subscriber sensors (type `subscriber`) cannot share connections. Each subscriber sensor needs its own dedicated Redis connection.

2. **Connection String Format**: Redis connection strings must follow the format `redis://[user:password@]host:port/db` as parsed by `redis.ParseURL()`.

### Sensor Configuration Gotchas

1. **Refresh Rate Parsing**: The `refresh_rate` parameter accepts strings that `time.ParseDuration()` understands, including compound durations like `5m 2s`.

2. **JSONPath Queries**: When using `json_path` with shell, file, or curl sensors, ensure the JSONPath query is valid. Invalid queries will cause sensor updates to fail.

3. **File Sensor Paths**: The `path_to_reading` parameter for file sensors must be accessible by the running process. For system files like `/sys/class/thermal/thermal_zone0/temp`, appropriate permissions are required.

### Security Considerations

1. **Sensitive Data Concealment**: Database credentials, tokens, and passwords are concealed from dashboard visitors. Only sensor values and metadata are displayed.

2. **WebUI Access Control**: Dashboard access should be restricted either by reverse proxy or by serving only in a local network.

3. **Endpoint Sensor Tokens**: When using endpoint sensors, ensure the token is sufficiently complex to prevent unauthorized updates.

### Performance Considerations

1. **Deferred Queue**: The application uses a deferred queue (`dqueue.Handler`) to manage sensor updates, preventing overwhelming the system with concurrent updates.

2. **Connection Pooling**: Database connections are pooled and reused across sensors of the same type.

3. **Memory Usage**: All sensor values are stored in memory, so the number of sensors should be balanced against available system memory.

## Deployment Examples

### Systemd Service

The project includes an example systemd service file in `contrib/systemd/dashboard.service` and NGINX configuration in `contrib/nginx/dashboard.conf`.

### Docker Deployment

The project provides Dockerfiles for both production (`Dockerfile`) and development (`Dockerfile_development`) environments.

### Podman Deployment

Makefile targets are provided for Podman deployment (`make podman/up`, `make podman/down`, etc.).

## Configuration Examples

### Basic Sensor Configuration

```yaml
sensors:
  - name: load1
    type: load1
    description: "Get system load average during last minute"
    refresh_rate: 5s
    tags:
      kind: load

  - name: free_ram
    type: free_ram
    description: "Current free RAM in megabytes"
    refresh_rate: 5s
    minimum: 500
    maximum: 8000
```

### Database Sensor Configuration

```yaml
# Database connections
database_connections:
  - name: mysql@container
    type: mysql
    connection_string: "root:dashboard@tcp(127.0.0.1:3306)/dashboard"
    max_open_cons: 2
    max_idle_cons: 1

  - name: postgres@container
    type: postgres
    connection_string: "postgres://dashboard:dashboard@127.0.0.1:5432/dashboard"
    max_open_cons: 2
    max_idle_cons: 1

# Sensors
sensors:
  - name: mysql_random
    type: mysql
    description: "Select random number from range"
    connection_name: "mysql@container"
    query: "SELECT rand()*99+1 as random"
    minimum: 30
    maximum: 60
    refresh_rate: 5s

  - name: postgres_random
    type: postgres
    description: "Select random number from range"
    connection_name: "postgres@container"
    query: "SELECT random()*99+1 as random"
    minimum: 30
    maximum: 60
    refresh_rate: 5s
```

### Redis Sensor Configuration

```yaml
# Database connections
database_connections:
  - name: redis@container
    type: redis
    connection_string: "redis://127.0.0.1:6379"
    max_open_cons: 2
    max_idle_cons: 1

# Sensors
sensors:
  - name: redis_key_value
    type: redis
    description: "Get value of redis key a"
    connection_name: "redis@container"
    query: "get a"
    refresh_rate: 5s
```

### External HTTP Sensor Configuration

```yaml
sensors:
  - name: curl_sensor
    type: curl
    description: "Fetch data from external API"
    http_method: "GET"
    endpoint: "https://api.example.com/data"
    json_path: "$.value"
    refresh_rate: 30s
```

## Additional Resources

- **Full Configuration Example**: `contrib/dashboard.yaml`
- **Development Configuration**: `contrib/dashboard_docker.yaml`
- **Sensor Documentation**: Files in `docs/` directory, particularly `sensor_shared.md`
- **Deployment Examples**: `contrib/systemd/` and `contrib/nginx/` directories
- **UI Customization**: `contrib/header.html` and `contrib/footer.html` for adding custom HTML to the dashboard

This document should be updated when significant changes are made to the project structure, new sensor types are added, or important patterns emerge during development.
