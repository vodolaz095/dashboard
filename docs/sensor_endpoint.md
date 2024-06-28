***Endpoint sensor***

Waits for incoming HTTP POST request from external scripts/applications to update value.
Incoming HTTP request should have `Token: ....` with value equal to the one in config.

Consider dashboard application is running on `localhost:3000`. For config like this

```yaml

- name: endpoint1
  type: endpoint
  description: "Update value by incoming POST request"
  token: "test321"
  
- name: endpoint2
  type: endpoint
  description: "Update value by incoming POST request"
  token: "test321forEndpoint2"
  

```

This curl command updates sensor `endpoint1` with value 21

```shell

curl -v -H "Host: localhost" \
  -H "Token: test321" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST \
  -d "name=endpoint1&value=21" \
  http://localhost:3000/update

```

and this one:

```shell

curl -v -H "Host: localhost" \
  -H "Token: test321forEndpoint2" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST \
  -d "name=endpoint2&value=53.5" \
  http://localhost:3000/update

```

updates `endpoint2` sensor with value 53.5
It is possible to increment/decrement in a race condition save manner values by calling `/increment` and `/decrement`
endpoints:
```shell

# set sensor `endpoint1` value to 10
curl -v -H "Host: localhost" \
  -H "Token: test321" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST \
  -d "name=endpoint1&value=10" \
  http://localhost:3000/update

# increment `endpoint1` value by 5. If there are few concurrent HTTP POST requests, all data will be applied in save manner
curl -v -H "Host: localhost" \
  -H "Token: test321" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST \
  -d "name=endpoint1&value=5" \
  http://localhost:3000/increment

# decrement `endpoint1` value by 3. If there are few concurrent HTTP POST requests, all data will be applied in save manner
curl -v -H "Host: localhost" \
  -H "Token: test321" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST \
  -d "name=endpoint1&value=3" \
  http://localhost:3000/decrement

```
