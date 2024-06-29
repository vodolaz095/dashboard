Deployment examples
===============================

Application can be deployed via various approaches. Since it can monitor and be dependent on different remote resources,
there is no `one approach fits all` way to deploy it. 

Systemd (system service) + NGINX
===============================
1. Limited system user `dashboard` with home directory in `/var/lib/dashboard` was created
2. [Systemd unit file](https://github.com/vodolaz095/dashboard/blob/master/contrib/systemd/dashboard.service) was tested on Fedora 40 server.
3. [NGINX site config](https://github.com/vodolaz095/dashboard/blob/master/contrib/nginx/dashboard.conf) was placed in `/etc/nginx/sites/`

Systemd (user service)
===============================
TODO

Docker Swarm
===============================
TODO

Kubernetes
===============================
TODO
