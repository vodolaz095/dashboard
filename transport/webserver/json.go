package webserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/dashboard/model"
)

func (tr *Transport) exposeJSON() {
	tr.engine.GET("/api/v1/sensor/:name", func(c *gin.Context) {
		name := c.Param("name")
		sensor, found := tr.SensorsService.Sensors[name]
		if !found {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		var errMsg string
		if sensor.GetLastError() != nil {
			errMsg = sensor.GetLastError().Error()
		}
		c.JSON(http.StatusOK, model.Sensor{
			Name:        sensor.GetName(),
			Type:        sensor.GetType(),
			Description: sensor.GetDescription(),
			Link:        sensor.GetLink(),
			Minimum:     sensor.GetMinimum(),
			Maximum:     sensor.GetMaximum(),
			Value:       sensor.GetValue(),
			Error:       errMsg,
			Tags:        sensor.GetTags(),
			UpdatedAt:   sensor.GetUpdatedAt(),
		})
	})

	tr.engine.GET("/api/v1/sensor", func(c *gin.Context) {
		sensors, _ := tr.listFilteredSensors(c)
		c.JSON(http.StatusOK, sensors)
	})

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
