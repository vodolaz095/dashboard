package webserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (tr *Transport) exposeIndex() {
	tr.engine.GET("/", func(c *gin.Context) {
		sensors, filtered := tr.listFilteredSensors(c)
		stats := tr.SensorsService.Stats()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":               tr.Title,
			"description":         tr.Description,
			"keywords":            strings.Join(tr.Keywords, ", "),
			"doIndex":             tr.DoIndex,
			"sensors":             sensors,
			"filtered":            filtered,
			"version":             tr.Version,
			"header":              string(tr.header),
			"footer":              string(tr.footer),
			"sensors_updated_now": stats.QueueLength,
			"queue_length":        stats.SensorsUpdatedNow,
			"subscribers":         stats.Subscribers,
		})
	})
}
