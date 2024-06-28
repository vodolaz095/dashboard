***CURL sensor***

This sensor sends periodical HTTP requests to external endpoint providing sensor readings in form of raw string or JSON data.

```yaml

- name: curl1
  type: curl
  description: "Sensor sends simple HTTP GET request expecting float string in response with latitude of IP address origin"
  http_method: "GET"
  link: "https://ip-api.com/"
  endpoint: "http://ip-api.com/line/193.41.76.51?fields=lat"
  
  
- name: curl2
  type: curl
  description: "Sensor sends simple HTTP GET request expecting JSON response"
  link: "https://ip-api.com/"
  http_method: "GET"
  endpoint: "http://ip-api.com/json/193.41.76.51"
  headers:
     User-Agent: "Vodolaz095's Dashboard"
  json_path: "@.lat"

- name: curl3
  type: curl
  description: "Sensor sends POST request expecting JSON response"
  http_method: "POST"
  endpoint: "https://example.org/api/v1/rpc"
  headers:
     User-Agent: "Vodolaz095's Dashboard"
     Authorization: "Bearer: EFLXCXxv7QCU7GyDvE36Azl8e8gIc0kG0BvGHNEnxAYA"
     Content-Type: "application/x-www-form-urlencoded"
  json_path: "@.balance"
  body: "entity=portfolio&action=get"


```
