package webserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (tr *Transport) exposeJSON() {
	tr.engine.GET("/json", func(c *gin.Context) {
		sensors, filtered := tr.listFilteredSensors(c)
		c.JSON(http.StatusOK, gin.H{
			"title":       tr.Title,
			"description": tr.Description,
			"keywords":    strings.Join(tr.Keywords, ", "),
			"sensors":     sensors,
			"filtered":    filtered,
			"version":     tr.Version,
			"stats":       tr.SensorsService.Stats(),
		})
	})
}
