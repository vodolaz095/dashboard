package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example

func (tr *Transport) exposeMetrics() {
	tr.engine.GET("/metrics", func(c *gin.Context) {
		sensors, _ := tr.listFilteredSensors(c)
		c.Header("Content-Type", "text/plain; version=0.0.4")
		stats := tr.SensorsService.Stats()
		for i := range sensors {
			fmt.Fprint(c.Writer, sensors[i].String())
		}
		fmt.Fprintln(c.Writer, "# HELP dashboard_queue_length Number of items in deferred queue for sensors updates")
		fmt.Fprintln(c.Writer, "# TYPE dashboard_queue_length gauge")
		fmt.Fprintf(c.Writer, "dashboard_queue_length %v %v\n", stats.QueueLength, time.Now().Unix())

		fmt.Fprintln(c.Writer, "# HELP dashboard_subscribers Number of subscribers for event feeds")
		fmt.Fprintln(c.Writer, "# TYPE dashboard_subscribers gauge")
		fmt.Fprintf(c.Writer, "dashboard_subscribers %v %v\n", stats.Subscribers, time.Now().Unix())

		fmt.Fprintln(c.Writer, "# HELP dashboard_sensors_updated Number of sensors being updated now")
		fmt.Fprintln(c.Writer, "# TYPE dashboard_sensors_updated gauge")
		fmt.Fprintf(c.Writer, "dashboard_sensors_updated %v %v\n", stats.SensorsUpdatedNow, time.Now().Unix())
		c.AbortWithStatus(http.StatusOK)
	})
}
