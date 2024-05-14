package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (tr *Transport) exposeIndex() {
	tr.engine.GET("/", func(c *gin.Context) {
		sensors := tr.SensorsService.List()
		c.Header("Content-Type", "text/plain; version=0.0.4")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"sensors": sensors,
		})
	})
}
