package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/dashboard/sensors/endpoint"
	"github.com/vodolaz095/dashboard/transport/webserver/dto"
)

// exposeEndpoint exposes endpoint used to update sensor value by incoming HTTP POST request
func (tr *Transport) exposeEndpoint() {
	tr.engine.POST("/update", func(c *gin.Context) {
		var data dto.UpdateSensorRequest
		err := c.Bind(&data)
		if err != nil {
			c.String(http.StatusBadRequest, "Malformed request body: %s", err)
			return
		}
		sensor, found := tr.SensorsService.Sensors[data.Name]
		if !found {
			c.String(http.StatusBadRequest, "sensor %s not found", data.Name)
			return
		}
		casted, ok := sensor.(*endpoint.Sensor)
		if !ok {
			c.String(http.StatusBadRequest, "sensor %s type is wrong", data.Name)
			return
		}
		if casted.Token != "" {
			if casted.Token != c.GetHeader("Token") {
				c.String(http.StatusBadRequest, "Header `Token` has wrong value")
				return
			}
		}
		err = casted.Update(c.Request.Context(), data.Value)
		if err != nil {
			c.String(http.StatusBadRequest, "Error updating sensor: %s", err)
			return
		}
		c.AbortWithStatus(http.StatusNoContent)
	})
}