Dashboard customization
=============================
WebUI can be customized by setting page title, description, keywords and adding static HTML
code snippets to page with, for example, custom menu and tracking pixels for analytics platforms.
Just create/edit `header.html` and `footer.html` files with your HTML code in it and mention them in configuration file.

```yaml

web_ui:
  listen: "0.0.0.0:3000"
  domain: "localhost"
  title: "dashboard"
  description: "dashboard"
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
