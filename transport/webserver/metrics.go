package webserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example

func (tr *Transport) exposeMetrics() {
	tr.engine.GET("/metrics", func(c *gin.Context) {
		sensors, _ := tr.listFilteredSensors(c)
		c.Header("Content-Type", "text/plain; version=0.0.4")
		for i := range sensors {
			fmt.Fprint(c.Writer, sensors[i].String())
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
}
