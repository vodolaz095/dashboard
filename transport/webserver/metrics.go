package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example

func (tr *Transport) exposeMetrics() {
	tr.engine.GET("/metrics", func(c *gin.Context) {
		sensors := tr.SensorsService.List()
		c.Header("Content-Type", "text/plain; version=0.0.4")
		c.HTML(http.StatusOK, "metrics.html", gin.H{
			"sensors": sensors,
		})
	})
}
