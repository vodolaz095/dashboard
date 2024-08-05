Deployment examples
===============================

Application can be deployed via various approaches. Since it can monitor and be dependent on different remote resources,
there is no `one approach fits all` way to deploy it. 

Systemd (system service) + NGINX
===============================
NGINX as reverse proxy, encryption and authorization is done by NGINX.

1. Limited system user `dashboard` with home directory in `/var/lib/dashboard` was created
2. [Systemd unit file](https://github.com/vodolaz095/dashboard/blob/master/contrib/systemd/dashboard.service) was tested on Fedora 40 server.
3. [NGINX site config](https://github.com/vodolaz095/dashboard/blob/master/contrib/nginx/dashboard.conf) was placed in `/etc/nginx/sites/`

Good read:
- https://nginx.org/ru/docs/http/ngx_http_proxy_module.html
- https://stackoverflow.com/questions/23844761/upstream-sent-too-big-header-while-reading-response-header-from-upstream
- https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-http-basic-authentication/


Systemd (user service)
===============================
TODO

Docker Swarm
===============================
TODO

Kubernetes
===============================
TODO

Configuring HTTP server
===============================

In general, it can be a good idea to make dashboard listen on separate network interface by using 
`listen` configuration parameter in `web_ui` part of config and make it respond to
HTTP requests having `HOST header` matching `domain` configuration parameter.
If you use Cloudflare and other CDN you can define HTTP header name used to extract clients IP address by setting 
parameter `header_for_client_ip` as explained in GIN documentation - https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies

```yaml

web_ui:
  listen: "192.168.3.5:3000"
  domain: "dashboard.local"
  title: "dashboard"
  description: "dashboard"
  header_for_client_ip: "CF-Connecting-IP"
  keywords:
    - "dashboard"
    - "vodolaz095"
    - "golang"
    - "redis"
    - "postgresql"
    - "mysql"
  do_index: true
  path_to_header: ./contrib/header.html
  path_to_footer: ./contrib/footer.html

```
