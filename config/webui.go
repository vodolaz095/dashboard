package config

// WebUI allows to customize webserver - where it listens, what title does index page have and so on
type WebUI struct {
	// Listen sets address, where application is listening, for example, 127.0.0.1:3000
	Listen string `yaml:"listen" validate:"required,hostname_port"`
	// Domain sets HTTP HOST where application accepts requests
	Domain string `yaml:"domain" validate:"hostname_rfc1123"`
	// Title sets title of index page
	Title string `yaml:"title"`
	// Description sets description of index page
	Description string `yaml:"description"`
	// Keywords sets keywords of index page
	Keywords []string `yaml:"keywords"`
	// DoIndex sets http header equivalents to allow page indexing by search engine crawlers
	DoIndex bool `yaml:"do_index"`
	// PathToHeader contains path to file for header which will be included in dashboard template as header
	PathToHeader string `yaml:"path_to_header"`
	// PathToFooter contains path to file for footer which will be included in dashboard template as footer
	PathToFooter string `yaml:"path_to_footer"`
	// HeaderForClientIP defines header name used to extract clients IP address. For example,
	// for Google App Engine it is "X-Appengine-Remote-Addr", for Cloudflare - "CF-Connecting-IP" and for
	// Fly.io - "Fly-Client-IP"
	HeaderForClientIP string `yaml:"header_for_client_ip"`
	// TrustProxies - contains list of IP addresses of proxies to trust. See
	// https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies
	TrustProxies []string `yaml:"trust_proxies"`
}
