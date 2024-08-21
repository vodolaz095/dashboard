package webserver

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (tr *Transport) exposeFeed() {
	tr.engine.GET("/feed", func(c *gin.Context) {
		clientIP := c.ClientIP() // since c.ClientIP has some dramatic computations inside
		log.Info().Msgf("Client %s via %s subscribes to feed", clientIP, c.Request.RemoteAddr)
		defer log.Info().Msgf("Client %s via %s finished subscription to feed", clientIP, c.Request.RemoteAddr)
		events, err := tr.SensorsService.Subscribe(c.Request.Context(), c.Request.RemoteAddr)
		if err != nil {
			c.String(http.StatusBadRequest, "error subscribing: %s", err)
			return
		}
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		c.Stream(func(w io.Writer) bool {
			if msg, ok := <-events; ok {
				if msg.Error != "" {
					log.Trace().
						Str("sensor", msg.Name).
						Str("Client-IP", clientIP).
						Str("User-Agent", c.GetHeader("User-Agent")).
						Float64("value", msg.Value).
						Err(errors.New(msg.Error)).
						Msgf("Broadcasting %s: %v with error %s", msg.Name, msg.Value, msg.Error)
				} else {
					log.Trace().
						Str("Client-IP", clientIP).
						Str("User-Agent", c.GetHeader("User-Agent")).
						Str("sensor", msg.Name).
						Float64("value", msg.Value).
						Msgf("Broadcasting %s: %v", msg.Name, msg.Value)
				}
				c.SSEvent(msg.Name, msg)
				return true
			}
			return false
		})
	})
}
