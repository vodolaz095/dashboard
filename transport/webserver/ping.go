package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (tr *Transport) exposeHealthcheck() {
	tr.engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong!")
	})
	tr.engine.GET("/healthcheck", func(c *gin.Context) {
		err := tr.SensorsService.Ping(c.Request.Context())
		if err != nil {
			log.Err(err).Msgf("Healthcheck failed with error: %s", err)
			c.String(http.StatusInternalServerError, "System malfunction confirmed!")
			return
		}
		c.String(http.StatusOK, "System integrity confirmed!")
	})
}
