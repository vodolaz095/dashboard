package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/internal/sensors/endpoint"
	"github.com/vodolaz095/dashboard/internal/transport/webserver/dto"
)

// exposeUpdate exposes endpoint used to update sensor value by incoming HTTP POST request
func (tr *Transport) exposeUpdate() {
	handle := func(c *gin.Context) {
		var data dto.UpdateSensorRequest
		err := c.Bind(&data)
		if err != nil {
			c.String(http.StatusBadRequest, "Malformed request body: %s", err)
			return
		}
		sensor, found := tr.SensorsService.Sensors[data.Name]
		if !found {
			log.Info().
				Str("sensor", data.Name).
				Float64("reading", data.Value).
				Str("User-Agent", c.GetHeader("User-Agent")).
				Str("Client-IP", c.ClientIP()).
				Msgf("Sensor %s is not found", data.Name)
			c.String(http.StatusBadRequest, "sensor %s not found", data.Name)
			return
		}
		casted, ok := sensor.(*endpoint.Sensor)
		if !ok {
			log.Info().
				Str("sensor", data.Name).
				Float64("reading", data.Value).
				Str("User-Agent", c.GetHeader("User-Agent")).
				Str("Client-IP", c.ClientIP()).
				Msgf("Sensor %s's type is not `endpoint`", data.Name)
			c.String(http.StatusBadRequest, "sensor %s type is wrong", data.Name)
			return
		}
		if casted.Token != "" {
			if casted.Token != c.GetHeader("Token") {
				log.Warn().
					Str("sensor", casted.Name).
					Float64("reading", casted.GetValue()).
					Str("User-Agent", c.GetHeader("User-Agent")).
					Str("Client-IP", c.ClientIP()).
					Msgf("Updating endpoint sensor %s with value %v failed for token mismatch",
						casted.Name, data.Value,
					)
				c.String(http.StatusBadRequest, "Header `Token` has wrong value")
				return
			}
		}
		log.Info().
			Str("sensor", casted.Name).
			Float64("reading", casted.GetValue()).
			Str("User-Agent", c.GetHeader("User-Agent")).
			Str("Client-IP", c.ClientIP()).
			Msgf("Updating endpoint sensor %s with value %v",
				casted.Name, data.Value,
			)
		switch c.FullPath() {
		case "/update":
			casted.Set(data.Value)
			break
		case "/increment":
			casted.Increment(data.Value)
			break
		case "/decrement":
			casted.Increment(-data.Value)
			break
		default:
			c.String(http.StatusBadRequest, "unknown method")
		}
		if data.Description != "" {
			casted.SetDescription(data.Description)
		}
		tr.SensorsService.Broadcast(casted.Name, "", casted.GetStatus(), casted.GetValue())
		c.AbortWithStatus(http.StatusNoContent)
	}

	tr.engine.POST("/update", handle)
	tr.engine.POST("/increment", handle)
	tr.engine.POST("/decrement", handle)
}
