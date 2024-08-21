package middlewares

import (
	"net/http"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

// see https://github.com/gin-contrib/secure

//Strict-Transport-Security	HTTP Strict Transport Security is an excellent feature to support on your site and strengthens your implementation of TLS by getting the User Agent to enforce the use of HTTPS. Recommended value "Strict-Transport-Security: max-age=31536000; includeSubDomains".
//Content-Security-Policy	Content Security Policy is an effective measure to protect your site from XSS attacks. By whitelisting sources of approved content, you can prevent the browser from loading malicious assets.
//X-Frame-Options tells the browser whether you want to allow your site to be framed or not. By preventing a browser from framing your site you can defend against attacks like clickjacking. Recommended value "X-Frame-Options: SAMEORIGIN".
//X-Content-Type-Options	X-Content-Type-Options stops a browser from trying to MIME-sniff the content type and forces it to stick with the declared content-type. The only valid value for this header is "X-Content-Type-Options: nosniff".
//Referrer-Policy	Referrer Policy is a new header that allows a site to control how much information the browser includes with navigations away from a document and should be set by all sites.
//Permissions-Policy	Permissions Policy is a new header that allows a site to control which features and APIs can be used in the browser.

// Secure makes things more funny
func Secure(router *gin.Engine, domain string) {
	router.Use(secure.New(secure.Config{
		AllowedHosts:          []string{domain},
		SSLRedirect:           true,
		SSLTemporaryRedirect:  true,
		SSLHost:               domain,
		STSSeconds:            24 * 60 * 60,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'unsafe-inline' 'self'",
		ReferrerPolicy:        "no-referrer",
		IsDevelopment:         gin.IsDebugging(),
		BadHostHandler: func(c *gin.Context) {
			if c.ClientIP() == "127.0.0.1" {
				c.Next()
				return
			}
			c.String(http.StatusNotFound, "Domain %s is unknown", c.Request.Header.Get("HOST"))
			c.Abort()
		},
		IENoOpen:                  true,
		DontRedirectIPV4Hostnames: true,
		SSLProxyHeaders:           map[string]string{"X-Forwarded-Proto": "https"},
		FeaturePolicy:             "geolocation 'none'; microphone 'none'",
	}))
	router.Use(func(c *gin.Context) {
		c.Header("Permissions-Policy", "geolocation=(), microphone=()")
	})
}
