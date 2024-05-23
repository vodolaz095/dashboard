package webserver

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (tr *Transport) exposeFeed() {
	tr.engine.GET("/feed", func(c *gin.Context) {
		log.Info().Msgf("Client %s subscribes to feed", c.Request.RemoteAddr)
		defer log.Info().Msgf("Client %s finished subscription to feed", c.Request.RemoteAddr)
		events, err := tr.SensorsService.Subscribe(c.Request.Context(), c.Request.RemoteAddr)
		if err != nil {
			c.String(http.StatusBadRequest, "error subscribing: %s", err)
			return
		}
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-events; ok {
				log.Debug().Msgf("Broadcasting %s: %v", msg.Name, msg.Value)
				c.SSEvent(msg.Name, msg)
				return true
			}
			return false
		})
	})
}
