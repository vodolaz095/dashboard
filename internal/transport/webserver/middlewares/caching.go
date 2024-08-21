package middlewares

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control

// UseCaching exposes cache headers to make things faster
func UseCaching(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		switch {
		case c.FullPath() == "/favicon.ico",
			c.FullPath() == "/robots.txt",
			strings.HasPrefix(c.FullPath(), "/assets"):
			c.Header("Cache-Control", "public,max-age=300,s-maxage=900")

		default:
			c.Header("Cache-Control", "no-store")
			c.Header("expires", time.Now().Add(-time.Minute).Format(time.RFC1123))
			c.Header("pragma", "no-cache")
		}
		c.Next()
	})
}
